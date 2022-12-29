package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/ribincao/ribin-game-server/codec"
	errs "github.com/ribincao/ribin-game-server/error"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-game-server/utils"
	"github.com/ribincao/ribin-protocol/base"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type OnCloseFunc func(conn *network.WrapConnection)
type OnConnectFunc func(conn *network.WrapConnection) bool
type Handler func(ctx context.Context, conn *network.WrapConnection, req *base.Test) (*base.Test, error)

var upgrader = websocket.Upgrader{}

type roomServer struct {
	opts                *ServerOptions
	ConnCloseCallback   OnCloseFunc
	ConnConnectCallback OnConnectFunc
	MessageHandler      Handler
}

func (s *roomServer) Close() {
}

func (s *roomServer) SetHandler(handler Handler) {
	s.MessageHandler = handler
}

func (s *roomServer) GetPort() string {
	if s.opts.address == "" {
		return ""
	}
	return strings.Split(s.opts.address, ":")[1]
}

func (s *roomServer) Serve() {
	http.Serve(s.opts.listener, s)
}

func (s *roomServer) SetConnCloseCallback(closeFunc OnCloseFunc) {
	s.ConnCloseCallback = closeFunc
}

func (s *roomServer) SetConnConnectCallback(connectFunc OnConnectFunc) {
	s.ConnConnectCallback = connectFunc
}

func (s *roomServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	utils.GoWithRecover(func() {
		conn := network.NewWrapConn(c, "", s)
		conn.Run()
	})
}

func (s *roomServer) OnConnect(conn *network.WrapConnection) bool {
	if s.ConnConnectCallback == nil {
		return true
	}
	return s.ConnConnectCallback(conn)
}

func (s *roomServer) OnClose(conn *network.WrapConnection) {
	if s.ConnCloseCallback == nil {
		conn.Close()
		return
	}
	s.ConnCloseCallback(conn)
}

func (s *roomServer) GetOpt() *ServerOptions {
	return s.opts
}

func (s *roomServer) OnMessage(c *network.WrapConnection, packet *network.Message) bool {
	frame, err := codec.DefaultCodec.Decode(packet.Data)
	if err != nil {
		return false
	}
	data, err := s.handleFrame(c, frame)
	if err != nil && err.Error() != "" {
		if e, ok := err.(*errs.Error); ok && e.Code != errs.RoomUnexistErrorCode && e.Code != errs.PlayerNotInRoomErrorCode {
			logger.Error("handleFrame error", zap.Error(e))
		}
	}

	rspbuf, _ := codec.DefaultCodec.Encode(data, codec.RPC)
	err = c.Write(packet.MsgType, rspbuf)
	return err == nil
}

func (s *roomServer) handleFrame(conn *network.WrapConnection, frame *codec.Frame) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.opts.timeout)
	defer cancel()

	req := &base.Test{}
	err := proto.Unmarshal(frame.Data, req)
	if err != nil {
		return nil, err
	}

	rsp, err := s.MessageHandler(ctx, conn, req)
	rspbuf, _ := proto.Marshal(rsp)
	if err != nil {
		return rspbuf, err
	}
	return rspbuf, nil
}
