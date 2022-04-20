package conn

import (
	"context"
	"p2p/log"
	"p2p/msg"
	"syscall"
)

type OnMessage func(msg msg.Msg, addr Addr)

func ReceiveLoop(conn *Conn, onMessage OnMessage, ctx context.Context) {

	type Msg struct {
		msg  msg.Msg
		addr Addr
	}

	var msgChan = make(chan Msg, 100)
	go func() {
		for {
			var m, addr, err = conn.ReadMsg()
			if err == nil {
				msgChan <- Msg{
					msg:  m,
					addr: addr,
				}
				continue
			}
			switch err {
			case msg.ErrDataLength, msg.ErrMagicNum, msg.ErrHeadLength:
				log.Errorf("receive %v msg error:%v", addr.String(), err)
			case syscall.EINTR:
				continue
			default:
				log.Errorf("receive %v net error:%v", addr.String(), err)
				ctx.Done()
			}
		}
	}()
	for {
		select {
		case m := <-msgChan:
			onMessage(m.msg, m.addr)
		case <-ctx.Done():
			return
		}
	}
}
