package network

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	errs "github.com/ribincao/ribin-game-server/error"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/utils"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

const DefaultReadTimeout = 2 * time.Second
const DefaultWriteTimeout = 2 * time.Second
const DefaultWriteBuff = 1024

type Message struct {
	MsgType int
	Data    []byte
}

type ConnCallback interface {
	OnConnect(*WrapConnection) bool
	OnMessage(*WrapConnection, *Message) bool
	OnClose(*WrapConnection)
}

type WrapConnection struct {
	Connection     *Connection
	PlayerId       string
	IsClosed       atomic.Bool
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	LastActiveTime *atomic.Int64
	msgChan        chan *Message
	server         ConnCallback
	closeOnce      sync.Once
}

func (wc *WrapConnection) Read() (messageType int, p []byte, err error) {
	if wc.IsClosed.Load() {
		return 0, nil, errs.ConnectionClosedError
	}
	msgType, data, err := wc.Connection.ReadMessage()
	if err != nil {
		logger.Warn("ConnectionReadFailWarn",
			zap.Any("Connection", wc),
			zap.String("ErrMsg", err.Error()))
		wc.Close()
		return 0, nil, err
	}
	wc.UpdateLastActiveTime(time.Now().UnixMilli())

	return msgType, data, err
}

func (wc *WrapConnection) readInLoop() {
	for {
		if wc.IsClosed.Load() {
			return
		}
		msgType, data, err := wc.Connection.ReadMessage()
		if err != nil {
			logger.Warn("ConnectionReadInLoopFailWarn",
				zap.Any("Connection", wc),
				zap.String("ErrMsg", err.Error()))
			wc.Close()
			return
		}
		packet := &Message{
			MsgType: msgType,
			Data:    data,
		}

		wc.UpdateLastActiveTime(time.Now().UnixMilli())
		if !wc.server.OnMessage(wc, packet) {
			wc.Close()
			return
		}
	}
}

func (wc *WrapConnection) Write(messageType int, data []byte) (err error) {
	if wc.IsClosed.Load() {
		return errs.ConnectionClosedError
	}
	packet := &Message{
		MsgType: messageType,
		Data:    data,
	}
	utils.GoWithRecover(func() {
		err := wc.Connection.WriteMessage(packet.MsgType, packet.Data)
		if err != nil {
			logger.Warn("ConnectionWriteFailWarn",
				zap.Any("Connection", wc),
				zap.String("ErrMsg", err.Error()))
			wc.Close()
			return
		}
		wc.UpdateLastActiveTime(time.Now().UnixMilli())
	})
	return nil
}
func (wc *WrapConnection) writeInLoop() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case packet := <-wc.msgChan:
			if wc.IsClosed.Load() {
				return
			}
			utils.GoWithRecover(func() {
				err := wc.Connection.WriteMessage(packet.MsgType, packet.Data)
				if err != nil {
					logger.Warn("ConnectionWriteInLoopFailWarn",
						zap.Any("Connection", wc),
						zap.String("ErrMsg", err.Error()))
					wc.Close()
					return
				}
				wc.UpdateLastActiveTime(time.Now().UnixMilli())
			})
		case <-ticker.C:
			if wc.IsClosed.Load() {
				return
			}
		}
	}
}

func (wc *WrapConnection) Close() {
	wc.closeOnce.Do(func() {
		wc.IsClosed.Store(true)
		if wc.Connection != nil {
			wc.Connection.Close()
		}
		close(wc.msgChan)
		wc.server.OnClose(wc)
	})
}

func (wc *WrapConnection) OnConnect() {
	if wc.Connection != nil {
		wc.Connection.onConnect()
	}
}

func (wc *WrapConnection) AddOnConnectHandler(name string, function *Functor) {
	wc.Connection.RegisterConnectFunctor(name, function)
}

func (wc *WrapConnection) AddOnCloseHandler(name string, function *Functor) {
	wc.Connection.RegisterCloseFunctor(name, function)
}

func (wc *WrapConnection) UpdateLastActiveTime(timestamp int64) {
	wc.LastActiveTime.Store(timestamp)
}

func (wc *WrapConnection) Run() {
	if !wc.server.OnConnect(wc) {
		return
	}

	utils.GoWithRecover(func() {
		wc.readInLoop()
	})
	utils.GoWithRecover(func() {
		wc.writeInLoop()
	})
}

func NewWrapConn(conn *websocket.Conn, playerId string, server ConnCallback) *WrapConnection {
	c := &Connection{
		WebsocketConn: conn,
	}
	c.InitCloseHandler()
	wc := &WrapConnection{
		Connection:     c,
		PlayerId:       playerId,
		ReadTimeout:    DefaultReadTimeout,
		WriteTimeout:   DefaultWriteTimeout,
		LastActiveTime: atomic.NewInt64(time.Now().UnixMilli()),
		msgChan:        make(chan *Message, DefaultWriteBuff),
		server:         server,
	}
	return wc
}
