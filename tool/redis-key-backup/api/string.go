package api

import (
	"context"

	"github.com/redis/go-redis/v9"

	"redis-key-backup/config"
)

func init() {
	registerKeyOperation(config.TypeString, stringOperation{})
}

type stringOperation struct{}

func (s stringOperation) Dump(client *redis.Client, key string) (string, error) {
	var ret, err = client.Get(context.Background(), key).Result()
	return ret, checkRedisError(err)
}

func (s stringOperation) Restore(client *redis.Client, key, val string) error {
	var _, err = client.Set(context.Background(), key, val, 0).Result()
	return err
}
