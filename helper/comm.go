package helper

import "log"

func PanicIfErr(reason string, err error) {
	if err != nil {
		log.Panicf("[%v]:%v", reason, err)
	}
}
