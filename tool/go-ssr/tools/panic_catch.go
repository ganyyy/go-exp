package tools

import "log"

func PanicCatch() {
	if err := recover(); err != nil {
		log.Panicln(err)
	}
}
