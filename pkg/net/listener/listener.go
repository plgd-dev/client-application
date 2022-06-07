package listener

import (
	"fmt"
	"net"

	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/pkg/net/listener"
)

type Listener interface {
	AddCloseFunc(f func())
	Close() error
	Accept() (net.Conn, error)
	Addr() net.Addr
}

// Server handles gRPC requests to the service.
type Server struct {
	listener  net.Listener
	closeFunc []func()
}

// NewServer instantiates a listen server.
// When passing addr with an unspecified port or ":", use Addr().
func New(config Config, logger log.Logger) (Listener, error) {
	if config.TLS.Enabled {
		return listener.New(listener.Config{
			Addr: config.Addr,
			TLS:  config.TLS.Config,
		}, logger)
	}

	lis, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("listening failed: %w", err)
	}

	return &Server{listener: lis}, nil
}

// AddCloseFunc adds a function to be called by the Close method.
// This eliminates the need for wrapping the Server.
func (s *Server) AddCloseFunc(f func()) {
	s.closeFunc = append(s.closeFunc, f)
}

func (s *Server) Close() error {
	err := s.listener.Close()
	for _, f := range s.closeFunc {
		f()
	}
	return err
}

func (s *Server) Accept() (net.Conn, error) {
	return s.listener.Accept()
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}
