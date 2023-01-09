package server

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Server interface {
	Serve()
	Close()
	GetPort() string
	GetOpt() *ServerOptions
	SetConnConnectCallback(OnConnectFunc)
	SetConnCloseCallback(OnCloseFunc)
	SetHandler(Handler)
}

type ServerType uint32

const (
	RoomServer  ServerType = 1
	FrameServer ServerType = 2
)

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
	address string // listening address, e.g. :( ip://127.0.0.1:8080、 dns://www.google.com)
	// network           string        // network type, e.g. : tcp、udp
	// serializationType string        // serialization type, default: proto
	timeout  time.Duration // timeout
	listener net.Listener  // net listener
	// httpWriter        http.ResponseWriter
	// httpRequest       *http.Request
}

type ServerOption func(*ServerOptions)

func WithAddress(address string) ServerOption {
	return func(o *ServerOptions) {
		o.address = address
	}
}

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
