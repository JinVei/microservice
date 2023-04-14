package main

import (
	"github.com/jinvei/microservice/app/userservice/domain/service/auth/middleware"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/log"
	"github.com/jinvei/microservice/base/framework/web"
	"github.com/labstack/echo/v4"
)

const (
	systemID = 11001
)

func main() {
	conf := configuration.DefaultOrDie()
	authm, err := middleware.NewAuthMiddleware()
	if err != nil {
		panic(err)
	}

	web.App(conf, systemID, func(e *echo.Echo) {
		e.Use(authm)
		e.GET("/demo/test", func(c echo.Context) error {
			jwt := middleware.GetJwt(c)
			log.Default.Info("log jwt", "jwt", jwt)
			c.JSON(200, "Well Done!")
			return nil
		})
	})
}
