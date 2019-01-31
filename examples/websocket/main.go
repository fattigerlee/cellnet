package main

import (
	"github.com/davyxu/golog"
	"github.com/fattigerlee/cellnet"
	"github.com/fattigerlee/cellnet/peer"
	"github.com/fattigerlee/cellnet/proc"

	"fmt"
	"github.com/fattigerlee/cellnet/codec"
	_ "github.com/fattigerlee/cellnet/codec/json"
	_ "github.com/fattigerlee/cellnet/peer/gorillaws"
	_ "github.com/fattigerlee/cellnet/proc/gorillaws"
	"reflect"
)

var log = golog.New("websocket_server")

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

// 将消息注册到系统
func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    1234,
	})
}

// 运行服务器, 在浏览器(Chrome)中打开index.html, F12打开调试窗口->Console标签 查看命令行输出
// 注意, 如果http代理/VPN在运行时可能会导致无法连接, 请关闭
func main() {

	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("gorillaws.Acceptor", "server", "http://127.0.0.1:18802~18803/echo", queue)

	proc.BindProcessorHandler(p, "gorillaws.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {

		case *cellnet.SessionAccepted:
			log.Debugln("server accepted")
			// 有连接断开
		case *cellnet.SessionClosed:
			log.Debugln("session closed: ", ev.Session().ID())
		case *TestEchoACK:
			log.Debugf("recv: %+v", msg)

			ev.Session().Send(&TestEchoACK{
				Msg:   "roger",
				Value: 1234,
			})
		}
	})

	// 开始侦听
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	queue.Wait()

}
