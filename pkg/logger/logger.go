package logger

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if msg == "" {
		fmt.Println("")
		return
	}

	c := color.New(color.FgHiCyan)
	if len(args) > 0 {
		c.Printf(msg+"\n", args...)
	} else {
		c.Println(msg)
	}
}

func (l *Logger) Error(err error) {
	c := color.New(color.FgHiRed)
	c.Println(fmt.Sprintf("%#v", err))
}

func (l *Logger) Instructions(msg string, args ...interface{}) {
	white := color.New(color.FgHiWhite)
	white.Println("")
	white.Println(fmt.Sprintf(msg, args...))
}
