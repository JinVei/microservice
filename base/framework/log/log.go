package log

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// TODO: framework log
type Log struct {
	logger logr.Logger
}

func NewLogger() (logr.Logger, error) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return logr.Logger{}, err
	}
	return zapr.NewLogger(zapLogger), nil
}

func New() *Log {
	l, err := NewLogger()
	if err != nil {
		panic(err)
	}
	return &Log{
		logger: l,
	}
}

func (il *Log) Infof(format string, v ...any) {
	il.logger.Info(fmt.Sprintf(format, v...))
}

func (il *Log) Info(msg string, v ...any) {
	il.logger.Info(msg, v...)
}

func (il *Log) Debug(msg string, v ...any) {
	il.logger.V(5).Info(msg, v...)
}
func (il *Log) Debugf(format string, v ...any) {
	il.logger.V(5).Info(fmt.Sprintf(format, v...))
}

func (il *Log) Error(err error, msg string, v ...any) {
	il.logger.Error(err, msg, v...)
}

func (il *Log) Errorf(err error, format string, v ...any) {
	il.logger.Error(err, fmt.Sprintf(format, v...))
}

func (il *Log) Warn(msg string, v ...any) {
	il.logger.V(3).Info(msg, v...)
}

func (il *Log) Warnf(format string, v ...any) {
	il.logger.V(3).Info(fmt.Sprintf(format, v...))
}

// func (l *ilog) With(fields ...string) Log {
// 	return l
// }
