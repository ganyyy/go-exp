package main

import (
	"flag"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	version = "0.0.1"
)

var (
	showVersion = flag.Bool("version", false, "print version and exit")

	pkgName   = flag.String("pkg", "dirtygen", "package name")
	pkgPath   = flag.String("pkgpath", "dirtygen", "package path")
	pbAlias   = flag.String("pb", "pb", "protobuf package alias")
	pbPath    = flag.String("pbpath", "pb", "protobuf package path")
	metaAlias = flag.String("meta", "meta", "meta package alias")
	metaPath  = flag.String("metapath", "meta", "meta package path")
)

func main() {

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(p *protogen.Plugin) error {

		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		if *showVersion {
			println("versio:", version)
			return nil
		}

		var importInfo = ImportInfo{
			PBAlias:   *pbAlias,
			MetaAlias: *metaAlias,
			Imports: map[string]string{
				*pbPath:   *pbAlias,
				*metaPath: *metaAlias,
			},
		}

		println("pkgName:", *pkgName)
		println("pkgPath:", *pkgPath)
		println("pbAlias:", *pbAlias)
		println("pbPath:", *pbPath)
		println("metaAlias:", *metaAlias)
		println("metaPath:", *metaPath)

		for _, file := range p.Files {
			var outputFile File
			outputFile.Structs = make(map[string]*Struct)
			outputFile.ImportInfo = &importInfo
			outputFile.Name = filepath.Join(*pkgPath, filepath.Base(file.GeneratedFilenamePrefix)+".dirty.go")
			outputFile.Package = *pkgName
			// println(file.GeneratedFilenamePrefix, file.GoPackageName, file.GoImportPath)
			for _, msg := range file.Messages {
				var s Struct
				s.Name = msg.GoIdent.GoName
				for _, field := range msg.Fields {
					parseField(field, &s)
				}
				outputFile.Structs[s.Name] = &s
			}
			content, err := outputFile.Render()
			if err != nil {
				return err
			}
			f := p.NewGeneratedFile(outputFile.Name, protogen.GoImportPath(*pkgPath))
			f.Write(content)
		}
		return nil
	})

}
