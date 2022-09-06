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

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/client-application/service/remoteProvisioning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *ClientApplicationServer) GetJSONWebKeys(ctx context.Context, req *pb.GetJSONWebKeysRequest) (*structpb.Struct, error) {
	if s.remoteProvisioningConfig.Mode != remoteProvisioning.Mode_UserAgent {
		return nil, status.Errorf(codes.Unimplemented, "not supported")
	}
	jwksCache := s.jwksCache.Load()
	if jwksCache == nil {
		return nil, status.Errorf(codes.Unavailable, "not available")
	}
	keys := make([]jwk.Key, 0, jwksCache.keys.Len())
	for i := 0; i < jwksCache.keys.Len(); i++ {
		k, ok := jwksCache.keys.Get(i)
		if ok {
			keys = append(keys, k)
		}
	}
	resp, err := structpb.NewStruct(map[string]interface{}{
		"keys": keys,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert to struct: %v", err)
	}
	return resp, nil
}
