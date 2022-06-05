package cmd

import (
	"io"
	"log"

	"ganyyy.com/go-exp/demo/kcp-go/common"
	"github.com/xtaci/kcp-go/v5"
)

func Server() {
	var listener, err = kcp.Listen(common.ADDR)
	if err != nil {
		panic(err)
	}

	for {
		ses, err := listener.Accept()
		if err != nil {
			break
		}
		go func(ses *kcp.UDPSession) {
			log.Println("accept session", ses.RemoteAddr())
			for {
				var buf [1024]byte
				var data, err = ses.Read(buf[:])
				if err != nil && err != io.EOF {
					log.Println("read error", err)
					return
				}
				log.Println("read data length:", data)
			}
		}(ses.(*kcp.UDPSession))
	}
}
