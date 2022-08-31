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
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/hub/v2/http-gateway/serverMux"
	"github.com/plgd-dev/hub/v2/http-gateway/uri"
	kitNetGrpc "github.com/plgd-dev/hub/v2/pkg/net/grpc"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
)

func createContentBody(body io.ReadCloser) (io.ReadCloser, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	req := commands.Content{
		ContentType:       message.AppJSON.String(),
		CoapContentFormat: int32(message.AppJSON),
		Data:              data,
	}
	reqData, err := protojson.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal to protojson: %w", err)
	}

	return io.NopCloser(bytes.NewReader(reqData)), nil
}

func (requestHandler *RequestHandler) updateResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars[uri.DeviceIDKey]
	href := vars[uri.ResourceHrefKey]

	contentType := r.Header.Get(uri.ContentTypeHeaderKey)
	if contentType == uri.ApplicationProtoJsonContentType {
		requestHandler.mux.ServeHTTP(w, r)
		return
	}

	newBody, err := createContentBody(r.Body)
	if err != nil {
		serverMux.WriteError(w, kitNetGrpc.ForwardErrorf(codes.InvalidArgument, "cannot update resource('%v%v'): %v", deviceID, href, err))
		return
	}

	r.Body = newBody
	requestHandler.mux.ServeHTTP(w, r)
}
