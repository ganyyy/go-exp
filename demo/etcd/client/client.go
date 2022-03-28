package client

import (
	"context"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfig struct {
	Endpoints []string
	Username  string
	Password  string
	Root      string
}

var (
	initOnce      sync.Once
	initError     error
	closeOnce     sync.Once
	defaultClient struct {
		Root string
		*clientv3.Client
	}
)

func Init(config *EtcdConfig) error {

	initOnce.Do(func() {
		defaultClient.Client, initError = clientv3.New(clientv3.Config{
			Endpoints: config.Endpoints,
			Username:  config.Username,
			Password:  config.Password,
		})
		defaultClient.Root = config.Root
	})

	return initError
}

func Stop() {
	closeOnce.Do(func() {
		if initError != nil {
			return
		}
		if defaultClient.Client == nil {
			return
		}
		_ = defaultClient.Close()
	})
}

func Put(key, val string) error {
	_, err := defaultClient.Put(context.Background(), key, val)
	return err
}

type KV struct {
	Key string
	Val string
}

func Get(key string, prefix bool) ([]KV, error) {
	var opt []clientv3.OpOption
	if prefix {
		opt = append(opt, clientv3.WithPrefix())
	}
	resp, err := defaultClient.Get(context.Background(), key, opt...)
	if err != nil {
		return nil, err
	}
	var ret = make([]KV, 0, resp.Count)
	for i := 0; i < int(resp.Count); i++ {
		ret = append(ret, KV{
			Key: string(resp.Kvs[i].Key),
			Val: string(resp.Kvs[i].Value),
		})
	}
	return ret, nil
}
