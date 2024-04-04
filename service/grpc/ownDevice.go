// ************************************************************************
// Copyright (C) 2022 plgd.dev, s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ************************************************************************

package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/schema"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"go.uber.org/atomic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type remoteSign struct {
	state               uuid.UUID
	closed              atomic.Bool
	errChan             chan error
	getCSRChan          chan []byte
	certificateSignChan chan []byte
	cancel              context.CancelFunc
	ctx                 context.Context
}

func (s *remoteSign) Close(err error) {
	if !s.closed.CompareAndSwap(false, true) {
		return
	}
	if err != nil {
		s.SendErr(err)
	}
	close(s.errChan)
	s.cancel()
	close(s.getCSRChan)
	close(s.certificateSignChan)
}

func (s *remoteSign) Sign(_ context.Context, csr []byte) ([]byte, error) {
	select {
	case <-s.ctx.Done():
		return nil, fmt.Errorf("cannot send request with CSR: %w", s.ctx.Err())
	case s.getCSRChan <- csr:
	}
	select {
	case <-s.ctx.Done():
		return nil, fmt.Errorf("cannot wait for certificate: %w", s.ctx.Err())
	case cert := <-s.certificateSignChan:
		return cert, nil
	}
}

func (s *remoteSign) SendErr(err error) {
	select {
	case s.errChan <- err:
	default:
	}
}

func (s *remoteSign) SendCertificate(ctx context.Context, certificate []byte) error {
	select {
	case err := <-s.errChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("cannot send certificate: %w", ctx.Err())
	case <-s.ctx.Done():
		return fmt.Errorf("cannot send certificate: %w", s.ctx.Err())
	case s.certificateSignChan <- certificate:
		return nil
	}
}

func (s *remoteSign) ReadError(ctx context.Context) error {
	select {
	case err := <-s.errChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("cannot read error: %w", ctx.Err())
	case <-s.ctx.Done():
		return fmt.Errorf("cannot read error: %w", s.ctx.Err())
	}
}

func (s *remoteSign) ReadCSR(ctx context.Context) ([]byte, error) {
	select {
	case err := <-s.errChan:
		return nil, err
	case csr := <-s.getCSRChan:
		return csr, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("cannot read CSR: %w", ctx.Err())
	case <-s.ctx.Done():
		return nil, fmt.Errorf("cannot read CSR: %w", s.ctx.Err())
	}
}

func deviceStateID(device uuid.UUID, state uuid.UUID) uuid.UUID {
	return uuid.NewSHA1(device, state[:])
}

func newRemoteSign(timeout time.Duration) *remoteSign {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &remoteSign{
		state:               uuid.New(),
		errChan:             make(chan error, 10),
		getCSRChan:          make(chan []byte),
		certificateSignChan: make(chan []byte),
		cancel:              cancel,
		ctx:                 ctx,
	}
}

func (s *ClientApplicationServer) ownDeviceGetCSR(ctx context.Context, timeoutValue int64, dev *device, links schema.ResourceLinks) (*pb.OwnDeviceResponse, error) {
	timeout := time.Second * 15
	if timeoutValue > 0 {
		timeout = time.Duration(timeoutValue) * time.Nanosecond
	}
	remoteSign := newRemoteSign(timeout)
	_, loaded := s.remoteOwnSignCache.LoadOrStore(deviceStateID(dev.ID, remoteSign.state), remoteSign)
	if loaded {
		remoteSign.Close(nil)
		return nil, status.Errorf(codes.Unavailable, "cannot get CSR: state %v is already in progress", remoteSign.state)
	}
	go func() {
		devService := s.serviceDevice.Load()
		if devService == nil {
			remoteSign.Close(errors.New("device service is not initialized"))
			return
		}
		defer s.remoteOwnSignCache.Delete(deviceStateID(dev.ID, remoteSign.state))
		ownOpts, err := devService.GetOwnOptions()
		if err != nil {
			remoteSign.Close(fmt.Errorf("cannot get own options: %w", err))
			return
		}
		ownOpts = append(ownOpts, core.WithSetupCertificates(remoteSign.Sign))
		err = dev.Own(remoteSign.ctx, links, devService.GetOwnershipClients(), ownOpts...)
		remoteSign.Close(err)
	}()
	csr, err := remoteSign.ReadCSR(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.OwnDeviceResponse{
		IdentityCertificateChallenge: &pb.IdentityCertificateChallenge{
			CertificateSigningRequest: csr,
			State:                     remoteSign.state.String(),
		},
	}, nil
}

func (s *ClientApplicationServer) ownDeviceSetCertificate(ctx context.Context, devID uuid.UUID, req *pb.FinishOwnDeviceRequest) error {
	state, err := uuid.Parse(req.GetState())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot parse state: %v", err)
	}

	remoteSign, ok := s.remoteOwnSignCache.Load(deviceStateID(devID, state))
	if !ok {
		return status.Errorf(codes.NotFound, "cannot find remote sign for state: %v", req.GetState())
	}
	if err = remoteSign.SendCertificate(ctx, req.GetCertificate()); err != nil {
		return err
	}
	if err = remoteSign.ReadError(ctx); err != nil {
		return err
	}
	return nil
}

func (s *ClientApplicationServer) OwnDevice(ctx context.Context, req *pb.OwnDeviceRequest) (*pb.OwnDeviceResponse, error) {
	devService := s.serviceDevice.Load()
	if devService == nil {
		return nil, status.Errorf(codes.Unavailable, "device service is not initialized")
	}
	devID, err := strDeviceID2UUID(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	dev, err := s.getDevice(devID)
	if err != nil {
		return nil, err
	}
	links, err := dev.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return nil, err
	}
	if s.signIdentityCertificateRemotely() {
		resp, err2 := s.ownDeviceGetCSR(ctx, req.GetTimeout(), dev, links)
		if err2 != nil {
			return nil, grpc.ForwardErrorf(codes.InvalidArgument, "cannot own device mediated by user agent: %v", err2)
		}
		return resp, nil
	}

	ownOptions, err := devService.GetOwnOptions()
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get own options: %w", err)).Err()
	}
	err = dev.Own(ctx, links, devService.GetOwnershipClients(), ownOptions...)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot own device %v: %w", dev.ID, err)).Err()
	}
	dev.updateOwnershipStatus(grpcgwPb.Device_OWNED)

	return &pb.OwnDeviceResponse{}, nil
}

func (s *ClientApplicationServer) FinishOwnDevice(ctx context.Context, req *pb.FinishOwnDeviceRequest) (*pb.FinishOwnDeviceResponse, error) {
	if !s.signIdentityCertificateRemotely() {
		return nil, status.Errorf(codes.Unimplemented, "initialize with certificate is disabled")
	}
	devID, err := strDeviceID2UUID(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	if err := s.ownDeviceSetCertificate(ctx, devID, req); err != nil {
		return nil, err
	}
	return &pb.FinishOwnDeviceResponse{}, nil
}
