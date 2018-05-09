package stormtf

import (
	"github.com/fatih/color"
)

type Logger struct{}

func (*Logger) Log(s string, a ...interface{}) {
	color.Blue(s, a...)
}

func (*Logger) Info(s string, a ...interface{}) {
	color.White(s, a...)
}

func (*Logger) Warn(s string, a ...interface{}) {
	color.Red(s, a...)
}

func (*Logger) Debug(s string, a ...interface{}) {
	color.Yellow(s, a...)
}

func (*Logger) Error(s string, a ...interface{}) {
	color.Red(s, a...)
}

var (
	logger = &Logger{}
)
