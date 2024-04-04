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
	"encoding/pem"

	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"google.golang.org/grpc/codes"
)

func (s *ClientApplicationServer) GetIdentityCertificate(_ context.Context, _ *pb.GetIdentityCertificateRequest) (*pb.GetIdentityCertificateResponse, error) {
	devService := s.serviceDevice.Load()
	if devService == nil {
		return nil, grpc.ForwardErrorf(codes.Unavailable, "service is not initialized")
	}
	tls, err := devService.GetIdentityCertificate()
	if err != nil {
		return nil, grpc.ForwardErrorf(codes.Unavailable, "cannot get identity certificate %v", err)
	}
	if len(tls.Certificate) == 0 {
		return nil, grpc.ForwardErrorf(codes.Unavailable, "cannot get identity certificate: certificate is not set")
	}
	certificate := make([]byte, 0, len(tls.Certificate[0])*len(tls.Certificate))
	for _, c := range tls.Certificate {
		certificate = append(certificate, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c})...)
	}
	return &pb.GetIdentityCertificateResponse{
		Certificate: certificate,
	}, nil
}
