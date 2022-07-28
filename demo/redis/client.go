package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Addr string
	DB   int
	Auth string
)

func init() {
	log.SetPrefix("[Redis]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var (
	client *redis.Client
	ctx    = context.Background()
)

func InitClient() {
	client = redis.NewClient(&redis.Options{
		Addr:         Addr,
		DB:           DB,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		PoolSize:     2,
	})

	var ret, err = client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Printf("ping result:%v", ret)
}
