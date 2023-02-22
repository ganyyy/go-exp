package api

import (
	"encoding/json"

	"context"

	"github.com/go-redis/redis/v8"

	"redis-key-backup/config"
)

type hMapOperation struct{}

func init() {
	registerKeyOperation(config.TypeHMap, hMapOperation{})
}

func (h hMapOperation) Dump(client *redis.Client, key string) (string, error) {
	var ret, err = client.HGetAll(context.Background(), key).Result()
	if checkRedisError(err) != nil {
		return "", err
	}
	if len(ret) == 0 {
		return "", redis.Nil
	}
	var bs, _ = json.Marshal(ret)
	return BytesToString(bs), nil
}

func (h hMapOperation) Restore(client *redis.Client, key, val string) error {
	var hash = make(map[string]string)
	err := json.Unmarshal(StringToBytes(val), &hash)
	if err != nil {
		return err
	}
	if len(hash) == 0 {
		return redis.Nil
	}
	var args = make([]interface{}, 0, len(hash)*2)
	for k, v := range hash {
		args = append(args, k, v)
	}
	_, err = client.HMSet(context.Background(), key, args...).Result()
	return err
}
