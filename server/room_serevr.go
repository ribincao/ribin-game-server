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
)

type OnCloseFunc func(conn *network.WrapConnection)
type OnConnectFunc func(conn *network.WrapConnection) bool
type Handler func(ctx context.Context, conn *network.WrapConnection, req *base.Client2ServerReq) (*base.Server2ClientRsp, error)

var upgrader = websocket.Upgrader{}

type RoomServer struct {
	opts                *ServerOptions
	ConnCloseCallback   OnCloseFunc
	ConnConnectCallback OnConnectFunc
	MessageHandler      Handler
	CodeType            string
	MarshalType         string
}

func (s *RoomServer) Close() {
}

func (s *RoomServer) SetHandler(handler Handler) {
	s.MessageHandler = handler
}

func (s *RoomServer) GetPort() string {
	if s.opts.address == "" {
		return ""
	}
	return strings.Split(s.opts.address, ":")[1]
}

func (s *RoomServer) SetCodecType(codecType string) {
	s.CodeType = codecType
}

func (s *RoomServer) SetMarshalType(marshalType string) {
	s.MarshalType = marshalType
}

func (s *RoomServer) Serve() {
	http.Serve(s.opts.listener, s)
	logger.Info("Server Start, ", zap.String("Address", s.GetOpt().address))
}

func (s *RoomServer) SetConnCloseCallback(closeFunc OnCloseFunc) {
	s.ConnCloseCallback = closeFunc
}

func (s *RoomServer) SetConnConnectCallback(connectFunc OnConnectFunc) {
	s.ConnConnectCallback = connectFunc
}

func (s *RoomServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	utils.GoWithRecover(func() {
		conn := network.NewWrapConn(c, "", s)
		conn.Run()
	})
}

func (s *RoomServer) OnConnect(conn *network.WrapConnection) bool {
	if s.ConnConnectCallback == nil {
		return true
	}
	return s.ConnConnectCallback(conn)
}

func (s *RoomServer) OnClose(conn *network.WrapConnection) {
	if s.ConnCloseCallback == nil {
		conn.Close()
		return
	}
	s.ConnCloseCallback(conn)
}

func (s *RoomServer) GetOpt() *ServerOptions {
	return s.opts
}

func (s *RoomServer) OnMessage(c *network.WrapConnection, packet *network.Message) bool {
	frame, err := codec.GetCodec(s.CodeType).Decode(packet.Data)
	if err != nil {
		return false
	}
	data, err := s.handleFrame(c, frame)
	if err != nil && err.Error() != "" {
		if e, ok := err.(*errs.Error); ok && e.Code != errs.RoomUnexistErrorCode && e.Code != errs.PlayerNotInRoomErrorCode {
			logger.Error("handleFrame error", zap.Error(e))
		}
	}

	rspbuf, _ := codec.GetCodec(s.CodeType).Encode(data, codec.RPC)
	err = c.Write(packet.MsgType, rspbuf)
	return err == nil
}

func (s *RoomServer) handleFrame(conn *network.WrapConnection, frame *codec.Frame) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.opts.timeout)
	defer cancel()

	req := &base.Client2ServerReq{}
	err := codec.GetMarshaller(s.MarshalType).Unmarshal(frame.Data, req)
	if err != nil {
		return nil, err
	}

	rsp, err := s.MessageHandler(ctx, conn, req)
	rspbuf, _ := codec.GetMarshaller(s.MarshalType).Marshal(rsp)
	if err != nil {
		return rspbuf, err
	}
	return rspbuf, nil
}
