// Generated by github.com/davyxu/cellnet/protoc-gen-msg
// DO NOT EDIT!
// Source: pb.proto

package test

import (
	"github.com/fattigerlee/cellnet"
	"github.com/fattigerlee/cellnet/codec"
	_ "github.com/fattigerlee/cellnet/codec/gogopb"
	"reflect"
)

func init() {

	// pb.proto
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*ContentACK)(nil)).Elem(),
		ID:    60952,
	})
}
