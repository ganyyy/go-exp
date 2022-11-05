package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	version = "0.0.1"
)

var (
	showVersion = flag.Bool("version", false, "print version and exit")
	prefix      = flag.String("prefix", "Logic", "logicCommand and logic ")
)

// protoc --go_out=. -I. --go-handle_out=. --go-handle_opt=prefix=Logic  proto/*.proto

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println("versio:", version)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(f)
		}
		return nil
	})
}
