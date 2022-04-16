package restore

import (
	"fmt"

	"github.com/urfave/cli"

	"redis-key-backup/api"
	"redis-key-backup/cmd"
	"redis-key-backup/config"
)

func init() {
	cmd.Register(cli.Command{
		Name:   "restore",
		Usage:  "restore redis key",
		Action: Restore,
		Flags:  append(config.RedisFlags, config.InputFlag...),
	})
}

func Restore(ctx *cli.Context) {
	config.Init(ctx)

	var fromFile = ctx.String(config.File)
	var fromTerminal = ctx.Bool(config.Input)
	var key = ctx.String(config.Key)

	var data string
	var err error
	if fromTerminal {
		_, err = fmt.Scan(&data)
		if err != nil {
			panic(fmt.Errorf("scan from input error:%w", err))
		}
	} else if fromFile != "" {
		data, err = api.ReadFromFile(fromFile)
		if err != nil {
			panic(err)
		}
	}
	if data == "" {
		panic("must input valid data!")
	}
	var saveData api.SaveStruct
	err = saveData.FromVal(data)
	if err != nil {
		panic(fmt.Sprintf("unpack error:%v", err))
	}
	err = saveData.DoRestore(key)
	if err != nil {
		panic(fmt.Sprintf("restore error:%v", err))
	}
	fmt.Printf("restore %v success!", key)
}
