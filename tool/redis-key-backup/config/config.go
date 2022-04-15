package config

import (
	"github.com/urfave/cli"
)

type RedisConnect struct {
	Host string
	Auth string
	DB   int
}

var RedisFlags = []cli.Flag{
	cli.StringFlag{
		Name:  HostFull,
		Usage: "redis addr",
		Value: "localhost:6379",
	},
	cli.StringFlag{
		Name:  AuthFull,
		Usage: "redis password",
		Value: "",
	},
	cli.IntFlag{
		Name:  DBFull,
		Usage: "redis db",
		Value: 0,
	},
	cli.StringFlag{
		Name:     KeyFull,
		Usage:    "redis key",
		Required: true,
	},
	cli.StringFlag{
		Name:  FileFull,
		Usage: "out file",
	},
}

var OutputFlag = []cli.Flag{
	cli.BoolFlag{
		Name:  OutputFull,
		Usage: "output to terminal",
	},
}

var InputFlag = []cli.Flag{
	cli.BoolFlag{
		Name:  InputFull,
		Usage: "input from terminal",
	},
}

var Config RedisConnect

func Init(cmd *cli.Context) {
	if err := parse(cmd); err != nil {
		panic(err)
	}
	if err := initClient(); err != nil {
		panic(err)
	}
}

func parse(cmd *cli.Context) error {
	var host = cmd.String(Host)
	if host == "" {
		host = "localhost:6379"
	}
	var auth = cmd.String(Auth)
	var db = cmd.Int(DB)

	Config.Host = host
	Config.Auth = auth
	Config.DB = db

	return nil
}
