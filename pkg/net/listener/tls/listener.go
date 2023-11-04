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

package tls

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/plgd-dev/hub/v2/pkg/fsnotify"
	"github.com/plgd-dev/hub/v2/pkg/log"
	"github.com/plgd-dev/hub/v2/pkg/security/certManager/server"
)

type ConnectionType uint8

const (
	ConnectionTypeHTTP ConnectionType = iota
	ConnectionTypeTLS
)

type Conn struct {
	net.Conn
	ConnectionType ConnectionType
	buf            *bufio.Reader
}

func (c *Conn) Read(b []byte) (int, error) {
	return c.buf.Read(b)
}

type chConn struct {
	c   net.Conn
	err error
}

type SplitListener struct {
	net.Listener
	config *tls.Config
	cons   chan chConn
	logger log.Logger
}

// tlsRecordHeaderLooksLikeHTTP reports whether a TLS record header
// looks like it might've been a misdirected plaintext HTTP request.
func tlsRecordHeaderLooksLikeHTTP(hdr []byte) bool {
	if len(hdr) < 5 {
		return false
	}
	switch string(hdr[:5]) {
	case "GET /", "HEAD ", "POST ", "PUT /", "OPTIO":
		return true
	}
	return false
}

// prepareConnection accepts connections asynchronously to avoid blocking subsequent connection attempts due to potential blocking caused by c.Peek.
// The use of asynchronous connection acceptance ensures that the program can continue accepting incoming connections while waiting for c.Peek to complete
func (l *SplitListener) prepareConnection(c net.Conn) {
	// buffer reads on our conn
	bconn := &Conn{
		Conn: c,
		buf:  bufio.NewReader(c),
	}
	n := 5
	_ = c.SetReadDeadline(time.Now().Add(1 * time.Second))
	hdr, err := bconn.buf.Peek(n)
	_ = c.SetReadDeadline(time.Time{})
	if err != nil {
		l.logger.Warnf("closing connection for error while peeking for first 5bytes in 1sec: %v", err)
		_ = bconn.Close()
		return
	}
	// The first 5 bytes is a TLS header, so we can use it to check for HTTP
	if tlsRecordHeaderLooksLikeHTTP(hdr) {
		bconn.ConnectionType = ConnectionTypeHTTP
		l.cons <- chConn{c: bconn, err: nil}
		return
	}
	bconn.ConnectionType = ConnectionTypeTLS
	l.cons <- chConn{c: tls.Server(bconn, l.config), err: nil}
}

func (l *SplitListener) run() {
	for {
		c, err := l.Listener.Accept()
		if err != nil {
			l.cons <- chConn{nil, err}
			return
		}
		go l.prepareConnection(c)
	}
}

func (l *SplitListener) Accept() (net.Conn, error) {
	d := <-l.cons
	return d.c, d.err
}

func NewSplitListener(l net.Listener, config *tls.Config, logger log.Logger) net.Listener {
	sl := &SplitListener{
		Listener: l,
		config:   config,
		cons:     make(chan chConn),
		logger:   logger,
	}
	go sl.run()
	return sl
}

// Server handles gRPC requests to the service.
type Server struct {
	listener  net.Listener
	closeFunc []func()
}

// New instantiates a listening server.
// When passing addr with an unspecified port or ":", use Addr().
func New(config Config, fileWatcher *fsnotify.Watcher, logger log.Logger) (*Server, error) {
	certManager, err := server.New(config.TLS, fileWatcher, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create cert manager %w", err)
	}

	lis, err := net.Listen("tcp", config.Addr)
	if err != nil {
		certManager.Close()
		return nil, fmt.Errorf("listening failed: %w", err)
	}
	splitListener := NewSplitListener(lis, certManager.GetTLSConfig(), logger)

	return &Server{listener: splitListener, closeFunc: []func(){certManager.Close}}, nil
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
