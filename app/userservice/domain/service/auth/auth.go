package auth

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinvei/microservice/app/userservice/domain"
	"github.com/jinvei/microservice/app/userservice/domain/entity"
	"github.com/jinvei/microservice/app/userservice/domain/service/auth/config"
	apicode "github.com/jinvei/microservice/base/api/codes"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/api/proto/v1/dto"
	"github.com/jinvei/microservice/base/framework/cache"
	"github.com/jinvei/microservice/base/framework/codes"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/log"
	"github.com/jinvei/microservice/pkg/emailsender"
	"github.com/jinvei/microservice/pkg/rand"
	"golang.org/x/crypto/bcrypt"
)

var flog = log.New()

const (
	VerifyCodeFormat = "{\"verify_code\":\"%s\"}"
)

type Auth struct {
	userRepo  domain.IUserRepository
	cfg       config.AuthConfig
	userSess  *UserSession
	esender   emailsender.EmailSender
	svcID     string
	tokenDura time.Duration
}

// type svcConfig struct {
// 	JwtSecret  string `json:"jwtSecret"`
// 	TokenDura  string `json:"tokenDuration"`
// 	MaxSession int64  `json:"maxSession"`

// 	evconf emailsender.Config
// }

func NewAuth(conf configuration.Configuration, userRepo domain.IUserRepository) domain.IAuthService {
	rdb := cache.RedisClient(conf)
	cfg, err := config.GetAuthConfig(conf)
	if err != nil {
		panic(err)
	}

	esender, err := emailsender.New(cfg.ESconf)
	if err != nil {
		panic(err)
	}
	tdura := time.Hour * 24 * 7 // 7h
	d, err := time.ParseDuration(cfg.TokenDura)
	if err != nil {
		flog.Warn("ParseDuration()", "err", err)
	} else {
		tdura = d
	}

	return &Auth{
		userRepo:  userRepo,
		cfg:       cfg,
		userSess:  NewUserSession(rdb, tdura, cfg.MaxSession),
		esender:   esender,
		svcID:     conf.GetSystemID(),
		tokenDura: tdura,
	}
}

func (a *Auth) SignInByEmail(ctx context.Context, in *app.SignInByEmailReq) (*app.SignInByEmailResp, error) {
	resp := app.SignInByEmailResp{}

	handleError := func(c codes.Code, err error) (*app.SignInByEmailResp, error) {
		if err != nil {
			flog.Error(err, "handleError")
		}
		resp := app.SignInByEmailResp{}
		resp.Status = c.ToStatus()
		return &resp, nil
	}

	user, err := a.userRepo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return handleError(apicode.ErrInternalXorm, err)
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password+user.Salt)); err != nil {
		flog.Debug("CompareHashAndPassword() err", "err", err)
		return handleError(apicode.ErrPassword, nil)
	}

	sid := rand.RandStringRunes(10)
	uid := strconv.Itoa(int(user.Id))
	token, err := a.gennerateToken(uid, sid)
	if err != nil {
		handleError(apicode.ErrUnknownInternal, err)
	}

	if err := a.userSess.AddSession(ctx, uid, sid); err != nil {
		handleError(apicode.ErrUnknownInternal, err)
	}

	resp.Token = []byte(token)

	return &resp, nil
}

func (a *Auth) SignOut(ctx context.Context, req *app.SignOutReq) (*app.SignOutResp, error) {
	// parse jwt
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.cfg.JwtSecret), nil
	})
	if err != nil {
		flog.Error(err, "jwt.Parse()", "JwtSecret", a.cfg.JwtSecret, "req.Token", req.Token)
		s := apicode.ErrParseJwt.ToStatus()
		return &app.SignOutResp{Status: s}, nil
	}

	claims, ok := token.Claims.(*entity.Jwt)
	if ok && token.Valid {
		s := apicode.ErrParseJwt.ToStatus()
		return &app.SignOutResp{Status: s}, nil
	}

	if status := a.userSess.DelUserSession(ctx, claims.Uid, claims.Sid); status != nil {
		return &app.SignOutResp{
			Status: status,
		}, nil
	}
	return &app.SignOutResp{}, nil
}

