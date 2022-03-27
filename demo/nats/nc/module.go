package nc

import (
	"errors"
	"sync"
)

var (
	ErrServiceNotInit = errors.New("server not initial")
)

const (
	reqSubject  = "req"
	pushSubject = "psh"
)

func GenRequestSubject(subject string) string {
	return reqSubject + ":" + subject
}

func GenPushSubject(subject string) string {
	return pushSubject + ":" + subject
}

func GenClientSubject(subject string) (string, string) {
	return GenRequestSubject(subject), GenPushSubject(subject)
}

func GenServerSubject(subject string) (string, string) {
	return GenPushSubject(subject), GenRequestSubject(subject)
}

type PushCallBack func(msg NatsMessage)

type NatsService struct {
	reqSubject string
	group      string
	msgChan    chan NatsMessage
	once       sync.Once
	done       func()
}

func NewNatsServer(group, subject string) (*NatsService, error) {
	var req, receive = GenServerSubject(subject)
	return NewNatsService(group, req, receive)
}

func NewNatsClient(group, subject string) (*NatsService, error) {
	var req, receive = GenClientSubject(subject)
	return NewNatsService(group, req, receive)
}

func NewNatsService(group, req, receive string) (*NatsService, error) {
	var msgChan, cancel, err = defaultClient.Subcribe(group, receive)
	if err != nil {
		return nil, err
	}
	var service NatsService
	service.group = group
	service.reqSubject = req
	service.msgChan = msgChan
	service.done = cancel
	return &service, nil
}

func (ns *NatsService) Message() <-chan NatsMessage {
	return ns.msgChan
}

func (ns *NatsService) Request(meth string, req, rsp interface{}) error {
	return defaultClient.Request(ns.reqSubject, meth, req, rsp)
}

func (ns *NatsService) Push(meth string, data interface{}) error {
	return defaultClient.Push(ns.reqSubject, meth, data)
}

func (ns *NatsService) Stop() {
	if ns == nil {
		return
	}
	ns.once.Do(ns.done)
}

func Decode(data []byte, v interface{}) error {
	return defaultClient.codec.Decode(data, v)
}
