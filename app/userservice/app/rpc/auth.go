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

func (a *AuthServer) SignOut(ctx context.Context, in *app.SignOutReq) (*app.SignOutResp, error) {
	return a.svc.SignOut(ctx, in)
}

func (a *AuthServer) SignUpByEmail(ctx context.Context, in *app.SignUpByEmailReq) (*app.SignUpByEmailResp, error) {
	return a.svc.SignUpByEmail(ctx, in)
}

func (a *AuthServer) SendEmailVerifyCode(ctx context.Context, in *app.SendEmailVerifyCodeReq) (*app.SendEmailVerifyCodeResp, error) {
	return a.svc.SendEmailVerifyCode(ctx, in)
}
