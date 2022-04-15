package api

import (
	"github.com/go-redis/redis/v8"

	"redis-key-backup/config"
)

func init() {
	registerKeyOperation(config.TypeZSet, zSetOperation{})
}

type zSetOperation struct{}

func (z zSetOperation) Dump(client *redis.Client, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (z zSetOperation) Restore(client *redis.Client, key, val string) error {
	//TODO implement me
	panic("implement me")
}
