package test

import (
	"fmt"
	codec "ribin-server/codec"
	"testing"

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
	err := codec.DefaultSerialization.Unmarshal(frame.Payload, out)
	if err != nil {
		fmt.Println("Serialization Error", err)
		return
	}
	fmt.Println("Decode ", out.Seq)
}
