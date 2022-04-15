package cmd

import (
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func Register(command cli.Command) {
	app.Commands = append(app.Commands, command)
}

func Run() {
	_ = app.Run(os.Args)
}
