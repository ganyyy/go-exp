package netpoll

import (
	"encoding/binary"

	"ganyyy.com/go-exp/net/msg"

	"github.com/cloudwego/netpoll"
)

func Decode(reader netpoll.Reader, m *msg.Message) (err error) {
	msgLen, err := reader.Next(4)
	if err != nil {
		return
	}
	bs, err := reader.ReadString(int(binary.BigEndian.Uint32(msgLen)))
	if err != nil {
		return
	}
	m.Message = string(bs)
	return reader.Release()
}

func Encode(write netpoll.Writer, msg *msg.Message) (err error) {
	// 头部空间
	var header, _ = write.Malloc(4)
	// 写入内容
	binary.BigEndian.PutUint32(header, uint32(len(msg.Message)))
	write.WriteString(msg.Message)
	// flush 刷新缓冲区, 写入到套接字中
	return write.Flush()
}
