package helper

import (
	"log"
	"syscall"
)

func FilterSyscallError(reason string, err syscall.Errno) {
	if err != syscall.Errno(0) {
		PanicIfErr(reason, err)
	}
}

func PanicIfErr(reason string, err error) {
	if err != nil {
		log.Panicf("[%v]:%v", reason, err)
	}
}
