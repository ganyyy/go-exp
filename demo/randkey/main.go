package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli"
)

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
		),
	)
}

func main() {

	var redisServer string
	var redisPort int
	var redisAuth string
	var redisDB int
	var concurrency int
	var keyNum int
	var keyExpire int
	var setStep int

	var app cli.App
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "server, s",
			Usage:       "server address",
			Destination: &redisServer,
			Value:       "127.0.0.1",
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "port, p",
			Usage:       "server port",
			Destination: &redisPort,
			Value:       6379,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "auth, a",
			Usage:       "server auth",
			Value:       "",
			Destination: &redisAuth,
		},
		&cli.IntFlag{
			Name:        "db, d",
			Usage:       "server db",
			Destination: &redisDB,
			Value:       5,
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "concurrency, c",
			Usage:       "concurrency",
			Destination: &concurrency,
			Value:       10,
		},
		&cli.IntFlag{
			Name:        "keynum, k",
			Usage:       "key number",
			Destination: &keyNum,
			Value:       10_000_000,
		},
		&cli.IntFlag{
			Name:        "expire, e",
			Usage:       "key expire",
			Destination: &keyExpire,
			Value:       60,
		},
		&cli.IntFlag{
			Name:        "setstep, ss",
			Usage:       "set step",
			Destination: &setStep,
			Value:       10000,
		},
	}

	app.Action = func(c *cli.Context) error {

		// 检查 addr+port 是否可用

		var tcpAddr *net.TCPAddr
		var err error

		{
			// 解析 addr+port, 检查是否可用
			addr := fmt.Sprintf("%s:%d", redisServer, redisPort)
			tcpAddr, err = net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				slog.Error("net.ResolveTCPAddr", slog.String("error", err.Error()))
				return err
			}
		}

		var client = redis.NewClient(&redis.Options{
			Addr:         tcpAddr.String(),
			DB:           redisDB,
			Password:     redisAuth,
			DialTimeout:  time.Second * 5,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			PoolSize:     10,
			MinIdleConns: 1,
		})

		{
			// 判断是否需要 auth
			err := client.Ping(context.Background()).Err()
			if err != nil {
				slog.Error("client.Ping", slog.String("error", err.Error()))
				return err
			}
		}

		{
			// 起多个 goroutine, 每个 goroutine 负责写入一部分 key
			var wg sync.WaitGroup
			wg.Add(concurrency)

			for i := 0; i < concurrency; i++ {
				go func(i int) {
					defer wg.Done()
					var ctx = context.Background()

					slog.Info("pipeline.Exec start", slog.Int("i", i))
					var pipeline = client.Pipeline()
					defer pipeline.Close()
					doExec := func(idx int) {
						if pipeline.Len() == 0 {
							slog.Info("pipeline.Len() == 0", slog.Int("i", i), slog.Int("idx", idx))
							return
						}
						_, err := pipeline.Exec(ctx)
						if err != nil {
							slog.Error("pipeline.Exec", slog.String("error", err.Error()))
							return
						}
						slog.Info("pipeline.Exec success", slog.Int("i", i), slog.Int("idx", idx))
						// 重置 pipeline
						pipeline.Discard()
					}

					for j := 0; j < keyNum/concurrency; j++ {
						var key = fmt.Sprintf("key-%d-%d", i, j)
						err := pipeline.SetEX(ctx, key, "value", time.Duration(keyExpire)*time.Second).Err()
						if err != nil {
							slog.Error("pipeline.Set", slog.String("error", err.Error()))
							return
						}
						if j > 0 && j%setStep == 0 {
							doExec(j)
						}
					}
					doExec(keyNum / concurrency)
					slog.Info("pipeline.Exec end", slog.Int("i", i))
				}(i)
			}

			wg.Wait()
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error("app.Run", slog.String("error", err.Error()))
	}
}
