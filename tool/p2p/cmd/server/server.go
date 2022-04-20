package server

import (
	"context"
	"fmt"
	"p2p/cmd"
	"p2p/conn"
	"p2p/log"
	"p2p/msg"
	"p2p/tool"
	"strings"
	"syscall"

	"github.com/urfave/cli"
)

const (
	Addr     = "addr"
	AddrFull = "addr, a"
)

func init() {
	cmd.Register(cli.Command{
		Name:   "server",
		Usage:  "p2p server",
		Action: Server,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     AddrFull,
				Usage:    "server addr, include HOST:PORT",
				Required: true,
			},
		},
	})
}

func Server(ctx *cli.Context) {
	var listenAddr conn.Addr
	var addr = ctx.String(Addr)
	log.Infof("cli addr %v", addr)
	if err := listenAddr.FromString(addr); err != nil {
		log.Errorf("parse %v error %v", addr, err)
		return
	}

	var listenConn, err = conn.NewConn()
	if err != nil {
		log.Errorf("create listen socket error:%v", err)
		return
	}

	if err := syscall.Bind(listenConn.Fd(), listenAddr.ToSysInet4Addr()); err != nil {
		log.Errorf("bind error:%v", err)
		return
	}
	var wait tool.WaitGroup
	var done, cancel = context.WithCancel(context.Background())
	defer cancel()
	wait.Do(func() {
		conn.ReceiveLoop(listenConn, func(m msg.Msg, addr conn.Addr) {
			switch m.Type {
			case msg.MsgLogin:
				conn.Add(addr)
				log.Infof("peer(%s) login!", addr.String())
				listenConn.SendMsg(msg.NewMsg(msg.MsgReply, "Login Success!"), addr)
			case msg.MsgLogout:
				conn.Del(addr)
				log.Infof("peer(%s) logout!", addr.String())
				listenConn.SendMsg(msg.NewMsg(msg.MsgReply, "Logout Success!"), addr)
			case msg.MsgList:
				var data strings.Builder
				for _, ca := range conn.List() {
					if ca == addr {
						data.WriteString("(you);")
						continue
					}
					data.WriteString(ca.String())
					data.WriteString(";")
				}
				listenConn.SendMsg(msg.NewMsg(msg.MsgReply, data.String()), addr)
			case msg.MsgPunch:
				var peerAddr conn.Addr
				if err := peerAddr.FromBytes(m.Data); err != nil {
					listenConn.SendMsg(msg.NewMsg(msg.MsgReply, err.Error()), addr)
					return
				}
				if !conn.Has(peerAddr) {
					listenConn.SendMsg(msg.NewMsg(msg.MsgReply,
						fmt.Sprintf("cannot found peer(%s)", peerAddr.String())),
						addr)
					return
				}
				log.Infof("peer(%s) punch to peer(%s)", addr.String(), peerAddr.String())
				listenConn.SendMsg(msg.NewMsg(msg.MsgPunch, addr.String()), peerAddr)
				listenConn.SendMsg(msg.NewMsg(msg.MsgReply, "punch request send"), addr)
			default:
				log.Infof("peer(%s) send %s", addr.String(), m.Data)
			}
		}, done)
	})

	wait.Wait()

}
