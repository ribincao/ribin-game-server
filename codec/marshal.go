package codec

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Marshaller interface {
	Marshal(m protoreflect.ProtoMessage) ([]byte, error)
	Unmarshal(b []byte, m protoreflect.ProtoMessage) error
}

type defaultMarshal struct{}

var marshalMap = make(map[string]Marshaller)
var DefaultMarshal = NewMarshal()
var NewMarshal = func() Marshaller {
	return &defaultMarshal{}
}

func GetMarshaller(name string) Marshaller {
	if marshaller, ok := marshalMap[name]; ok {
		return marshaller
	}
	return DefaultMarshal
}

func RegisterMarshaller(name string, marshaller Marshaller) {
	if marshalMap == nil {
		marshalMap = make(map[string]Marshaller)
	}
	marshalMap[name] = marshaller
}

func (c *defaultMarshal) Marshal(m protoreflect.ProtoMessage) ([]byte, error) {
	return proto.Marshal(m)
}

func (c *defaultMarshal) Unmarshal(b []byte, m protoreflect.ProtoMessage) error {
	return proto.Unmarshal(b, m)
}
