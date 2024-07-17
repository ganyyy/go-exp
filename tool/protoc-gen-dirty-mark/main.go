package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	version = "0.0.1"
)

var (
	showVersion = flag.Bool("version", false, "print version and exit")

	pkgName = flag.String("pkg", "data", "package name")
	pkgPath = flag.String("pkgpath", "data", "package path")

	pbAlias   = flag.String("pb", "pb", "protobuf package alias")
	pbPath    = flag.String("pbpath", "pb", "protobuf package path")
	metaAlias = flag.String("meta", "meta", "meta package alias")
	metaPath  = flag.String("metapath", "meta", "meta package path")
)

func main() {

	flag.Parse()

	if *showVersion {
		println("versio:", version)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range p.Files {
			println(file.GeneratedFilenamePrefix, file.GoPackageName, file.GoImportPath)
			for _, msg := range file.Messages {
				for _, field := range msg.Fields {
					if field.Desc.IsList() {
						println("list", field.GoName, field.Desc.Kind().GoString())
					} else if field.Desc.IsMap() {
						println("map", field.GoName, field.Desc.MapKey().Kind().GoString(), field.Desc.MapValue().Kind().GoString())
					} else {
						println("field", field.GoName, field.Desc.Kind().GoString())
					}
				}
			}
			f := p.NewGeneratedFile("test", file.GoImportPath)
			f.P("package ", "world")
		}
		return nil
	})

}
