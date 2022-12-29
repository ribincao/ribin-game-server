package test

import (
	"fmt"
	"testing"

	codec "github.com/ribincao/ribin-game-server/codec"

	"github.com/ribincao/ribin-protocol/base"
	"google.golang.org/protobuf/proto"
)

type TestMessage struct {
	Id   int
	Name string
}

func TestCodecEncode(t *testing.T) {
	message := &base.Test{
		Seq: "test",
	}
	rspbuf, _ := proto.Marshal(message)

	c := codec.NewCodec()
	buf, _ := c.Encode(rspbuf, codec.RPC)
	fmt.Println("Encode ", buf)
}

func TestCodecDecode(t *testing.T) {
	message := &base.Test{
		Seq: "test",
	}
	rspbuf, _ := proto.Marshal(message)

	c := codec.NewCodec()
	buf, _ := c.Encode(rspbuf, codec.RPC)
	fmt.Println("Encode ", buf)

	frame, _ := c.Decode(buf)
	out := &base.Test{}
	proto.Unmarshal(frame.Data, out)
	fmt.Println("Decode ", out.Seq)
}
