package middleware

import (
	"net/http"
	"strconv"

	"github.com/jinvei/microservice/base/framework/log"

	"github.com/golang-jwt/jwt"
	"github.com/jinvei/microservice/app/userservice/domain/entity"
	"github.com/jinvei/microservice/app/userservice/domain/service/auth/config"
	"github.com/jinvei/microservice/app/userservice/wire"
	"github.com/jinvei/microservice/base/api/codes"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/labstack/echo/v4"
)

const (
	CONTEXT_JWT_KEY = "_auth_jwt"
)

var flog = log.Default

func NewAuthMiddleware() (echo.MiddlewareFunc, error) {
	conf, err := configuration.Default()
	if err != nil {
		return nil, err
	}
	conf.SetSystemID(strconv.Itoa(wire.SystemID))

	cfg, err := config.GetAuthConfig(conf)
	if err != nil {
		flog.Error(err, "config.GetAuthConfig")
		return nil, err
	}

	jwtsecret := cfg.JwtSecret
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				token = c.Request().URL.Query().Get("token")
			}

			if token == "" {
				s := codes.ErrUserAuthInvalidToken.ToStatus()
				return c.JSON(http.StatusUnauthorized, s)

			}

			ejwt := entity.Jwt{}
			// 解码 JWT
			jwt, err := jwt.ParseWithClaims(token, &ejwt, func(token *jwt.Token) (interface{}, error) {
				// 返回用于解码签名的密钥，这里可以根据实际情况进行设置
				return []byte(jwtsecret), nil
			})

			if err != nil {
				s := codes.ErrUserAuthInvalidToken.ToStatus()
				return c.JSON(http.StatusUnauthorized, s)
			}

			if !jwt.Valid {
				s := codes.ErrUserAuthInvalidJwt.ToStatus()
				return c.JSON(http.StatusUnauthorized, s)
			}
			// TODO: store jwt to context
			c.Set(CONTEXT_JWT_KEY, ejwt)
			flog.Debug("Auth token", "jwt", ejwt)
			// 调用下一个处理函数
			return next(c)
		}
	}, nil
}

func GetJwt(c echo.Context) entity.Jwt {
	return c.Get(CONTEXT_JWT_KEY).(entity.Jwt)
}
