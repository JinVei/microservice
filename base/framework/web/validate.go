package web

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var formValidator = &echoValidator{validator: validator.New()}

type echoValidator struct {
	validator *validator.Validate
}

func (e echoValidator) Validate(i interface{}) error {
	err := e.validator.Struct(i)
	if err == nil {
		return nil
	}
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}
