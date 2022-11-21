package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"ganyyy.com/go-exp/rpc/grpc/proto"
)

var (
	port = flag.Int("port", 10086, "The server port")
)

func main() {
	flag.Parse()

	// logger.SetGRPCLogger()

	var keepParam = keepalive.ClientParameters{
		Time:                10 * time.Second, // 没有活跃的情况下, 最长多久发一次心跳包
		Timeout:             time.Second,      // 心跳包超过多久会认为链接已断开
		PermitWithoutStream: true,             // 没有激活的流, 是否允许发送心跳包?(可以理解为, 是针对RPC的链接保活. 这个应该开启更好点?)
	}

	var end, cancel = context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	var conn, err = grpc.DialContext(end,
		fmt.Sprintf(":%v", *port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), //
		grpc.WithKeepaliveParams(keepParam),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Second,      // 起始的重置延时
				Multiplier: 1.2,              // 下一次重试的递增系数
				Jitter:     0.2,              // 随机回退系数
				MaxDelay:   time.Second * 10, // 最大延时
			},
			MinConnectTimeout: time.Second * 10, // 多长时间连接不上表示需要进行重试
		}),
		// 连接创建时可以带上,  但是如果错误不为空, 只要保证地址是正确的, 也可以继续用
		// 阻塞式的情况下, 如果连接创建失败, 会返回一个空指针!
		// grpc.WithReturnConnectionError(), // 返回链接本身的错误, 而非context的错误, 阻塞式的等待连接成功
		// grpc.WithStatsHandler(logger.NewHandle("Client")),
	)
	if err != nil {
		log.Printf("dial error:%v", err)
	}
	defer func(conn *grpc.ClientConn) {
		closeError := conn.Close()
		if closeError != nil {
			log.Printf("close error:%v", err)
		}
	}(conn)

	var client = proto.NewGreeteClient(conn)

	dialContext, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	var doRPC = func() {
		rsp, err := client.SayHello(dialContext, &proto.HelloRequest{
			Name: "123131",
		})
		if err != nil {
			log.Printf("[SayHello]:%v, code:%v", err, status.Code(err))
		} else {
			log.Printf("[SayHello]:%v", rsp.GetMessage())
		}
	}
	var isClose = true

	var retryConn func(ctx context.Context)

	var streamContext context.Context
	var streamCancel context.CancelFunc

	retryConn = func(ctx context.Context) {

		select {
		case <-ctx.Done():
			log.Printf("[Stream] conn cancel")
			return
		default:
		}

		time.Sleep(time.Second * 3)
		var idBuf [12]byte
		// 流的重连必须要手动实现
		rand.Read(idBuf[:])
		id := hex.EncodeToString(idBuf[:])
		clientCtx := metadata.AppendToOutgoingContext(ctx, "id", id)
		log.Printf("[Stream:%v] start conn stream", id)
		stream, streamError := client.HelloStream(clientCtx)
		if streamError != nil {
			log.Printf("[Stream:%v] dial error:%v, code:%v", id, streamError, status.Code(streamError))
			go retryConn(ctx)
			return
		}
		go func() {
			for {
				resp, recvErr := stream.Recv()
				if recvErr != nil {
					if recvErr == io.EOF {
						// 客户端主动关闭时, 自己的发送端接收到的是 io.EOF, 此时不算错误
						log.Printf("[Stream:%v] close recv", id)
					} else {
						// 其他情况, 需要确定具体的错误类型
						log.Printf("[Stream:%v] recv error %v, code: %v", id, recvErr, status.Code(recvErr))
						// 重试应该在读的时候处理, 因为写的控制不确定, 但是读一定会立刻崩溃
						go retryConn(ctx)
					}
					return
				}
				log.Printf("[Stream:%v] resp:%v", id, resp)
			}
		}()
		var ticker = time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				var send = &proto.HelloRequest{Name: time.Now().String()}
				sendErr := stream.Send(send)
				if sendErr != nil {
					log.Printf("[Stream:%v] send error %v, code:%v", id, sendErr, status.Code(sendErr))
					return
				} else {
					log.Printf("[Stream:%v] send:%v", id, send)
				}
			case <-ctx.Done():
				// 客户端必须要执行CloseSend用来关闭流
				// 如果不关闭的话, 服务器会一直持有这个连接
				sendErr := stream.CloseSend()
				if sendErr != nil {
					log.Printf("[Stream:%v] send error %v", id, sendErr)
				} else {
					log.Printf("[Stream:%v] close send", id)
				}
				return
			}
		}
	}

	var cmd = make(chan string, 1)

	go func() {
		for {
			var s string
			n, err := fmt.Scanf("%s", &s)
			if err != nil || n != 1 {
				log.Printf("cmd n:%d, err:%v", n, err)
				continue
			}
			cmd <- s
		}
	}()

	var callTick = time.NewTicker(time.Second * 3)
	defer callTick.Stop()
	for {
		select {
		case <-callTick.C:
		case c, ok := <-cmd:
			if !ok {
				return
			}
			log.Printf("recv cmd:%s", c)
			switch c {
			case "c":
				doRPC()
			case "stream":
				if isClose {
					continue
				}
				isClose = true
				streamCancel()
			case "reconn":
				if !isClose {
					continue
				}
				isClose = false
				streamContext, streamCancel = context.WithCancel(context.Background())
				go retryConn(streamContext)
			case "quit":
				log.Printf("client quit")
				return
			}
		}
	}
}
