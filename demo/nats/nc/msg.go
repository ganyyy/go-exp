package nc

import "github.com/nats-io/nats.go"

// 服务, 方法, 数据

type NatsMessage struct {
	Method string
	Data   []byte

	conn *natsClient
	msg  *nats.Msg
}

func (m *NatsMessage) Init(msg *nats.Msg, conn *natsClient) {
	m.msg = msg
	m.conn = conn
}

func (m *NatsMessage) IsRequest() bool {
	if m.msg == nil || m.conn == nil {
		return false
	}
	return m.msg.Subject != "" && m.msg.Reply != ""
}

func (m *NatsMessage) Response(data interface{}) error {
	if !m.IsRequest() {
		return ErrNotValidRequest
	}
	if data == nil {
		return ErrNotValidResponse
	}

	return m.conn.Reply(m.msg.Reply, data)
}
