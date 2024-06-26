package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// var debugMode = flag.Bool("debug", false, "是否开启调试模式")
var debugMode = true

func init() {

}

func WithLocation(message ...any) {

	if !debugMode {
		log.Println(message)
	} else {
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

func FatalWithLocation(message ...any) {
	WithLocation(message)
	os.Exit(1)
}
