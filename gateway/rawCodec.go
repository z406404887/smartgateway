package gateway

import (
	"github.com/golang/protobuf/proto"
)

//rawCodec 原生codec类型
type rawCodec struct {
}

type bridge struct {
	rawBytes []byte
}

//NewRawCodec 返回原生编码
func NewRawCodec() (rc *rawCodec) {
	return &rawCodec{}
}

// 序列化
func (c *rawCodec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*bridge)
	if !ok {
		return proto.Marshal(v.(proto.Message))
	}
	return out.rawBytes, nil

}

// 反序列化
func (c *rawCodec) Unmarshal(data []byte, v interface{}) error {
	dst, ok := v.(*bridge)
	if !ok {
		return proto.Unmarshal(data, v.(proto.Message))
	}
	dst.rawBytes = data
	return nil
}

func (c *rawCodec) String() string {
	return "The RawCodec"
}
