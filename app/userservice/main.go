package main

import (
	"log"
	"sync"

	"github.com/jinvei/microservice/app/userservice/wire"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/datasource"
	"github.com/jinvei/microservice/base/framework/rpc"
	"github.com/jinvei/microservice/base/framework/web"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func main() {
	var waitSrv sync.WaitGroup
	waitSrv.Add(2)

	conf := configuration.DefaultOrDie()

	go func() {
		defer waitSrv.Done()

		err := rpc.Serve(nil, wire.SystemID, func(srv *grpc.Server) {
			db := datasource.New(conf, wire.SystemID)

			userrepo := wire.InitUserRepository(db.Orm())
			authSvc := wire.InitAuthServer(conf, userrepo)

			app.RegisterAuthServiceServer(srv, authSvc)
		})

		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer waitSrv.Done()
		web.App(nil, wire.SystemID, func(e *echo.Echo) {
			e.GET("/demo/test", func(c echo.Context) error {
				c.JSON(200, "Well Done!")
				return nil
			})
		})
	}()

	waitSrv.Wait()
}
