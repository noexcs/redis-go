//go:build debug
// +build debug

package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func Debug(message ...any) {
	_, file, line, ok := runtime.Caller(1) // 获取调用者的信息
	var Location string
	if ok {
		Location = fmt.Sprintf("%s:%d -", file, line)
	} else {
		Location = "???:0 -"
	}
	log.Println(Location, message)
}

func Fatal(message ...any) {
	Debug(message)
	os.Exit(1)
}

func Info(message ...any) {
	Debug(message)
}
