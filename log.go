package hclient

import (
	"log"
)

type (
	Logger struct {
	}
)

func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Println(format, v)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	log.Println(format, v)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	log.Println(format, v)
}
