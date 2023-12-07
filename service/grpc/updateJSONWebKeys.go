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

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/plgd-dev/hub/v2/identity-store/events"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc"
	plgdJwt "github.com/plgd-dev/hub/v2/pkg/security/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *ClientApplicationServer) updateJwkCache(newCache *JSONWebKeyCache) error {
	for {
		c := s.jwksCache.Load()
		if c != nil && c.owner != newCache.owner {
			return status.Errorf(codes.InvalidArgument, "cannot update jwks cache with owner %v for different owner %v", c.owner, newCache.owner)
		}
		if s.jwksCache.CompareAndSwap(c, newCache) {
			return nil
		}
	}
}

func (s *ClientApplicationServer) getOwnerForUpdateJSONWebKeys(ctx context.Context) (string, error) {
	token, err := grpc.TokenFromMD(ctx)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "cannot get token: %v", err)
	}
	owner := ""
	cfg := s.GetConfig()
	if s.jwksCache.Load() != nil {
		scopedClaims := plgdJwt.NewScopeClaims()
		err = s.ParseWithClaims(token, scopedClaims)
		if err != nil {
			return "", status.Errorf(codes.Unauthenticated, "cannot parse token: %v", err)
		}
		claims := plgdJwt.Claims(*scopedClaims)
		owner = claims.Owner(cfg.RemoteProvisioning.GetJwtOwnerClaim())
		if owner == "" {
			return "", status.Errorf(codes.Unauthenticated, "cannot get owner from token: claim %v is not set", cfg.RemoteProvisioning.GetJwtOwnerClaim())
		}
	} else {
		owner, err = grpc.OwnerFromTokenMD(ctx, cfg.RemoteProvisioning.GetJwtOwnerClaim())
		if err != nil {
			return "", status.Errorf(codes.Unauthenticated, "cannot get owner from token: %v", err)
		}
	}
	return owner, nil
}

func (s *ClientApplicationServer) UpdateJSONWebKeys(ctx context.Context, jwksReq *structpb.Struct) error {
	owner, err := s.getOwnerForUpdateJSONWebKeys(ctx)
	if err != nil {
		return err
	}

	keys, err := jwksReq.MarshalJSON()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot marshal keys: %v", err)
	}

	jwks, err := jwk.Parse(keys)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot marshal keys: %v", err)
	}
	ownerUuid, err := uuid.Parse(events.OwnerToUUID(owner))
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot parse owner: %v", err)
	}
	if err := s.updateJwkCache(NewJSONWebKeyCache(ownerUuid, jwks)); err != nil {
		return err
	}
	return nil
}
