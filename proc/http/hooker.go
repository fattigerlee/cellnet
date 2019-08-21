package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
)

type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteRecvLogger(log, "http", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteSendLogger(log, "http", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}
