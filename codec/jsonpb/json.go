package jsonpb

import (
	"bytes"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
)

type jsonPbCodec struct {
	marshal   jsonpb.Marshaler
	unmarshal jsonpb.Unmarshaler
}

// 编码器的名称
func (self *jsonPbCodec) Name() string {
	return "jsonpb"
}

func (self *jsonPbCodec) MimeType() string {
	return "application/jsonpb"
}

// 将结构体编码为JSON的字节数组
func (self *jsonPbCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {

	if _, ok := msgObj.(proto.Message); !ok {
		return nil, cellnet.NewError("jsonpb marshal type assert failed")
	}

	jsonString, err := self.marshal.MarshalToString(msgObj.(proto.Message))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader([]byte(jsonString)), nil
}

// 将JSON的字节数组解码为结构体
func (self *jsonPbCodec) Decode(data interface{}, msgObj interface{}) error {

	var reader io.Reader
	switch v := data.(type) {
	case *http.Request:
		reader = v.Body
	case io.Reader:
		reader = v
	}

	if _, ok := msgObj.(proto.Message); !ok {
		return cellnet.NewError("jsonpb unmarshal type assert failed")
	}

	return self.unmarshal.Unmarshal(reader, msgObj.(proto.Message))
}

func init() {

	c := new(jsonPbCodec)
	c.marshal = jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "",
		OrigName:     false,
		AnyResolver:  nil,
	}

	c.unmarshal = jsonpb.Unmarshaler{
		AllowUnknownFields: false,
		AnyResolver:        nil,
	}

	// 注册编码器
	codec.RegisterCodec(c)
}
