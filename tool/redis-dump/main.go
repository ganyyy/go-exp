package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	Host    = flag.String("host", "localhost", "地址")
	Port    = flag.Int("port", 6379, "端口")
	Auth    = flag.String("auth", "", "密码")
	DB      = flag.Int("db", 0, "数据库地址")
	Key     = flag.String("key", "", "备份的地址")
	Restore = flag.Bool("restore", false, "备份or写入")
)

func main() {
	flag.Parse()

	var client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", *Host, *Port),
		DB:       *DB,
		Password: *Auth,
	})

	var ctx = context.Background()

	if !*Restore {
		ret, err := client.HGetAll(ctx, *Key).Result()
		if err != nil && err != redis.Nil {
			panic(err)
		}

		data, _ := json.Marshal(ret)

		_ = os.WriteFile("data.json", data, 0644)
	} else {
		info, err := os.ReadFile("data.json")
		if err != nil {
			panic(err)
		}
		var m map[string]interface{}

		err = json.Unmarshal(info, &m)
		if err != nil {
			panic(err)
		}

		var pipeline = client.Pipeline()
		defer pipeline.Close()
		var args = make([]interface{}, 0, len(m)<<1)
		for key, val := range m {
			args = append(args, key, val)
		}
		pipeline.HMSet(ctx, *Key, args...)

		_, err = pipeline.Exec(ctx)
		if err != nil {
			panic(err)
		}

	}
}
