package server

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// Server defines a generic implementation interface for a specific Service
type Server interface {
	Serve()
	Close()
	GetPort() string
	GetOpt() *ServerOptions
}

// ServerType is the type of server
type ServerType uint32

const (
	RoomServer  ServerType = 1
	FrameServer ServerType = 2
)

// NewServer creates a server
func NewServer(serverType ServerType, opts ...ServerOption) Server {
	var s Server
	switch serverType {
	case RoomServer:
		s = &roomServer{
			opts: &ServerOptions{
				timeout: 10 * time.Second,
			},
		}

	}

	server := s
	for _, opt := range opts {
		opt(server.GetOpt())
	}

	if server.GetOpt().address != "" {
		listener, err := net.Listen("tcp", server.GetOpt().address)
		if err != nil {
			panic(err)
		}
		server.GetOpt().listener = listener
		return s
	}

	if listener, port, err := OpenFreePort(10000, 1000); err == nil {
		server.GetOpt().listener = listener
		server.GetOpt().address = fmt.Sprintf(":%d", port)
	}
	return s
}

type ServerOptions struct {
	address           string        // listening address, e.g. :( ip://127.0.0.1:8080、 dns://www.google.com)
	network           string        // network type, e.g. : tcp、udp
	serializationType string        // serialization type, default: proto
	timeout           time.Duration // timeout
	listener          net.Listener  // net listener
	httpWriter        http.ResponseWriter
	httpRequest       *http.Request
}

type ServerOption func(*ServerOptions)

// OpenFreePort opens free UDP port.
// This example does not actually use UDP ports,
// but to avoid port collisions with the HTTP server,
// it binds the same number of UDP port in advance.
func OpenFreePort(portBase int, num int) (net.Listener, int, error) {
	random := rand.Intn(num)
	for i := random; i < random+num; i++ {
		port := portBase + i
		listener, err := net.Listen("tcp", fmt.Sprint(":", port))
		if err != nil {
			continue
		}
		return listener, port, nil
	}
	return nil, 0, errors.New("failed to open free port")
}
