package codes

import (
	"github.com/jinvei/microservice/base/framework/codes"
)

const (
	StatusOK = codes.Code(0)
)

func init() {
	codes.Register(map[codes.Code]string{
		StatusOK: "ok",
	})
}
