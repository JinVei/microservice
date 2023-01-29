package configuration

import (
	"errors"
	"os"

	"github.com/jinvei/microservice/base/framework/configuration/store"
)

type Configuration interface {
	Get(path string) (string, error)
	GetJson(path string, obj interface{}) error
}

func Default() (Configuration, error) {
	token := os.Getenv("MICROSERVICE_CONFIGURATION_TOKEN")
	if token == "" {
		return nil, errors.New("env `MICROSERVICE_CONFIGURATION_TOKEN` is empty")
	}
	return store.NewEtcdStore(token)
}
