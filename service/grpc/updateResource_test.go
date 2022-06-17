//***************************************************************************
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
//*************************************************************************

package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	serviceHttp "github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/client-application/test"
	"github.com/plgd-dev/device/schema/device"
	"github.com/plgd-dev/device/schema/doxm"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClientApplicationServerUpdateResource(t *testing.T) {
	dev := test.MustFindDeviceByName(test.DevsimName, []pb.GetDevicesRequest_UseMulticast{pb.GetDevicesRequest_IPV4})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	type args struct {
		req *pb.UpdateResourceRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *grpcgwPb.UpdateResourceResponse
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "doxm update",
			args: args{
				req: &pb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     doxm.ResourceURI,
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"oxmsel":0}`),
					},
				},
			},
			want: &grpcgwPb.UpdateResourceResponse{
				Data: &events.ResourceUpdated{
					Content: &commands.Content{},
					Status:  commands.Status_OK,
				},
			},
		},
		{
			name: "device resource - fail",
			args: args{
				req: &pb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     device.ResourceURI,
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.PermissionDenied,
		},
		{
			name: "unknown device",
			args: args{
				req: &pb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: uuid.NewString(),
						Href:     device.ResourceURI,
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
		{
			name: "unknown href",
			args: args{
				req: &pb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     "/unknown",
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
		{
			name: "permission denied - cannot establish TLS connection",
			args: args{
				req: &pb.UpdateResourceRequest{
					ResourceId: &commands.ResourceId{
						DeviceId: dev.Id,
						Href:     "/light/1",
					},
					Content: &grpcgwPb.Content{
						ContentType: serviceHttp.ApplicationJsonContentType,
						Data:        []byte(`{"name":"test"}`),
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.PermissionDenied,
		},
	}

	s, teardown, err := test.NewClientApplicationServer(ctx)
	require.NoError(t, err)
	defer teardown()
	err = s.GetDevices(&pb.GetDevicesRequest{}, test.NewClientApplicationGetDevicesServer(ctx))
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.UpdateResource(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrCode.String(), status.Code(err).String())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
