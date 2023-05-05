package configuration

import (
	"errors"
	"os"

	"encoding/base64"

	"github.com/jinvei/microservice/base/framework/configuration/store"
	"github.com/jinvei/microservice/base/framework/log"
)

var flog = log.Default

type Configuration interface {
	Get(path string) (string, error)
	GetJson(path string, obj interface{}) error
	GetSvcJson(systemID, subpath string, obj interface{}) error
	SetSystemID(id string)
	GetSystemID() string
}

func Default() (Configuration, error) {
	token := os.Getenv("MICROSERVICE_CONFIGURATION_TOKEN")
	if token == "" {
		return nil, errors.New("env `MICROSERVICE_CONFIGURATION_TOKEN` is empty")
	}
	return store.NewEtcdStore(token)
}

func DefaultOrDie() Configuration {
	token := os.Getenv("MICROSERVICE_CONFIGURATION_TOKEN")
	if token == "" {
		token = base64.StdEncoding.EncodeToString([]byte("{}"))
		flog.Warn("env `MICROSERVICE_CONFIGURATION_TOKEN` is empty. use default value", "MICROSERVICE_CONFIGURATION_TOKEN", token)
	}
	conf, err := store.NewEtcdStore(token)
	if err != nil {
		panic(err)
	}
	return conf
}
