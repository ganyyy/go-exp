package main

import (
	"redis-key-backup/cmd"
	_ "redis-key-backup/cmd/dump"
	_ "redis-key-backup/cmd/restore"
)

func main() {
	cmd.Run()
}
