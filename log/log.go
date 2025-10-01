//go:build !debug

package log

import "log"

func Debug(message ...any) {
	// 在非调试模式下，不输出调试信息
}

func Fatal(message ...any) {

}

func Info(message ...any) {
	log.Println(message)
}
