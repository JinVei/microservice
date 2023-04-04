package cache

import (
	"path/filepath"

	"github.com/jinvei/microservice/base/framework/configuration"
	confkeys "github.com/jinvei/microservice/base/framework/configuration/keys"
	"github.com/jinvei/microservice/base/framework/log"
	"github.com/redis/go-redis/v9"
)

var flog = log.New()

type redisConfig struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func RedisClient(conf configuration.Configuration) *redis.Client {
	if conf == nil {
		conf = configuration.DefaultOrDie()
	}

	c := getRedsiConfig(conf)

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
	})

	return rdb

}

func getRedsiConfig(conf configuration.Configuration) redisConfig {
	systemID := configuration.GetSystemID()
	if systemID == "" {
		panic("systemID is empty. should set SystemID by 'configuration.SetSystemID(id)'")
	}
	path := filepath.Join(confkeys.FwCacheRedis, systemID)
	c := redisConfig{}
	err := conf.GetJson(path, &c)
	if err != nil {
		flog.Error(err, "conf.GetJson", "path", path, "redisconfig", c)
		panic(err)
	}
	return c
}
