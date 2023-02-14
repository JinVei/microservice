package rpc

import (
	"context"

	"github.com/jinvei/microservice/app/userservice/domain"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
)

type AuthServer struct {
	svc domain.IAuthService
	app.UnimplementedAuthServiceServer
}

func NewAuthServer(authsvc domain.IAuthService) *AuthServer {
	return &AuthServer{
		svc: authsvc,
	}
}

func (a *AuthServer) SignInByEmail(ctx context.Context, in *app.SignInByEmailReq) (*app.SignInByEmailResp, error) {
	return a.svc.SignInByEmail(ctx, in)
}

// func (a *AuthService) SignOut(context.Context, *app.SignOutReq) (*app.SignOutResp, error) {

// }

// func (a *AuthService) SignUpByEmail(context.Context, *app.SignUpByEmailReq) (*app.SignUpByEmailResp, error) {

// }

// func (a *AuthService) SendEmailVerifyCode(context.Context, *app.SendEmailVerifyCodeReq) (*app.SendEmailVerifyCodeResp, error) {

// }
