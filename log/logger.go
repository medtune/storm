package log

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	LoggingLevel = 2
)

func SetLoggingLevel(x int) {
	LoggingLevel = x
}

func Log(s string, a ...interface{}) {
	if LoggingLevel > 0 {
		v := fmt.Sprintf("[LOG]   %v", s)
		color.Blue(v, a...)
	}
}

func Info(s string, a ...interface{}) {
	if LoggingLevel > 1 {
		v := fmt.Sprintf("[INFO]  %v", s)
		color.White(v, a...)
	}
}

func Warn(s string, a ...interface{}) {
	v := fmt.Sprintf("[WARN]  %v", s)
	color.Red(v, a...)
}

func Debug(s string, a ...interface{}) {
	if LoggingLevel > 1 {
		v := fmt.Sprintf("[DEBUG] %v", s)
		color.Yellow(v, a...)
	}
}

func Error(s string, a ...interface{}) {
	v := fmt.Sprintf("[ERROR] %v", s)
	color.Red(v, a...)
}

func Fatal(s string, a ...interface{}) {
	v := fmt.Sprintf("[FATAL] %v", s)
	color.Red(v, a...)
}
