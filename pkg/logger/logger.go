package logger

import (
	"log"
	"os"
)

type LoggerType int64

const (
	Start LoggerType = iota
	DefaultLoggerType
	End // should always be last element in enum
)

func NewLogger(lType LoggerType) *log.Logger {
	switch lType {
	case DefaultLoggerType:
		fallthrough
	default:
		return log.New(os.Stdout, "", 0)
	}
}
