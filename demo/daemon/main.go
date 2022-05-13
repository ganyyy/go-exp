//go:build linux

package main

import (
	"log"
	"os"
	"syscall"
	"time"
)

func main() {

	var pid, _, err = syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != syscall.Errno(0) {
		log.Panicf("syscall.SYS_FORK error %v", err)
	}

	// 退出父进程
	if pid > 0 {
		log.Printf("Create Daemon process %v", pid)
		time.Sleep(time.Minute)
		os.Exit(0)
	}

	// 设置sid
	syscall.Setsid()

	// 转移根目录
	syscall.Chdir("/")

	// 设置umask
	syscall.Umask(0)

	// 关闭文件描述符
	// 就这么关吧..
	for i := 0; i < 3; i++ {
		syscall.Close(i)
	}

	for {
		var fd, err = syscall.Open("/tmp/go_daemon.log", syscall.O_CREAT|syscall.O_WRONLY|syscall.O_APPEND, 0600)
		if err != nil {
			log.Panicf("open daemon log error %v", err)
		}
		syscall.Write(fd, []byte(time.Now().String()))
		syscall.Close(fd)
		// log.Printf("112345")
		time.Sleep(time.Second * 10)
	}
}
