package codes

import (
	"github.com/jinvei/microservice/base/framework/codes"
)

const (
	ErrInternalXorm       = codes.Code(1000101)
	ErrInternalCache      = codes.Code(1000102)
	ErrPassword           = codes.Code(1000103)
	ErrParseJwt           = codes.Code(1000104)
	ErrJwtInvalid         = codes.Code(1000105)
	ErrInvalidEmail       = codes.Code(1000106)
	ErrInvalidPassword    = codes.Code(1000107)
	ErrInvalidUsername    = codes.Code(1000108)
	ErrInternalVerifyCode = codes.Code(1000109)
	ErrInvalidVerifyCode  = codes.Code(1000110)
	ErrCreateUser         = codes.Code(1000111)
	ErrSendEmail          = codes.Code(1000112)
	ErrVerifyCodeTooMany  = codes.Code(1000113)
	ErrUnknownInternal    = codes.Code(1000114)
)

func init() {
	codes.Register(map[codes.Code]string{
		ErrInternalXorm:       "internal xorm error",
		ErrInternalCache:      "internal cache error",
		ErrPassword:           "password error",
		ErrParseJwt:           "parse JWT error",
		ErrJwtInvalid:         "jwt invalid",
		ErrInvalidEmail:       "Invalid Email",
		ErrInvalidPassword:    "Invalid Password",
		ErrInvalidUsername:    "Invalid Username",
		ErrInternalVerifyCode: "Internal Error in verify code",
		ErrInvalidVerifyCode:  "Invalid Verify Code",
		ErrCreateUser:         "Create User Error",
		ErrSendEmail:          "Send Email Error",
		ErrUnknownInternal:    "Unknown Internal Error",
	})
}
