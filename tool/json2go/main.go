package main

import (
	"json2go/parse"
	"os"

	"github.com/urfave/cli"
)

func main() {
	var app = cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "input,i",
			Usage:    "input json file",
			Required: true,
			Value:    "./json",
		},
		cli.StringFlag{
			Name:  "output,o",
			Usage: "output path",
			Value: "./data",
		},
		cli.StringFlag{
			Name:  "package,pkg",
			Usage: "go package name",
			Value: "data",
		},
		cli.BoolFlag{
			Name:  "number,n",
			Usage: "try parse number to int",
		},
	}
	app.Action = cli.ActionFunc(run)

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(ctx *cli.Context) error {

	var input = ctx.String("input")
	var gopkg = ctx.String("package")
	var output = ctx.String("output")

	var param = parse.ParseParam{
		InputFile:  input,
		OutputPath: output,
		GoPackage:  gopkg,
		UseNumber:  ctx.Bool("number"),
	}

	if err := param.InitOutput(); err != nil {
		return err
	}
	return param.Parse()
}
