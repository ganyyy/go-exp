package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
)

func main() {

	var cat = exec.Command("bash", "-c", "cat ./test.csv")
	var grep = exec.Command("bash", "-c", "awk -F ',' '{print $1,$2,$3}'")
	var buf = bytes.NewBuffer(nil)
	var output bytes.Buffer
	cat.Stdout = buf

	runWait := func(cmd *exec.Cmd) {
		log.Printf("%v start error:%v", cmd.Path, cmd.Run())
	}

	runWait(cat)

	log.Printf("data:\n%v", buf.String())
	grep.Stdin = buf
	grep.Stdout = &output
	runWait(grep)

	// log.Printf("result:\n%s", output.String())

	for line, err := output.ReadString('\n'); err == nil; line, err = output.ReadString('\n') {
		if err != nil {
			if err != io.EOF {
				log.Printf("find other error %v", err)
			}
			break
		}
		log.Printf("line:%v", line)
	}
}
