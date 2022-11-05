package main

import (
	"log"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	allCode    = make(map[string]*protogen.Enum)
	allService = make(map[string]*protogen.Service)
)

func generateFile(f *protogen.File) {
	var prefix = *prefix
	for _, e := range f.Enums {
		if !strings.HasPrefix(e.GoIdent.GoName, prefix) {
			continue
		}
		allCode[e.GoIdent.GoName] = e
		for _, ele := range e.Values {
			log.Printf("%v", ele.GoIdent)
		}
	}

	for _, s := range f.Services {
		if !strings.HasPrefix(s.GoName, prefix) {
			continue
		}
		allService[s.GoName] = s
		log.Printf("service:%v", s.GoName)
		for _, r := range s.Methods {
			if r.Desc.IsStreamingClient() || r.Desc.IsStreamingServer() {
				continue
			}
			log.Printf("method: %v, Input:%v, Output:%v", r.GoName, r.Input.GoIdent, r.Output.GoIdent)
		}
	}
}
