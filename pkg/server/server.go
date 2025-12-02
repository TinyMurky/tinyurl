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
	"fmt"
	"net"
	"strconv"
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
	tcpListener, error := net.Listen("tcp", addr)

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
