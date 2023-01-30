package datasource

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jinvei/microservice/base/framework/configuration"
	confkey "github.com/jinvei/microservice/base/framework/configuration/keys"
	"github.com/jinvei/microservice/base/framework/log"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
	xname "xorm.io/xorm/names"
)

var flog = log.New()

type Config struct {
	Dialect      string        `json:"dialect"` // 数据库类型 Mysql/SqLite/PostgreSQL
	Dsn          string        `json:"dsn"`     // 数据库链接
	Debug        bool          `json:"debug"`
	EnableLog    bool          `json:"enableLog"`
	Prefix       string        `json:"prefix"`       // 表名前缀
	MinPoolSize  int           `json:"minPoolSize"`  // pool最大空闲数
	MaxPoolSize  int           `json:"maxPoolSize"`  // pool最大连接数
	IdleTimeout  time.Duration `json:"idleTimeout"`  // 连接最长存活时间
	QueryTimeout time.Duration `json:"queryTimeout"` // 查询超时时间
	ExecTimeout  time.Duration `json:"execTimeout"`  // 执行超时时间
	TranTimeout  time.Duration `json:"tranTimeout"`  // 事务超时时间
}

type DataSource struct {
	conf     configuration.Configuration
	systemID int
}

func (s *DataSource) Orm() *xorm.Engine {
	flog.Debugf("Init xorm. SystemID='%d'", s.systemID)

	c := s.getConfig()

	flog.Debugf("Datasource config: %v", c)

	xe, err := xorm.NewEngine(c.Dialect, c.Dsn)
	if err != nil {
		flog.Error(err)
		panic(fmt.Sprintf("Failed to init xorm: %+v", err))
	}

	xe.ShowSQL(c.Debug)
	if c.EnableLog {
		xe.Logger().SetLevel(xlog.LOG_DEBUG)
	}

	xe.SetTableMapper(xname.NewPrefixMapper(xname.SnakeMapper{}, c.Prefix))
	xe.SetMaxIdleConns(c.MinPoolSize)
	xe.SetMaxOpenConns(c.MaxPoolSize)
	xe.SetConnMaxLifetime(time.Duration(c.IdleTimeout))

	return nil
}

func New(conf configuration.Configuration, systemID int) *DataSource {
	if conf == nil {
		conf = configuration.DefaultOrDie()
	}

	return &DataSource{
		conf:     conf,
		systemID: systemID,
	}
}

func (s *DataSource) getConfig() *Config {
	c := Config{}
	path := filepath.Join(confkey.FwDatasource, strconv.Itoa(s.systemID))
	if err := s.conf.GetJson(path, &c); err != nil {
		flog.Error(err)
		panic(err)
	}
	return &c
}
