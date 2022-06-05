package cmd

import (
	"log"
	"time"

	"ganyyy.com/go-exp/demo/kcp-go/common"
	"github.com/xtaci/kcp-go/v5"
)

func Client() {
	session, err := kcp.Dial(common.ADDR)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		var data = make([]byte, (i+1)*1000)
		for j := range data {
			data[j] = 'A'
		}
		n, err := session.Write(data)
		log.Printf("loop %v write data %v, error:%v", i, n, err)
		time.Sleep(time.Second)
	}
	session.Close()
}
