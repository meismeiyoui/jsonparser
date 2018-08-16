package jsonparse

import (
	"log"
	"runtime/debug"
)

func Traceback(err error) {
	if EnableTraceBack {
		log.Println("[TRACEBACK]<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		log.Println(err)
		debug.PrintStack()
		log.Println("[TRACEBACK]>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	}
}
