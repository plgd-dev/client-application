package server

import (
	"fmt"

	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/pkg/net/grpc/server"
	certManager "github.com/plgd-dev/hub/v2/pkg/security/certManager/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func New(config Config, logger log.Logger, opts ...grpc.ServerOption) (*server.Server, error) {
	v := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(config.EnforcementPolicy.ToGrpc()),
		grpc.KeepaliveParams(config.KeepAlive.ToGrpc()),
	}
	var tlsClose func()
	if config.TLS.Enabled {
		tls, err := certManager.New(config.TLS.Config, logger)
		if err != nil {
			return nil, fmt.Errorf("cannot create cert manager %w", err)
		}
		v = append(v, grpc.Creds(credentials.NewTLS(tls.GetTLSConfig())))
		tlsClose = tls.Close
	} else {
		v = append(v, grpc.Creds(insecure.NewCredentials()))
	}

	if len(opts) > 0 {
		v = append(v, opts...)
	}
	server, err := server.NewServer(config.Addr, v...)
	if err != nil {
		if tlsClose != nil {
			tlsClose()
		}
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}
	server.AddCloseFunc(tlsClose)

	return server, nil
}
