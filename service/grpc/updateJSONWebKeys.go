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

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/remoteProvisioning"
	"github.com/plgd-dev/hub/v2/identity-store/events"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type JSONWebKeyCache struct {
	owner uuid.UUID
	keys  jwk.Set
}

func NewJSONWebKeyCache(owner uuid.UUID, keys jwk.Set) *JSONWebKeyCache {
	return &JSONWebKeyCache{
		owner: owner,
		keys:  keys,
	}
}

func (c *JSONWebKeyCache) GetKey(token *jwt.Token) (interface{}, error) {
	key, err := c.LookupKey(token)
	if err != nil {
		return nil, err
	}
	var v interface{}
	return v, key.Raw(&v)
}

func (c *JSONWebKeyCache) LookupKey(token *jwt.Token) (jwk.Key, error) {
	id, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing key id in token")
	}

	if c.keys == nil {
		return nil, fmt.Errorf("empty JWK cache")
	}
	if key, ok := c.keys.LookupKeyID(id); ok {
		if key.Algorithm() == token.Method.Alg() {
			return key, nil
		}
	}
	return nil, fmt.Errorf("could not find JWK")
}

func (s *ClientApplicationServer) ParseWithClaims(token string, claims jwt.Claims) error {
	c := s.jwksCache.Load()
	if c == nil {
		return status.Errorf(codes.Unauthenticated, "cannot validate token: missing JWK cache")
	}
	if token == "" {
		return status.Errorf(codes.Unauthenticated, "missing token")
	}

	_, err := jwt.ParseWithClaims(token, claims, c.GetKey)
	if err != nil {
		return fmt.Errorf("could not parse token: %w", err)
	}
	return nil
}

func (s *ClientApplicationServer) updateJwkCache(newCache *JSONWebKeyCache) error {
	for {
		c := s.jwksCache.Load()
		if c == nil {
			if s.jwksCache.CompareAndSwap(c, newCache) {
				return nil
			}
		} else {
			if c.owner == newCache.owner {
				if s.jwksCache.CompareAndSwap(c, newCache) {
					return nil
				}
			} else {
				return status.Errorf(codes.PermissionDenied, "cannot update keys for other owner")
			}
		}
	}
}

func (s *ClientApplicationServer) UpdateJSONWebKeys(ctx context.Context, req *structpb.Struct) (*pb.UpdateJSONWebKeysResponse, error) {
	if s.remoteProvisioningConfig.Mode != remoteProvisioning.Mode_UserAgent {
		return nil, status.Errorf(codes.Unimplemented, "not supported")
	}
	owner, err := grpc.OwnerFromTokenMD(ctx, s.remoteProvisioningConfig.Authorization.OwnerClaim)
	if err != nil {
		return nil, s.logger.LogAndReturnError(status.Errorf(codes.Unauthenticated, "cannot get owner from token: %v", err))
	}
	keys, err := req.MarshalJSON()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot marshal keys: %v", err)
	}

	jwks, err := jwk.Parse(keys)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot marshal keys: %v", err)
	}
	ownerUuid, err := uuid.Parse(events.OwnerToUUID(owner))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse owner: %v", err)
	}
	if err := s.updateJwkCache(NewJSONWebKeyCache(ownerUuid, jwks)); err != nil {
		return nil, err
	}
	return &pb.UpdateJSONWebKeysResponse{}, nil
}
