package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"p2p/cmd"
	"p2p/conn"
	"p2p/log"
	"p2p/msg"
	"p2p/tool"
	"strings"

	"github.com/urfave/cli"
)

const (
	Addr     = "addr"
	AddrFull = "addr, a"
)

func init() {
	cmd.Register(cli.Command{
		Name:   "client",
		Usage:  "p2p client",
		Action: Client,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     AddrFull,
				Usage:    "server addr, include HOST:PORT",
				Required: true,
			},
		},
	})
}

var (
	serverAddr conn.Addr
	client     *conn.Conn
)

func Client(ctx *cli.Context) {

	var addr = ctx.String(Addr)
	err := serverAddr.FromString(addr)
	if err != nil {
		log.Errorf("parse %v error %v", addr, err)
		return
	}
	// 创建客户端套接字
	client, err = conn.NewConn()
	if err != nil {
		log.Errorf("creat udp sock error:%v", err)
		return
	}
	var quit bool
	var wait tool.WaitGroup
	var done, cancel = context.WithCancel(context.Background())
	wait.Do(func() {
		defer func() {
			if !quit {
				client.SendMsg(msg.NewMsgNoData(msg.MsgLogout), serverAddr)
			}
		}()
		conn.ReceiveLoop(client, func(m msg.Msg, addr conn.Addr) {
			if addr == serverAddr {
				switch m.Type {
				case msg.MsgPunch:
					var peerAddr conn.Addr
					if err := peerAddr.FromBytes(m.Data); err != nil {
						log.Errorf("data %v parse addr error %v", m.Data, err)
						return
					}
					client.SendMsg(msg.NewMsgNoData(msg.MsgReply), peerAddr)
				case msg.MsgReply:
					log.Infof("SERVER MSG: %s", string(m.Data))
				}
			} else {
				switch m.Type {
				case msg.MsgText:
					log.Infof("Peer(%s): %s", addr.String(), m.Data)
				case msg.MsgReply:
					log.Warnf("Peer(%s) replied, you can talk now.", addr.String())
				case msg.MsgPing:
					//TODO 实现心跳机制
				case msg.MsgPunch:
					//TODO 实现去中心化(?)
				}
			}
		}, done)
	})

	wait.Do(func() {
		var scanner = bufio.NewScanner(os.Stdin)
		for {
			fmt.Print(">>> ")
			var cmd, host, data string
			if !scanner.Scan() {
				log.Errorf("scan console error:%v", err)
				break
			}
			var text = scanner.Text()
			var textSplit = strings.Split(text, " ")
			if len(textSplit) > 3 || len(textSplit) < 1 {
				log.Errorf("error input num")
				continue
			}

			n := len(textSplit)

			if n < 1 {
				log.Infof("must input valid cmd!")
				continue
			}

			if len(textSplit) > 0 {
				cmd = textSplit[0]
			}
			if len(textSplit) > 1 {
				host = textSplit[1]
			}
			if len(textSplit) > 2 {
				data = strings.Join(textSplit[2:], " ")
			}

			switch strings.ToLower(cmd) {
			case "list":
				client.SendMsg(msg.NewMsgNoData(msg.MsgList), serverAddr)
			case "login":
				client.SendMsg(msg.NewMsgNoData(msg.MsgLogin), serverAddr)
			case "logout":
				client.SendMsg(msg.NewMsgNoData(msg.MsgLogout), serverAddr)
			case "punch":
				if n != 2 {
					log.Errorf("punch need 2 param!")
					continue
				}
				var peerAddr conn.Addr
				if err := peerAddr.FromString(host); err != nil {
					log.Errorf("host %v error %v", host, err)
					continue
				}
				client.SendMsg(msg.NewMsg(msg.MsgPunch, peerAddr.Bytes()), serverAddr)
			case "send":
				if n != 3 {
					log.Errorf("send need 3 param!")
					continue
				}
				var peerAddr conn.Addr
				if err := peerAddr.FromString(host); err != nil {
					log.Errorf("host %v error %v", host, err)
					continue
				}
				client.SendMsg(msg.NewMsg(msg.MsgText, []byte(data)), peerAddr)
			case "quit":
				client.SendMsg(msg.NewMsgNoData(msg.MsgLogout), serverAddr)
				quit = true
				cancel()
				return
			case "help":
				fallthrough
			default:
				fmt.Printf(
					`
Usage:
login
  login to server so that other peer(s) can see you
logout
  logout from server
list
  list login peers
punch host:port
  punch a hole through UDP to [host:port]
  host:port must have been logged in to server
  Example:
  >>> punch 9.8.8.8:53
send host:port data
  send [data] to peer [host:port] through UDP protocol
  the other peer could receive your message if UDP hole punching succeed
  Example:
  >>> send 8.8.8.8:53 hello
help
  print this help message
quit
  logout and quit this program"
`,
				)
			}
		}
		cancel()
	})

	wait.Wait()
}
