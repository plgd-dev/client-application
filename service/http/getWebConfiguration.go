package http

import (
	"net/http"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/kit/v2/codec/json"
)

// WebConfiguration represents web configuration for user interface
type WebConfigurationConfig struct {
	HTTPGatewayAddress string `yaml:"httpGatewayAddress" json:"httpGatewayAddress"`
}

const contentTypeHeaderKey = "Content-Type"

func jsonResponseWriter(w http.ResponseWriter, v interface{}) error {
	if v == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	w.Header().Set(contentTypeHeaderKey, message.AppJSON.String())
	return json.WriteTo(w, v)
}

func getWebConfiguration(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	cfg := WebConfigurationConfig{
		HTTPGatewayAddress: scheme + "://" + r.Host,
	}
	if err := jsonResponseWriter(w, cfg); err != nil {
		log.Errorf("failed to write response: %w", err)
	}
}
