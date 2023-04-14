package config

import (
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/log"
	"github.com/jinvei/microservice/pkg/emailsender"
)

var flog = log.Default

type AuthConfig struct {
	JwtSecret  string `json:"jwtSecret"`
	TokenDura  string `json:"tokenDuration"`
	MaxSession int64  `json:"maxSession"`

	ESconf emailsender.Config
}

func GetAuthConfig(conf configuration.Configuration) (AuthConfig, error) {
	// default value
	cfg := AuthConfig{
		JwtSecret:  "secret",
		TokenDura:  "168h", // 7 day
		MaxSession: 3,
	}
	if err := conf.GetSvcJson(conf.GetSystemID(), "", &cfg); err != nil {
		flog.Error(err, "conf.GetSvcJson()")
		return cfg, err
	}

	evconf := emailsender.Config{}
	if err := conf.GetSvcJson(conf.GetSystemID(), "/emailverify", &evconf); err != nil {
		flog.Error(err, "conf.GetSvcJson", "path", "/emailverify")
		return cfg, err
	}
	cfg.ESconf = evconf

	return cfg, nil
}
