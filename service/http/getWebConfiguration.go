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
	"fmt"
	"net/http"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/kit/v2/codec/json"
)

const contentTypeHeaderKey = "Content-Type"

func jsonResponseWriter(w http.ResponseWriter, v interface{}) error {
	if v == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	w.Header().Set(contentTypeHeaderKey, message.AppJSON.String())
	return json.WriteTo(w, v)
}

// WebConfigurationConfig represents web configuration for user interface exposed via getOAuthConfiguration handler
type WebConfigurationConfig struct {
	HTTPGatewayAddress string `yaml:"-" json:"httpGatewayAddress"`
}

func (requestHandler *RequestHandler) getWebConfiguration(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	cfg := WebConfigurationConfig{
		HTTPGatewayAddress: fmt.Sprintf("%s://%s", scheme, r.Host),
	}
	if err := jsonResponseWriter(w, cfg); err != nil {
		log.Errorf("failed to write response: %w", err)
	}
}
