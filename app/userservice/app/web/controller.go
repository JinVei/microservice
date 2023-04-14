package web

import (
	"github.com/jinvei/microservice/app/userservice/domain"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/framework/log"
	"github.com/labstack/echo/v4"
)

var flog = log.Default

type AuthController struct {
	pb domain.IAuthService
}

func NewAuthService(pb domain.IAuthService) *AuthController {
	return &AuthController{
		pb: pb,
	}
}

func (a *AuthController) InitRoute(e *echo.Echo) {
	e.POST("/v1/signin", a.Signin)
	e.POST("/v1/signup", a.Signup)
	e.POST("/v1/signout", a.Signout)
}

func (a *AuthController) Signin(c echo.Context) error {
	req := &app.SignInByEmailReq{}
	err := c.Bind(req)
	if err != nil {
		flog.Error(err, "c.Bind(req)")
		return err
	}

	resp, err := a.pb.SignInByEmail(c.Request().Context(), req)
	if err != nil {
		flog.Error(err, "a.pb.SignInByEmail()")
		return err
	}

	return c.JSON(200, resp)
}

func (a *AuthController) Signout(c echo.Context) error {
	req := &app.SignOutReq{}
	err := c.Bind(req)
	if err != nil {
		flog.Error(err, "c.Bind(req)")
		return err
	}

	resp, err := a.pb.SignOut(c.Request().Context(), req)
	if err != nil {
		flog.Error(err, "a.pb.SignOut()")
		return err
	}
	return c.JSON(200, resp)
}

func (a *AuthController) Signup(c echo.Context) error {
	req := &app.SignUpByEmailReq{}
	err := c.Bind(req)
	if err != nil {
		flog.Error(err, "c.Bind(req)")
		return err
	}

	resp, err := a.pb.SignUpByEmail(c.Request().Context(), req)
	if err != nil {
		flog.Error(err, "a.pb.SignOut()")
		return err
	}
	return c.JSON(200, resp)
}
