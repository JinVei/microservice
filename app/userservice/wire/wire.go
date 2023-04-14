//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jinvei/microservice/app/userservice/app/rpc"
	"github.com/jinvei/microservice/app/userservice/app/web"
	"github.com/jinvei/microservice/app/userservice/domain"
	"github.com/jinvei/microservice/app/userservice/domain/repository"
	"github.com/jinvei/microservice/app/userservice/domain/service/auth"
	"github.com/jinvei/microservice/base/framework/configuration"
	"xorm.io/xorm"
)

const (
	SystemID = 10001
)

func InitUserRepository(engine *xorm.Engine) domain.IUserRepository {
	panic(wire.Build(repository.NewUserRepository))
}

func InitAuthServer(configuration.Configuration, domain.IUserRepository) *rpc.AuthServer {
	panic(wire.Build(auth.NewAuth, rpc.NewAuthServer))
	return &rpc.AuthServer{}
}

func InitAuthWeb(configuration.Configuration, domain.IUserRepository) *web.AuthService {
	panic(wire.Build(auth.NewAuth, web.NewAuthService))
	return &web.AuthService{}
}
