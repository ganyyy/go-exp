package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	version = "0.0.1"
)

var (
	commonPkgPtr = flag.String(
		commonPkgParam,
		"protoc-gen-go-handle/common",
		"common package name",
	)
	moduleBasePathPtr = flag.String(
		moduleBasePathParam,
		"module",
		"module base path",
	)
)

const (
	commonPkgParam      = "common_pkg"
	moduleBasePathParam = "module_base_path"
)

var (
	commonPkgPath  string
	moduleBasePath string
)

func main() {

	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println("version:", version)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		commonPkgPath = *commonPkgPtr
		moduleBasePath = *moduleBasePathPtr

		log.Println(commonPkgParam, ":", commonPkgPath)
		log.Println(moduleBasePathParam, ":", moduleBasePath)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}
