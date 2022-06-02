package test

import (
	"context"
	"sync"
	"testing"

	_ "cloud.google.com/go"
	"github.com/plgd-dev/client-application/service"
	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/test/config"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

const (
	CLIENT_APPLICATIO_HTTP_HOST = "localhost:40050"
)

func MakeConfig(t require.TestingT) service.Config {
	var cfg service.Config
	cfg.Log = log.MakeDefaultConfig()
	cfg.APIs.HTTP = MakeHttpConfig()

	return cfg
}

func SetUp(t *testing.T) (TearDown func()) {
	return New(t, MakeConfig(t))
}

// New creates test coap-gateway.
func New(t *testing.T, cfg service.Config) func() {
	ctx := context.Background()
	logger := log.NewLogger(cfg.Log)

	s, err := service.New(ctx, cfg, logger)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = s.Serve()
	}()

	return func() {
		_ = s.Close()
		wg.Wait()
	}
}

func MakeHttpConfig() http.Config {
	return http.Config{
		Connection: config.MakeListenerConfig(CLIENT_APPLICATIO_HTTP_HOST),
	}
}

func NewHttpService(ctx context.Context, t *testing.T) (*http.Service, func()) {
	cfg := MakeConfig(t)
	cfg.APIs.HTTP.Connection.TLS.ClientCertificateRequired = false
	logger := log.NewLogger(cfg.Log)

	s, err := http.New(ctx, "dps-http", cfg.APIs.HTTP, logger, trace.NewNoopTracerProvider())
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = s.Serve()
	}()

	cleanUp := func() {
		err = s.Shutdown()
		require.NoError(t, err)
		wg.Wait()
	}

	return s, cleanUp
}
