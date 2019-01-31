package gorillaws

import (
	"github.com/fattigerlee/cellnet"
	"github.com/fattigerlee/cellnet/proc"
)

func init() {

	proc.RegisterProcessor("gorillaws.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(WSMessageTransmitter))
		bundle.SetHooker(new(MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))

	})
}
