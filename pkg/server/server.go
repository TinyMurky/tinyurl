// Copyright 2020 the Exposure Notifications Server authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package server provides an opinionated http server.
//
// Although exported, this package is non intended for general consumption. It
// is a shared dependency between multiple exposure notifications projects. We
// cannot guarantee that there won't be breaking changes in the future.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// Server accecpt ip, port, net.Listener
// and it can server net/http
//
// Server provides a gracefully-stoppable http server implementation. It is safe
// for concurrent use in goroutines.
type Server struct {
	ip       string
	port     string
	listener net.Listener
}

// New creates a new server listening on the provided address that responds to
// the http.Handler. It starts the listener (TCP), but does not start the server. If
// an empty port is given, the server randomly chooses one.
func New(port string) (*Server, error) {
	// :port will listen to 0.0.0.0
	addr := fmt.Sprintf(":%s", port)
	tcpListener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	server := Server{
		ip:       tcpListener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(tcpListener.Addr().(*net.TCPAddr).Port),
		listener: tcpListener,
	}

	return &server, nil
}

// NewFromListener creates a new server on the given listener. This is useful if
// you want to customize the listener type
func NewFromListener(listener net.Listener) (*Server, error) {
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)

	if !ok {
		return nil, fmt.Errorf("listener is not TCP")
	}

	server := Server{
		ip:       tcpAddr.IP.String(),
		port:     strconv.Itoa(tcpAddr.Port),
		listener: listener,
	}

	return &server, nil
}

// ServeHTTP starts the server and blocks until the provided context is closed.
// When the provided context is closed, the server is gracefully stopped with a
// timeout of 5 seconds.
//
// Once a server has been stopped, it is NOT safe for reuse.
func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	logger := logging.FromContext(ctx)

	errChan := make(chan error, 1)

	// Spawn a goroutine that listens for context closure. When the context is
	// closed, the server is stopped.
	go func() {
		// block until context closed
		<-ctx.Done()

		logger.Debug("server.Serve: Context close")

		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Debug("server.Serve: Graceful Shutting Down...")

		errChan <- srv.Shutdown(shutdownCtx)

	}()

	// Run the server. This will block until the provided context is closed.
	// If server shutdown correctly, it will return ErrServerClosed
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	var mutiErr *multierror.Error

	if err := <-errChan; err != nil {
		mutiErr = multierror.Append(mutiErr, fmt.Errorf("failed to shutdown server: %w", err))
	}

	return mutiErr.ErrorOrNil()
}

// ServeHTTPHandler is a convenience wrapper around ServeHTTP. It creates an
// HTTP server using the provided handler
func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	})
}

// Addr returns the server's listening address (ip + port).
func (s *Server) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

// IP returns the server's listening IP.
func (s *Server) IP() string {
	return s.ip
}

// Port returns the server's listening port.
func (s *Server) Port() string {
	return s.port
}
