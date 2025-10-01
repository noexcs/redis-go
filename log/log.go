//go:build !debug

package log

import "log"

func Debug(message ...any) {

}

func Fatal(message ...any) {

}

func Info(message ...any) {
	log.Println(message)
}
