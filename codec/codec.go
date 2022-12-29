package codec

import (
	"bytes"
	"encoding/binary"

	errs "github.com/ribincao/ribin-game-server/error"
)

type Codec interface {
	Encode([]byte, MsgType) ([]byte, error)
	Decode([]byte) (*Frame, error)
}

type MsgType int32
type DefaultCodec struct{}

const (
	RPC          MsgType = 1
	Broadcast    MsgType = 2
	FrameHeadLen         = 6
)

var codecMap = make(map[string]Codec)
var defaultCodec = NewCodec()
var NewCodec = func() Codec {
	return &DefaultCodec{}
}

type FrameHeader struct {
	MsgTypeStart uint8  // start of frame
	Length       uint32 // total packet length
	MsgTypeEnd   uint8  // end of frame
}

type Frame struct {
	Header *FrameHeader // header of Frame
	Data   []byte       // serialized data
}

func GetCodec(name string) Codec {
	if codec, ok := codecMap[name]; ok {
		return codec
	}
	return defaultCodec
}

func RegisterCodec(name string, codec Codec) {
	if codecMap == nil {
		codecMap = make(map[string]Codec)
	}
	codecMap[name] = codec
}

// |   1   |   4   |  ... |  1  |
// |   X   |  XXXX |  ... |  X  |
// | START |  LEN  | DATA | END |
func (c *DefaultCodec) Encode(data []byte, msgType MsgType) ([]byte, error) {

	totalLen := FrameHeadLen + len(data)
	buffer := bytes.NewBuffer(make([]byte, 0, totalLen))

	var msgTypeStart uint8 = 0x11
	var msgTypeEnd uint8 = 0x12
	if msgType == Broadcast {
		msgTypeStart = 0x21
		msgTypeEnd = 0x22
	}
	frame := FrameHeader{
		MsgTypeStart: msgTypeStart,
		Length:       uint32(len(data)),
		MsgTypeEnd:   msgTypeEnd,
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.MsgTypeStart); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Length); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.MsgTypeEnd); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *DefaultCodec) Decode(frameBytes []byte) (*Frame, error) {
	dataLen := binary.BigEndian.Uint32(frameBytes[1:5])
	if uint32(len(frameBytes)) < dataLen+5 {
		return nil, errs.MsgError
	}
	frame := &Frame{
		Header: &FrameHeader{
			MsgTypeStart: frameBytes[0],
			Length:       dataLen,
			MsgTypeEnd:   frameBytes[len(frameBytes)-1],
		},
		Data: frameBytes[5 : 5+dataLen],
	}
	return frame, nil
}
