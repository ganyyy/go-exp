package dump

import (
	"fmt"

	"github.com/urfave/cli"

	"redis-key-backup/api"
	"redis-key-backup/cmd"
	"redis-key-backup/config"
)

func init() {
	cmd.Register(cli.Command{
		Name:   "dump",
		Usage:  "dump redis key",
		Action: Dump,
		Flags:  append(config.RedisFlags, config.OutputFlag...),
	})
}

func Dump(ctx *cli.Context) {
	config.Init(ctx)

	var toFile = ctx.String(config.File)
	var toTerminal = ctx.Bool(config.Output)
	var key = ctx.String(config.Key)

	// 既没有输出地址, 也不输出到控制台
	if toFile == "" && !toTerminal {
		panic("invalid out path")
	}

	var err error
	var saveStruct api.SaveStruct
	err = saveStruct.DoDump(key)
	if err != nil {
		panic(fmt.Sprintf("do dump %v error:%v", key, err))
	}

	var data = saveStruct.String()
	if toFile != "" {
		err = api.ExportToFile(toFile, data)
		if err != nil {
			panic(fmt.Sprintf("write to %v error %v", toFile, err))
		}
	}
	if toTerminal {
		fmt.Print(data)
	}
}
