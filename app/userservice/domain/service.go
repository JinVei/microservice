package domain

import (
	"context"

	"github.com/jinvei/microservice/base/api/proto/v1/app"
)

type IAuthService interface {
	SignInByEmail(ctx context.Context, in *app.SignInByEmailReq) (*app.SignInByEmailResp, error)
	SignOut(context.Context, *app.SignOutReq) (*app.SignOutResp, error)
	SignUpByEmail(context.Context, *app.SignUpByEmailReq) (*app.SignUpByEmailResp, error)
	SendEmailVerifyCode(context.Context, *app.SendEmailVerifyCodeReq) (*app.SendEmailVerifyCodeResp, error)
}
