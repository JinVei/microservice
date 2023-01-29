package rpc

import (
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// tode
type ilog struct {
}

func newLog() *ilog {
	return &ilog{}
}
func (il *ilog) Log(level logging.Level, info string) {
	log.Println(info)
}

func (l *ilog) With(fields ...string) logging.Logger {
	return l
}
