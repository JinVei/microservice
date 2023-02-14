package log

import (
	"log"
)

// TODO: framework log
type Log struct {
}

func New() *Log {
	return &Log{}
}

func (il *Log) Infof(format string, v ...any) {
	log.Printf(format, v...)
}

func (il *Log) Info(v ...any) {
	log.Println(v...)
}

func (il *Log) Debug(v ...any) {
	log.Println(v...)
}
func (il *Log) Debugf(format string, v ...any) {
	log.Printf(format, v...)
}

func (il *Log) Error(v ...any) {
	log.Println(v...)
}

func (il *Log) Errorf(format string, v ...any) {
	log.Printf(format, v...)
}

func (il *Log) Warn(v ...any) {
	log.Println(v...)
}

func (il *Log) Warnf(format string, v ...any) {
	log.Printf(format, v...)
}

// func (l *ilog) With(fields ...string) Log {
// 	return l
// }
