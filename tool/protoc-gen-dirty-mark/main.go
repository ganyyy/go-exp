package main

import (
	"flag"
	"io"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	version = "0.0.1"
)

var (
	showVersion = flag.Bool("version", false, "print version and exit")
	pbAlias     = flag.String("pb", "pb", "protobuf package alias")
	metaAlias   = flag.String("meta", "meta", "meta package alias")
)

func main() {

	flag.Parse()

	if *showVersion {
		println("versio:", version)
		return
	}

	var opts protogen.Options
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var req pluginpb.CodeGeneratorRequest

	if err := proto.Unmarshal(in, &req); err != nil {
		panic(err)
	}

	_ = opts
}
