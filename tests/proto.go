package tests

import (
	"fmt"
	"github.com/fattigerlee/cellnet"
	"github.com/fattigerlee/cellnet/codec"
	_ "github.com/fattigerlee/cellnet/codec/binary"
	"github.com/fattigerlee/cellnet/util"
	"reflect"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(util.StringHash("tests.TestEchoACK")),
	})
}
