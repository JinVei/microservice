// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/jinvei/microservice/app/userservice/app/rpc"
	"github.com/jinvei/microservice/app/userservice/domain"
	"github.com/jinvei/microservice/app/userservice/domain/repository"
	"github.com/jinvei/microservice/app/userservice/domain/service"
	"github.com/jinvei/microservice/base/framework/configuration"
	"strconv"
	"xorm.io/xorm"
)

// Injectors from wire.go:

func InitUserRepository(engine *xorm.Engine) domain.IUserRepository {
	iUserRepository := repository.NewUserRepository(engine)
	return iUserRepository
}

func InitAuthServer(configurationConfiguration configuration.Configuration, iUserRepository domain.IUserRepository) *rpc.AuthServer {
	iAuthService := service.NewAuth(configurationConfiguration, iUserRepository)
	authServer := rpc.NewAuthServer(iAuthService)
	return authServer
}

// wire.go:

const (
	SystemID = 10001
)

func init() {
	configuration.SetSystemID(strconv.Itoa(SystemID))
}
