package test

import (
	"flag"
	"redis-key-backup/config"

	"github.com/urfave/cli"
)

func init() {
	initTestClient()
}

func initTestClient() {

	var flagSet = flag.NewFlagSet("test", 0)
	flagSet.String(config.Host, "", "host")
	flagSet.Int(config.DB, 0, "db")
	flagSet.String(config.Auth, "", "auth")
	flagSet.Parse([]string{
		"-host", "127.0.0.1:6379",
		"-db", "15",
	})

	var ctx = cli.NewContext(nil, flagSet, nil)
	config.Init(ctx)
}
