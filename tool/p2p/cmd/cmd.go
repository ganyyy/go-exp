package cmd

import (
	"os"
	"p2p/log"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func Register(command cli.Command) {
	log.Infof("register %v", command.Name)
	app.Commands = append(app.Commands, command)
}

func Run() {
	_ = app.Run(os.Args)
}
