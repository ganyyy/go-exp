package api

import (
	"context"
	"encoding/json"
	"fmt"
	"redis-key-backup/config"

	"github.com/go-redis/redis/v8"
)

type listOperation struct{}

func init() {
	registerKeyOperation(config.TypeList, listOperation{})
}

func (l listOperation) Dump(client *redis.Client, key string) (string, error) {
	var ret, err = client.LRange(context.Background(), key, 0, -1).Result()
	if checkRedisError(err) != nil {
		return "", err
	}
	if len(ret) == 0 {
		return "", redis.Nil
	}
	bs, err := json.Marshal(listElements(ret))
	return BytesToString(bs), err
}

func (l listOperation) Restore(client *redis.Client, key string, val string) error {
	var elements listElements
	err := json.Unmarshal(StringToBytes(val), &elements)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return redis.Nil
	}
	var args = make([]interface{}, 0, len(elements))
	for _, ele := range elements {
		args = append(args, ele)
	}

	// 保持原有的顺序
	n, err := client.RPush(context.Background(), key, args...).Result()
	if err != nil {
		return err
	}
	if int(n) != len(elements) {
		return fmt.Errorf("elements len:%v, push num:%v", len(elements), n)
	}
	return nil
}
