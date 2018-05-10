package stormtf

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	logger       = &cLogger{}
	LoggingLevel = 2
)

func SetLoggingLevel(x int) {
	LoggingLevel = x
}

type cLogger struct{}

func Logger() *cLogger {
	return logger
}

func (*cLogger) Log(s string, a ...interface{}) {
	if LoggingLevel > 0 {
		v := fmt.Sprintf("[LOG]   %v", s)
		color.Blue(v, a...)
	}
}

func (*cLogger) Info(s string, a ...interface{}) {
	if LoggingLevel > 1 {
		v := fmt.Sprintf("[INFO]  %v", s)
		color.White(v, a...)
	}
}

func (*cLogger) Warn(s string, a ...interface{}) {
	v := fmt.Sprintf("[WARN]  %v", s)
	color.Red(v, a...)
}

func (*cLogger) Debug(s string, a ...interface{}) {
	if LoggingLevel > 1 {
		v := fmt.Sprintf("[DEBUG] %v", s)
		color.Yellow(v, a...)
	}
}

func (*cLogger) Error(s string, a ...interface{}) {
	v := fmt.Sprintf("[ERROR] %v", s)
	color.Red(v, a...)
}

func (*cLogger) Fatal(s string, a ...interface{}) {
	v := fmt.Sprintf("[FATAL] %v", s)
	color.Red(v, a...)
}
