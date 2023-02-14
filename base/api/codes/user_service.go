package codes

import (
	"github.com/jinvei/microservice/base/framework/codes"
)

const (
	ErrInternalXorm  = codes.Code(1000101)
	ErrInternalCache = codes.Code(1000102)
	ErrPassword      = codes.Code(1000103)
)

func init() {
	codes.Register(map[codes.Code]string{
		ErrInternalXorm:  "internal xorm error",
		ErrInternalCache: "internal cache error",
		ErrPassword:      "password error",
	})
}
