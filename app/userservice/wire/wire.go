//go:build wireinject
// +build wireinject

package wire

import (
	"strconv"

	"github.com/google/wire"
	"github.com/jinvei/microservice/app/userservice/app/rpc"
	"github.com/jinvei/microservice/app/userservice/domain"
	"github.com/jinvei/microservice/app/userservice/domain/repository"
	"github.com/jinvei/microservice/app/userservice/domain/service"
	"github.com/jinvei/microservice/base/framework/configuration"
	"xorm.io/xorm"
)

const (
	SystemID = 10001
)

func init() {
	configuration.SetSystemID(strconv.Itoa(SystemID))
}

func InitUserRepository(engine *xorm.Engine) domain.IUserRepository {
	panic(wire.Build(repository.NewUserRepository))
}

func InitAuthServer(configuration.Configuration, domain.IUserRepository) *rpc.AuthServer {
	panic(wire.Build(service.NewAuth, rpc.NewAuthServer))
	return &rpc.AuthServer{}
}
