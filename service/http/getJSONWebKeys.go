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
	"net/http"
	"net/http/httptest"

	"github.com/plgd-dev/hub/v2/http-gateway/serverMux"
	"github.com/plgd-dev/hub/v2/pkg/log"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"github.com/plgd-dev/kit/v2/codec/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func writeError(w http.ResponseWriter, rec *httptest.ResponseRecorder) {
	// copy everything from response recorder
	// to actual response writer
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	w.WriteHeader(rec.Code)
	_, _ = w.Write(rec.Body.Bytes())
}

func (requestHandler *RequestHandler) getJSONWebKeys(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	requestHandler.mux.ServeHTTP(rec, r)
	if rec.Code != http.StatusOK {
		writeError(w, rec)
		return
	}
	var resp structpb.Struct
	err := protojson.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		serverMux.WriteError(w, kitNetGrpc.ForwardErrorf(codes.Internal, "cannot unmarshal response to structpb.Struct: %v", err))
		return
	}
	data, err := json.Encode(resp.AsMap())
	if err != nil {
		serverMux.WriteError(w, kitNetGrpc.ForwardErrorf(codes.Internal, "cannot marshal structpb.Struct to json-map: %v", err))
		return
	}
	if err := jsonResponseWriter(w, data); err != nil {
		log.Errorf("failed to write response: %w", err)
	}
}