func (a *Auth) SignUpByEmail(ctx context.Context, in *app.SignUpByEmailReq) (*app.SignUpByEmailResp, error) {
	// verify email
	handleBadStatus := func(s *dto.Status) (*app.SignUpByEmailResp, error) {
		return &app.SignUpByEmailResp{
			Status: s,
		}, nil
	}

	if !isValidEmail(in.Email) {
		return handleBadStatus(apicode.ErrInvalidEmail.ToStatus())
	}
	// verify password
	if !isValidPassword(in.Password) {
		return handleBadStatus(apicode.ErrInvalidPassword.ToStatus())
	}
	// verify username
	if !isValidUsername(in.Username) {
		return handleBadStatus(apicode.ErrInvalidUsername.ToStatus())
	}
	// check verify code
	verifyCode, ok := a.userSess.GetVerifyCode(ctx, in.Email)
	if !ok {
		return handleBadStatus(apicode.ErrInvalidVerifyCode.ToStatus())
	}

	if verifyCode != in.VerifyCode {
		return handleBadStatus(apicode.ErrInvalidVerifyCode.ToStatus())
	}

	salt := rand.RandStringRunes(4)
	password, err := bcrypt.GenerateFromPassword([]byte(in.Password+salt), bcrypt.MinCost)
	if err != nil {
		flog.Error(err, "GenerateFromPassword()")
		return handleBadStatus(apicode.ErrInvalidPassword.ToStatus())
	}

	createby, err := strconv.Atoi(a.svcID)
	if err != nil {
		flog.Error(err, "strconv.Atoi", "a.svcID", a.svcID)
	}
	timeNow := uint64(time.Now().Unix())
	user := entity.Users{
		Username:       in.Username,
		Password:       string(password[:]),
		Email:          in.Email,
		Salt:           salt,
		Status:         entity.UserStatusNormal,
		CreateBy:       uint64(createby),
		CreateTime:     timeNow,
		CreatedAt:      time.Now(),
		LastModifyBy:   uint64(createby),
		LastModifyTime: timeNow,
		UpdatedAt:      time.Now(),
	}

	err = a.userRepo.CreateUser(ctx, &user)
	if err != nil {
		flog.Error(err, "CreateUser()", "user", user)
		return handleBadStatus(apicode.ErrCreateUser.ToStatus())
	}

	ret, err := a.SignInByEmail(ctx, &app.SignInByEmailReq{
		Email:    in.Email,
		Password: in.Password,
	})

	if err != nil {
		flog.Error(err, "SignInByEmail", "Email", in.Email)
		return nil, err
	}

	if ret.Status != nil {
		return &app.SignUpByEmailResp{
			Status: ret.Status,
		}, nil
	}

	return &app.SignUpByEmailResp{
		Token: ret.Token,
	}, nil
}

func (a *Auth) SendEmailVerifyCode(ctx context.Context, in *app.SendEmailVerifyCodeReq) (*app.SendEmailVerifyCodeResp, error) {
	if !isValidEmail(in.Email) {
		return &app.SendEmailVerifyCodeResp{Status: apicode.ErrInvalidEmail.ToStatus()}, nil
	}
	// has already send verify code to email before
	code, ok := a.userSess.GetVerifyCode(ctx, in.Email)
	if ok {
		return &app.SendEmailVerifyCodeResp{Status: apicode.ErrVerifyCodeTooMany.ToStatus()}, nil //ok
	}

	code = rand.RandStrNumber(4)
	param := fmt.Sprintf(VerifyCodeFormat, code)
	if err := a.esender.Send(ctx, in.Email, param); err != nil {
		flog.Error(err, "esender.Send", "email", in.Email, "param", param)
		return &app.SendEmailVerifyCodeResp{Status: apicode.ErrSendEmail.ToStatus()}, nil
	}

	a.userSess.PutVerifyCode(ctx, in.Email, code)
	return &app.SendEmailVerifyCodeResp{}, nil
}

//func (a *Auth) ResetPasswordByEmail(){}

// jwt
func (a *Auth) gennerateToken(userID, sid string) (string, error) {
	claims := entity.Jwt{
		Uid: userID,
		Sid: sid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.tokenDura).Unix(),
			Issuer:    "auth-svc",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.cfg.JwtSecret))

	// Create token
	return tokenString, err
}

// func getSvcConfig(conf configuration.Configuration) (svcConfig, error) {
// 	// default value
// 	cfg := svcConfig{
// 		JwtSecret:  "secret",
// 		TokenDura:  "168h", // 7 day
// 		MaxSession: 3,
// 	}
// 	if err := conf.GetSvcJson(configuration.GetSystemID(), "", &cfg); err != nil {
// 		flog.Error(err, "conf.GetSvcJson()")
// 		return cfg, err
// 	}

// 	evconf := emailsender.Config{}
// 	if err := conf.GetSvcJson(configuration.GetSystemID(), "/emailverify", &evconf); err != nil {
// 		flog.Error(err, "conf.GetSvcJson", "path", "/emailverify")
// 		return cfg, err
// 	}
// 	cfg.evconf = evconf

// 	return cfg, nil
// }

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]{4,20}$`)
	return re.MatchString(password)
}

func isValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]{4,20}$`)
	return re.MatchString(username)
}
