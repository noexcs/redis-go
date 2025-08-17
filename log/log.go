package log

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

var debugMode = flag.Bool("debug", false, "是否开启调试模式")

func init() {
	flag.Parse()
}

func WithLocation(message ...any) {
	if *debugMode {
		_, file, line, ok := runtime.Caller(1) // 获取调用者的信息
		var Location string
		if ok {
			Location = fmt.Sprintf("%s:%d -", file, line)
		} else {
			Location = "???:0 -"
		}
		log.Println(Location, message)
	}
}

func ForceWithLocation(message ...any) {
	_, file, line, ok := runtime.Caller(1)
	var Location string
	if ok {
		Location = fmt.Sprintf("%s:%d -", file, line)
	} else {
		Location = "???:0 -"
	}
	log.Println(Location, message)
}

func FatalWithLocation(message ...any) {
	WithLocation(message)
	os.Exit(1)
}
