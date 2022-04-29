package config

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func initClient() error {
	client = redis.NewClient(&redis.Options{
		Addr:         redisConfig.Host,
		Password:     redisConfig.Auth,
		DB:           redisConfig.DB,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		PoolSize:     1,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	return nil
}

func GetClient() *redis.Client {
	return client
}
