package rpc

import (
	"github.com/jinvei/microservice/app/userservice/domain"
)

type AuthServer struct {
	domain.IAuthService
	//app.UnimplementedAuthServiceServer

}

func NewAuthServer(authsvc domain.IAuthService) *AuthServer {
	return &AuthServer{
		IAuthService: authsvc,
	}
}
