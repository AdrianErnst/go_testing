package logger

import (
	"log"
	"os"
)

type LoggerType int64

const (
	DefaultLoggerType LoggerType = iota
)

func NewLogger(lType LoggerType) *log.Logger {
	switch lType {
	case DefaultLoggerType:
		fallthrough
	default:
		return log.New(os.Stdout, "", 0)
	}
}
