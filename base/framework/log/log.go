package log

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

var (
	Default     *Log
	enableDebug = false
)

func init() {
	Default = New()
}

type Log struct {
	logger *zap.Logger
}

func NewLogger(scope string) (zapLogger *zap.Logger, err error) {
	zapconfig := zap.NewDevelopmentConfig()
	if !enableDebug {
		zapconfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zapLogger, err = zapconfig.Build()
	if err != nil {
		return nil, err
	}

	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(1))
	if scope != "" {
		zapLogger = zapLogger.Named(scope)
	}

	return zapLogger, nil
}

func New(scopes ...string) *Log {
	scope := ""
	if 1 <= len(scopes) {
		scope = "[" + scopes[0] + "] "
	}

	l, err := NewLogger(scope)
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

func (il *Log) Info(msg string, kv ...any) {
	il.logger.Info(msg, il.sweetenFields(kv)...)
}

func (il *Log) Debug(msg string, kv ...any) {
	il.logger.Debug(msg, il.sweetenFields(kv)...)
}
func (il *Log) Debugf(format string, v ...any) {
	il.logger.Debug(fmt.Sprintf(format, v...))
}

func (il *Log) Error(err error, msg string, kv ...any) {
	il.logger.Error(msg, il.sweetenFields(append([]any{}, zap.Any("err", err), kv))...)
}

func (il *Log) Errorf(err error, format string, v ...any) {
	il.logger.Error(fmt.Sprintf(format, v...), zap.Any("err", err))
}

func (il *Log) Warn(msg string, kv ...any) {
	il.logger.Warn(msg, il.sweetenFields(kv)...)
}

func (il *Log) Warnf(format string, v ...any) {
	il.logger.Warn(fmt.Sprintf(format, v...))
}

func EnableDebug(enable bool) {
	enableDebug = enable
}

func (il *Log) sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]zap.Field, 0, len(args))
	var invalid []any

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			il.logger.Info("_oddNumberErrMsg", zap.Any("ignored", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make([]any, 0, len(args)/2)
			}
			invalid = append(invalid, []any{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		il.logger.Info(_nonStringKeyErrMsg, zap.Any("invalid", invalid))
	}
	return fields
}
