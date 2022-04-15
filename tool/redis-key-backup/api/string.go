package api

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"redis-key-backup/config"
)

func init() {
	registerKeyOperation(config.TypeString, stringOperation{})
}

type stringOperation struct{}

func (s stringOperation) Dump(client *redis.Client, key string) (string, error) {
	var ret, err = client.Get(context.Background(), key).Result()
	return ret, err
}

func (s stringOperation) Restore(client *redis.Client, key, val string) error {
	var _, err = client.Set(context.Background(), key, val, 0).Result()
	return err
}
