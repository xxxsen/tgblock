package cache

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type PBCodec struct {
}

func (c *PBCodec) Encode(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("should be proto.Message")
	}
	return proto.Marshal(msg)
}

func (c *PBCodec) Decode(data []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("should be proto.Message")
	}
	return proto.Unmarshal(data, msg)
}
