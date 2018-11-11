package client

import (
	"log"
	"runtime"
)

func recoverPanic() {
	if p := recover(); p != nil {
		ReportPanic(p)
	}
}

// ReportPanic ...
func ReportPanic(p interface{}) {
	log.Print(p)

	// Nicer format than what debug.PrintStack() gives us
	var pc [32]uintptr
	n := runtime.Callers(3, pc[:]) // skip the defer, this func, and runtime.Callers
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		log.Printf("%v:%v in %v", file, line, fn.Name())
	}
}
