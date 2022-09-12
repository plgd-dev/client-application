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
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/device/client/core"
	"github.com/plgd-dev/device/schema"
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
	getCSRChan          chan *pb.GetIdentityCSRResponse
	certificateSignChan chan *pb.UpdateIdentityCertificateRequest
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

func (s *remoteSign) Sign(ctx context.Context, csr []byte) ([]byte, error) {
	select {
	case <-s.ctx.Done():
		return nil, fmt.Errorf("cannot send request with CSR: %w", s.ctx.Err())
	case s.getCSRChan <- &pb.GetIdentityCSRResponse{
		CertificateSigningRequest: string(csr),
		State:                     s.state.String(),
	}:
	}
	select {
	case <-s.ctx.Done():
		return nil, fmt.Errorf("cannot wait for certificate: %w", s.ctx.Err())
	case cert := <-s.certificateSignChan:
		return []byte(cert.GetCertificate()), nil
	}
}

func (s *remoteSign) SendErr(err error) {
	select {
	case s.errChan <- err:
	default:
	}
}

func (s *remoteSign) SendCertificate(ctx context.Context, req *pb.UpdateIdentityCertificateRequest) error {
	select {
	case err := <-s.errChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("cannot send certificate: %w", ctx.Err())
	case <-s.ctx.Done():
		return fmt.Errorf("cannot send certificate: %w", ctx.Err())
	case s.certificateSignChan <- req:
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

func (s *remoteSign) ReadCSR(ctx context.Context) (*pb.GetIdentityCSRResponse, error) {
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
		errChan:             make(chan error, 1),
		getCSRChan:          make(chan *pb.GetIdentityCSRResponse),
		certificateSignChan: make(chan *pb.UpdateIdentityCertificateRequest),
		cancel:              cancel,
		ctx:                 ctx,
	}
}

func (s *ClientApplicationServer) ownDeviceGetCSR(ctx context.Context, req *pb.OwnDeviceRequest_GetIdentityCsr, dev *device, links schema.ResourceLinks) (*pb.OwnDeviceResponse, error) {
	timeout := time.Second * 15
	if req.GetTimeout() > 0 {
		timeout = time.Duration(req.GetTimeout()) * time.Nanosecond
	}
	remoteSign := newRemoteSign(timeout)
	_, loaded := s.remoteOwnSignCache.LoadOrStore(deviceStateID(dev.ID, remoteSign.state), remoteSign)
	if loaded {
		remoteSign.Close(nil)
		return nil, status.Errorf(codes.Unavailable, "cannot get CSR: state %v is already in progress", remoteSign.state)
	}
	go func() {
		ownOpts := s.serviceDevice.GetOwnOptions()
		ownOpts = append(ownOpts, core.WithSetupCertificates(remoteSign.Sign))
		err := dev.Own(remoteSign.ctx, links, s.serviceDevice.GetOwnershipClients(), ownOpts...)
		remoteSign.Close(err)
		s.remoteOwnSignCache.Delete(deviceStateID(dev.ID, remoteSign.state))
	}()
	csr, err := remoteSign.ReadCSR(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.OwnDeviceResponse{
		GetIdentityCsr: csr,
	}, nil
}

func (s *ClientApplicationServer) ownDeviceSetCertificate(ctx context.Context, devID uuid.UUID, req *pb.UpdateIdentityCertificateRequest) (*pb.OwnDeviceResponse, error) {
	state, err := uuid.Parse(req.GetState())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse state: %v", err)
	}

	remoteSign, ok := s.remoteOwnSignCache.Load(deviceStateID(devID, state))
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot find remote sign for state: %v", req.GetState())
	}
	if err = remoteSign.SendCertificate(ctx, req); err != nil {
		return nil, err
	}
	if err = remoteSign.ReadError(ctx); err != nil {
		return nil, err
	}
	return &pb.OwnDeviceResponse{}, nil
}

func (s *ClientApplicationServer) OwnDevice(ctx context.Context, req *pb.OwnDeviceRequest) (*pb.OwnDeviceResponse, error) {
	devID, err := strDeviceID2UUID(req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	if s.signIdentityCertificateRemotely() && req.GetSetIdentityCertificate() != nil {
		return s.ownDeviceSetCertificate(ctx, devID, req.GetSetIdentityCertificate())
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
		resp, err2 := s.ownDeviceGetCSR(ctx, req.GetGetIdentityCsr(), dev, links)
		if err2 != nil {
			return nil, grpc.ForwardErrorf(codes.InvalidArgument, "cannot own device mediated by user agent: %v", err2)
		}
		return resp, nil
	}

	err = dev.Own(ctx, links, s.serviceDevice.GetOwnershipClients(), s.serviceDevice.GetOwnOptions()...)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot own device %v: %w", dev.ID, err)).Err()
	}
	dev.updateOwnershipStatus(grpcgwPb.Device_OWNED)

	return &pb.OwnDeviceResponse{}, nil
}
