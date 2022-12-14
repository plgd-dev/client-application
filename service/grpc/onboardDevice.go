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
	"crypto/x509"
	"fmt"

	"github.com/google/uuid"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/device/v2/client/core"
	"github.com/plgd-dev/device/v2/schema"
	"github.com/plgd-dev/device/v2/schema/acl"
	"github.com/plgd-dev/device/v2/schema/cloud"
	"github.com/plgd-dev/device/v2/schema/softwareupdate"
	"github.com/plgd-dev/kit/v2/security"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setACLForCloud(ctx context.Context, p *core.ProvisioningClient, cloudID string, links schema.ResourceLinks) error {
	link, err := core.GetResourceLink(links, acl.ResourceURI)
	if err != nil {
		return err
	}

	var acls acl.Response
	err = p.GetResource(ctx, link, &acls)
	if err != nil {
		return err
	}

	for _, acl := range acls.AccessControlList {
		if acl.Subject.Subject_Device != nil {
			if acl.Subject.Subject_Device.DeviceID == cloudID {
				return nil
			}
		}
	}
	confResources := acl.AllResources
	for _, href := range links.GetResourceHrefs(softwareupdate.ResourceType) {
		confResources = append(confResources, acl.Resource{
			Href:       href,
			Interfaces: []string{"*"},
		})
	}

	cloudACL := acl.UpdateRequest{
		AccessControlList: []acl.AccessControl{
			{
				Permission: acl.AllPermissions,
				Subject: acl.Subject{
					Subject_Device: &acl.Subject_Device{
						DeviceID: cloudID,
					},
				},
				Resources: confResources,
			},
		},
	}

	return p.UpdateResource(ctx, link, cloudACL, nil)
}

func validateOnboardDeviceRequest(req *pb.OnboardDeviceRequest) (uuid.UUID, error) {
	devID, err := strDeviceID2UUID(req.GetDeviceId())
	if err != nil {
		return uuid.UUID{}, err
	}
	if req.GetAuthorizationProviderName() == "" {
		return uuid.UUID{}, status.Error(codes.InvalidArgument, "invalid authorizationProviderName")
	}
	if req.GetCoapGatewayAddress() == "" {
		return uuid.UUID{}, status.Error(codes.InvalidArgument, "invalid coapGatewayAddress")
	}
	if req.GetAuthorizationCode() == "" {
		return uuid.UUID{}, status.Error(codes.InvalidArgument, "invalid authorizationCode")
	}
	if req.GetHubId() == "" {
		return uuid.UUID{}, status.Error(codes.InvalidArgument, "invalid hubId")
	}
	if req.GetCertificateAuthorities() != "" {
		_, err := security.ParseX509FromPEM([]byte(req.GetCertificateAuthorities()))
		if err != nil {
			return uuid.UUID{}, status.Errorf(codes.InvalidArgument, "invalid certificateAuthorities: %v", err)
		}
	}
	return devID, nil
}

func (s *ClientApplicationServer) getDeviceForSetupCloud(ctx context.Context, devID uuid.UUID) (*device, schema.ResourceLinks, error) {
	dev, err := s.getDevice(devID)
	if err != nil {
		return nil, nil, err
	}
	links, err := dev.getResourceLinksAndRefreshCache(ctx)
	if err != nil {
		return nil, nil, err
	}
	cloudLinks := links.GetResourceLinks(cloud.ResourceType)
	if len(cloudLinks) == 0 {
		return nil, nil, status.Errorf(codes.NotFound, "cannot find cloud resource for device %v", devID)
	}
	if err = dev.checkAccess(cloudLinks[0]); err != nil {
		return nil, nil, err
	}
	return dev, links, nil
}

func (s *ClientApplicationServer) OnboardDevice(ctx context.Context, req *pb.OnboardDeviceRequest) (resp *pb.OnboardDeviceResponse, err error) {
	devID, err := validateOnboardDeviceRequest(req)
	if err != nil {
		return nil, err
	}
	dev, links, err := s.getDeviceForSetupCloud(ctx, devID)
	if err != nil {
		return nil, err
	}
	if err = dev.provision(ctx, links, func(ctx context.Context, pc *core.ProvisioningClient) error {
		if errPro := setACLForCloud(ctx, pc, req.GetHubId(), links); errPro != nil {
			return errPro
		}
		if req.GetCertificateAuthorities() != "" {
			if errPro := pc.AddCertificateAuthority(ctx, req.GetHubId(), &x509.Certificate{
				Raw: []byte(req.GetCertificateAuthorities()),
			}); errPro != nil {
				return errPro
			}
		}
		return pc.SetCloudResource(ctx, cloud.ConfigurationUpdateRequest{
			AuthorizationProvider: req.GetAuthorizationProviderName(),
			AuthorizationCode:     req.GetAuthorizationCode(),
			URL:                   req.GetCoapGatewayAddress(),
			CloudID:               req.GetHubId(),
		})
	}); err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot provision onboard configuration for device %v: %w", dev.ID, err)).Err()
	}
	return &pb.OnboardDeviceResponse{}, nil
}
