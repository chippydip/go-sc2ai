package client

import (
	"fmt"
	"os"
	"runtime"
)

func recoverPanic() {
	if p := recover(); p != nil {
		ReportPanic(p)
	}
}

// ReportPanic ...
func ReportPanic(p interface{}) {
	fmt.Fprintln(os.Stderr, p)

	// Nicer format than what debug.PrintStack() gives us
	var pc [32]uintptr
	n := runtime.Callers(3, pc[:]) // skip the defer, this func, and runtime.Callers
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		fmt.Fprintf(os.Stderr, "%v:%v in %v\n", file, line, fn.Name())
	}
}
