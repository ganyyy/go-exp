package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Os      string `toml:"os"`
	Version int    `toml:"version"`
}

type TomlServer struct {
	Addr string `toml:"addr"`
}

type TomlStudent struct {
	Name string `toml:"name"`
	Age  int    `toml:"age"`
}

func main() {
	var bs, _ = os.ReadFile("./conf/config.toml")

	var cfg TomlConfig
	state, _ := toml.Decode(string(bs), &cfg)

	log.Printf("%+v", state)

	var server struct {
		Server TomlServer `toml:"server"`
	}
	_ = toml.Unmarshal(bs, &server)

	var info struct {
		Info struct {
			Arr1 []int    `toml:"arr1"`
			Arr2 []string `toml:"arr2"`
		} `toml:"info"`
	}
	_ = toml.Unmarshal(bs, &info)

	var stu struct {
		Stu map[string]TomlStudent `toml:"students"`
	}
	_ = toml.Unmarshal(bs, &stu)

	log.Printf("%+v, %+v, %+v, %+v", cfg, server, info, stu)
}
