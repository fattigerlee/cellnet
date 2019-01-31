package gorillaws

import (
	"github.com/fattigerlee/cellnet"
	"github.com/fattigerlee/cellnet/msglog"
)

// 带有RPC和relay功能
type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteRecvLogger(log, "ws", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteSendLogger(log, "ws", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}
