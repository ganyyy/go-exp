package client

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"ganyyy.com/go-exp/rpc/grpc/logger"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
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
			DialOptions: []grpc.DialOption{
				grpc.WithStatsHandler(logger.NewHandle("etcd")),
			},
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
		opt = append(opt, clientv3.WithSerializable())
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

func Watch(ctx context.Context, key string, prefix bool) {
	var opts []clientv3.OpOption
	if prefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	opts = append(opts, clientv3.WithPrevKV())
	var watchChan = defaultClient.Watch(ctx, key, opts...)
	go func() {
		for msg := range watchChan {
			var sb strings.Builder
			sb.WriteString("Header:" + msg.Header.String())
			sb.WriteByte('\t')
			var allEvents []string
			for _, e := range msg.Events {
				allEvents = append(allEvents, fmt.Sprintf("{Type:%v, Kv:%v, PreKv:%v}", e.Type.String(), e.Kv.String(), e.PrevKv.String()))
			}
			sb.WriteString(strings.Join(allEvents, ","))
			// log.Printf("watch result:%+v", sb.String())
		}
	}()

}

func CAS(key, val, oldVal string) (bool, string, error) {
	var resp, err = defaultClient.Txn(context.Background()).If(
		// 比较版本是否为0
		clientv3.Compare(clientv3.Version(key), "=", 0),
	).Then(
		clientv3.OpPut(key, val),
	).Else(
		// 如果失败了, 则比较值是否为oldVal
		clientv3.OpTxn(
			[]clientv3.Cmp{
				clientv3.Compare(clientv3.Value(key), "=", oldVal),
			},
			[]clientv3.Op{
				clientv3.OpPut(key, val),
			},
			[]clientv3.Op{
				clientv3.OpGet(key),
			},
		),
	).Commit()
	if err != nil {
		return false, "", err
	}
	if !resp.Succeeded {
		// 失败了, 获取ETCD中的值
		switch op := resp.Responses[0].GetResponse().(type) {
		case *etcdserverpb.ResponseOp_ResponseTxn:
			if !op.ResponseTxn.Succeeded {
				return false, string(op.ResponseTxn.GetResponses()[0].GetResponseRange().GetKvs()[0].Value), nil
			} else {
				return true, val, nil
			}
		default:
			return false, "", fmt.Errorf("unknown response type:%T", op)
		}
	}
	// 成功了, 那么返回的值就是val
	return true, val, nil
}

func SetNX(key, val string) (bool, error) {
	var resp, err = defaultClient.Txn(context.Background()).If(
		clientv3.Compare(clientv3.Version(key), "=", 0),
	).Then(
		clientv3.OpPut(key, val),
	).Else(
		clientv3.OpGet(key),
	).Commit()
	if err != nil {
		return false, err
	}
	if !resp.Succeeded {
		return false, nil
	}
	return true, nil
}

func DistributedLock(ctx context.Context, key string) (context.CancelFunc, error) {
	session, err := concurrency.NewSession(defaultClient.Client)
	if err != nil {
		return nil, err
	}
	mutex := concurrency.NewMutex(session, defaultClient.Root+key)
	if err := mutex.Lock(ctx); err != nil {
		session.Close()
		return nil, err
	}
	return func() {
		mutex.Unlock(context.Background())
		session.Close()
	}, nil
}

func DistributedLockWithSession(ctx context.Context, session *concurrency.Session, key string) (context.CancelFunc, error) {
	mutex := concurrency.NewMutex(session, defaultClient.Root+key)
	if err := mutex.Lock(ctx); err != nil {
		return nil, err
	}
	return func() {
		mutex.Unlock(context.Background())
	}, nil
}
