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

package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/plgd-dev/hub/v2/http-gateway/serverMux"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"github.com/plgd-dev/kit/v2/codec/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func (requestHandler *RequestHandler) updateJSONWebKeys(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		serverMux.WriteError(w, kitNetGrpc.ForwardErrorf(codes.InvalidArgument, "empty body"))
		return
	}
	var data map[string]interface{}
	err := json.ReadFrom(r.Body, &data)
	if err != nil {
		serverMux.WriteError(w, kitNetGrpc.ForwardErrorf(codes.InvalidArgument, "cannot read body: %v", err))
		return
	}
	req, err := structpb.NewStruct(data)
	if err != nil {
		serverMux.WriteError(w, kitNetGrpc.ForwardErrorf(codes.InvalidArgument, "cannot convert map to structpb.Struct: %v", err))
		return
	}
	res, err := protojson.Marshal(req)
	if err != nil {
		serverMux.WriteError(w, kitNetGrpc.ForwardErrorf(codes.InvalidArgument, "cannot marshal structpb.Struct to jsonpb: %v", err))
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(res))
	requestHandler.mux.ServeHTTP(w, r)
}
