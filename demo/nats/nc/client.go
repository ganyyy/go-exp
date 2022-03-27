package nc

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	ErrNotInit          = errors.New("nats client not init!")
	ErrNotValidRequest  = errors.New("cannot response invalid response")
	ErrNotValidResponse = errors.New("not valid response")
)

type NatsClientConfig struct {
	Addr  string
	Port  int
	Codec Codec
}

func (nc NatsClientConfig) URL() string {
	return fmt.Sprintf("nats://%s:%d", nc.Addr, nc.Port)
}

var (
	defaultClient *natsClient
	closeOnce     sync.Once
	initOnce      sync.Once
	initError     error
)

type natsClient struct {
	config NatsClientConfig
	conn   *nats.EncodedConn
	quit   chan struct{}
	wait   sync.WaitGroup
	codec  Codec
}

func Init(config NatsClientConfig) error {
	initOnce.Do(func() {
		nc, err := nats.Connect(config.URL())

		if err != nil {
			initError = err
			return
		}
		encodeConn, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
		if err != nil {
			initError = err
			return
		}
		if config.Codec == nil {
			config.Codec = JsonCodec
		}
		defaultClient = &natsClient{
			config: config,
			conn:   encodeConn,
			codec:  config.Codec,
			quit:   make(chan struct{}),
		}
		return
	})
	return initError
}

func Stop() {
	closeOnce.Do(func() {
		var nc = defaultClient
		if nc == nil {
			return
		}
		if nc.conn == nil {
			return
		}
		nc.conn.Close()
		close(nc.quit)
		nc.wait.Wait()
	})
}

//Subcribe
//group: 不同的服务器实例
//subject: 订阅的服务
func (nc *natsClient) Subcribe(group, subject string) (chan NatsMessage, func(), error) {
	if nc == nil || nc.conn == nil {
		return nil, nil, ErrNotInit
	}
	var msgChan = make(chan NatsMessage, 1024)
	var subscription, err = nc.conn.QueueSubscribe(subject, group, func(msg *nats.Msg) {
		var m NatsMessage
		_ = nc.conn.Enc.Decode(subject, msg.Data, &m)
		// log.Printf("[INF] Receive msg:%+v, %v, %v,err:%v", string(msg.Data), msg.Subject, msg.Reply, err)
		m.Init(msg, defaultClient)
		select {
		case msgChan <- m:
		default:
		}
	})
	if err != nil {
		return nil, nil, err
	}
	nc.wait.Add(1)
	return msgChan, func() {
		_ = subscription.Unsubscribe()
		nc.wait.Done()
	}, nil
}

//Request 同步请求
func (nc *natsClient) Request(subject, method string, req, rsp interface{}) error {
	if nc.conn == nil {
		return ErrNotInit
	}
	var request, reply NatsMessage
	request.Method = method
	data, err := nc.codec.Encode(req)
	if err != nil {
		return err
	}
	request.Data = data
	err = nc.conn.Request(subject, request, &reply, time.Second*3)
	if err != nil {
		return err
	}
	return nc.codec.Decode(reply.Data, rsp)
}

func (nc *natsClient) Reply(reply string, v interface{}) error {
	return nc.publish(reply, "", v)
}

func (nc *natsClient) Push(subject, method string, v interface{}) error {
	return nc.publish(subject, method, v)
}

func (nc *natsClient) publish(subject, method string, v interface{}) error {
	if nc.conn == nil {
		return ErrNotInit
	}
	data, err := nc.codec.Encode(v)
	if err != nil {
		return err
	}
	var msg NatsMessage
	msg.Method = method
	msg.Data = data
	return nc.conn.Publish(subject, msg)
}
