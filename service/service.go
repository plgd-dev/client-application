package service

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/plgd-dev/client-application/service/http"
	"github.com/plgd-dev/go-coap/v2/tcp"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	config         Config
	listener       tcp.Listener
	ctx            context.Context
	cancel         context.CancelFunc
	logger         log.Logger
	httpService    *http.Service
	tracerProvider trace.TracerProvider
	sigs           chan os.Signal
}

const DPSTag = "dps"

func setCAPools(roots *x509.CertPool, intermediates *x509.CertPool, certs []*x509.Certificate) {
	for _, cert := range certs {
		if !cert.IsCA {
			continue
		}
		if bytes.Equal(cert.RawIssuer, cert.RawSubject) {
			roots.AddCert(cert)
			continue
		}
		intermediates.AddCert(cert)
	}
}

const serviceName = "client-application"

// New creates server.
func New(ctx context.Context, config Config, logger log.Logger) (*Service, error) {
	tracerProvider := trace.NewNoopTracerProvider()
	httpService, err := http.New(ctx, serviceName, config.APIs.HTTP, logger, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("cannot create http service: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	s := Service{
		config: config,

		sigs: make(chan os.Signal, 1),

		ctx:    ctx,
		cancel: cancel,

		logger:         logger,
		httpService:    httpService,
		tracerProvider: tracerProvider,
	}

	return &s, nil
}

// Serve starts a device provisioning service on the configured address in *Service.
func (server *Service) Serve() error {
	return server.serveWithHandlingSignal()
}

func (server *Service) serveWithHandlingSignal() error {
	var wg sync.WaitGroup
	errCh := make(chan error, 3)
	wg.Add(1)
	go func(server *Service) {
		defer wg.Done()
		err := server.httpService.Serve()
		errCh <- err
	}(server)

	signal.Notify(server.sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-server.sigs

	err := server.httpService.Shutdown()
	errCh <- err
	wg.Wait()

	var errors []error
	for {
		select {
		case err := <-errCh:
			if err != nil {
				errors = append(errors, err)
			}
		default:
			switch len(errors) {
			case 0:
				return nil
			case 1:
				return errors[0]
			default:
				return fmt.Errorf("%v", errors)
			}
		}
	}
}

// Shutdown turn off server.
func (server *Service) Close() error {
	select {
	case server.sigs <- syscall.SIGTERM:
	default:
	}
	return nil
}
