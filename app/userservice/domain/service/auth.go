package service

import (
	"context"
	"crypto/sha1"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinvei/microservice/app/userservice/domain"
	apicode "github.com/jinvei/microservice/base/api/codes"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/framework/cache"
	"github.com/jinvei/microservice/base/framework/codes"
	"github.com/jinvei/microservice/base/framework/configuration"
	confkeys "github.com/jinvei/microservice/base/framework/configuration/keys"
	"github.com/jinvei/microservice/base/framework/log"
	"github.com/jinvei/microservice/pkg/rand"
	"github.com/redis/go-redis/v9"
)

var flog = log.New()

type Auth struct {
	userRepo domain.IUserRepository
	cfg      svcConfig
	userSess *UserSession
}

type svcConfig struct {
	JwtSecret  string        `json:"jwtSecret"`
	TokenDura  time.Duration `json:"tokenDuration"`
	MaxSession int64         `json:"maxSession"`
}

func NewAuth(conf configuration.Configuration, userRepo domain.IUserRepository) domain.IAuthService {
	var (
		rdb *redis.Client
	)
	cache.RedisClient(conf)
	cfg := getSvcConfig(conf)

	return &Auth{
		userRepo: userRepo,
		cfg:      cfg,
		userSess: NewUserSession(rdb, cfg.TokenDura, cfg.MaxSession),
	}
}

type jwtToken struct {
	UserID  string `json:"id"`
	Session string `json:"session"`
	jwt.StandardClaims
}

func (a *Auth) SignInByEmail(ctx context.Context, in *app.SignInByEmailReq) (*app.SignInByEmailResp, error) {
	resp := app.SignInByEmailResp{}

	handleError := func(c codes.Code, err error) (*app.SignInByEmailResp, error) {
		if err != nil {
			flog.Error(err)
		}
		resp := app.SignInByEmailResp{}
		s := apicode.ErrInternalXorm.ToStatus()
		resp.Status = &s
		return &resp, nil
	}

	user, err := a.userRepo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return handleError(apicode.ErrInternalXorm, err)
	}

	// verify password
	password := sha1.Sum([]byte(in.Password + user.Salt))
	if string(password[:]) != in.Password {
		return handleError(apicode.ErrPassword, nil)
	}
	sessionKey := strconv.Itoa(int(time.Now().Unix()))

	uid := strconv.Itoa(int(user.ID))
	token := a.gennerateToken(uid, sessionKey)

	sid := rand.RandStringRunes(10)
	a.userSess.AddSession(ctx, uid, sid)

	resp.Session = []byte(token)

	return &resp, nil
}

func (a *Auth) SignOut(context.Context, *app.SignOutReq) (*app.SignOutResp, error) {
	return nil, nil
}

func (a *Auth) SignUpByEmail(context.Context, *app.SignUpByEmailReq) (*app.SignUpByEmailResp, error) {
	return nil, nil
}

func (a *Auth) SendEmailVerifyCode(context.Context, *app.SendEmailVerifyCodeReq) (*app.SendEmailVerifyCodeResp, error) {
	return nil, nil
}

// jwt
func (a *Auth) gennerateToken(userID, sessionKey string) string {
	claims := jwtToken{
		UserID:  userID,
		Session: sessionKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.cfg.TokenDura).Unix(),
			Issuer:    "auth-svc",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(a.cfg.JwtSecret)

	// Create token
	return tokenString
}

func getSvcConfig(conf configuration.Configuration) svcConfig {
	confkey := filepath.Join(confkeys.FwService, configuration.GetSystemID())

	// default value
	cfg := svcConfig{
		JwtSecret:  "secret",
		TokenDura:  time.Hour * 24 * 7, // 7 day
		MaxSession: 3,
	}
	if err := conf.GetJson(confkey, &cfg); err != nil {
		flog.Error(err)
		return cfg
	}
	return cfg
}
