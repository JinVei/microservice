package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/jinvei/microservice/app/userservice/wire"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/datasource"
	"github.com/jinvei/microservice/base/framework/rpc"
	"github.com/jinvei/microservice/base/framework/web"
	"google.golang.org/grpc"
)

func main() {
	var waitSrv sync.WaitGroup
	waitSrv.Add(2)

	conf := configuration.DefaultOrDie()
	conf.SetSystemID(strconv.Itoa(wire.SystemID))

	db := datasource.New(conf, wire.SystemID)
	userrepo := wire.InitUserRepository(db.Orm())

	go func() {
		defer waitSrv.Done()

		err := rpc.Serve(conf, wire.SystemID, func(srv *grpc.Server) {
			authSvc := wire.InitAuthServer(conf, userrepo)

			app.RegisterAuthServiceServer(srv, authSvc)
		})

		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer waitSrv.Done()
		authctrl := wire.InitAuthWeb(conf, userrepo)

		web.App(nil, wire.SystemID, authctrl.InitRoute)
	}()

	waitSrv.Wait()
}
