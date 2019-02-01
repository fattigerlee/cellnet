package gorillaws

import (
	"encoding/binary"
	"github.com/fattigerlee/cellnet"
	"github.com/fattigerlee/cellnet/codec"
	"github.com/fattigerlee/cellnet/util"
	"github.com/gorilla/websocket"
)

const (
	MsgIDSize = 12 // uint16
)

type WSMessageTransmitter struct {
}

func (WSMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil, nil
	}

	var messageType int
	var raw []byte
	messageType, raw, err = conn.ReadMessage()

	if err != nil {
		return
	}

	if len(raw) < MsgIDSize {
		return nil, util.ErrMinPacket
	}

	switch messageType {
	case websocket.BinaryMessage:
		// 处理keep alive包
		cmd := binary.BigEndian.Uint32(raw)
		seq := binary.BigEndian.Uint32(raw[8:])
		if cmd == 1 && seq == 999 {
			conn.WriteMessage(websocket.BinaryMessage, raw)
			return
		}

		if len(raw) < MsgIDSize+2 {
			return nil, util.ErrMinPacket
		}

		// 处理其他消息包
		msgID := binary.BigEndian.Uint16(raw[12:])
		msgData := raw[MsgIDSize+2:]

		msg, _, err = codec.DecodeMessage(int(msgID), msgData)
	}

	return
}

func (WSMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	var (
		msgData []byte
		msgID   int
	)

	switch m := msg.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error
		var meta *cellnet.MessageMeta

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(msg, nil)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	pkt := make([]byte, MsgIDSize+2+len(msgData))
	binary.BigEndian.PutUint32(pkt, uint32(0))
	binary.BigEndian.PutUint32(pkt[4:], uint32(int32(len(msgData))+MsgIDSize+2))
	binary.BigEndian.PutUint32(pkt[8:], uint32(0))
	binary.BigEndian.PutUint16(pkt[12:], uint16(msgID))
	copy(pkt[MsgIDSize+2:], msgData)

	conn.WriteMessage(websocket.BinaryMessage, pkt)

	return nil
}
