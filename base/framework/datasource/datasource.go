package datasource

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinvei/microservice/base/framework/configuration"
	confkey "github.com/jinvei/microservice/base/framework/configuration/keys"
	"github.com/jinvei/microservice/base/framework/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
	xname "xorm.io/xorm/names"
)

var flog = log.New()

type Config struct {
	Dialect       string `json:"dialect"` // 数据库类型 Mysql/SqLite/PostgreSQL
	Dsn           string `json:"dsn"`     // 数据库链接
	Debug         bool   `json:"debug"`
	EnableLog     bool   `json:"enableLog"`
	Prefix        string `json:"prefix"`       // 表名前缀
	MinPoolSize   int    `json:"minPoolSize"`  // pool最大空闲数
	MaxPoolSize   int    `json:"maxPoolSize"`  // pool最大连接数
	IdleTimeout   string `json:"idleTimeout"`  // 连接最长存活时间
	QueryTimeout  string `json:"queryTimeout"` // 查询超时时间
	ExecTimeout   string `json:"execTimeout"`  // 执行超时时间
	TranTimeout   string `json:"tranTimeout"`  // 事务超时时间
	SingularTable bool   `json:"singularTable"`
}

type DataSource struct {
	conf     configuration.Configuration
	systemID int
}

func (s *DataSource) Orm() *xorm.Engine {
	flog.Debugf("Init xorm. SystemID='%d'", s.systemID)

	c := s.getConfig()

	flog.Info("Init Datasource config", "config", c)

	xe, err := xorm.NewEngine(c.Dialect, c.Dsn)
	if err != nil {
		flog.Error(err, "xorm.NewEngine()")
		panic(fmt.Sprintf("Failed to init xorm: %+v", err))
	}

	xe.ShowSQL(c.Debug)
	if c.EnableLog {
		xe.Logger().SetLevel(xlog.LOG_DEBUG)
	}

	xe.SetTableMapper(xname.NewPrefixMapper(xname.SnakeMapper{}, c.Prefix))
	xe.SetMaxIdleConns(c.MinPoolSize)
	xe.SetMaxOpenConns(c.MaxPoolSize)

	d, err := time.ParseDuration(c.IdleTimeout)
	if err != nil {
		flog.Warn("parse c.IdleTimeout error. set to default '10s'")
		d = 10 * time.Second
	}
	xe.SetConnMaxLifetime(d)

	return xe
}

func (s *DataSource) Gorm() *gorm.DB {
	flog.Debugf("Init gorm. SystemID='%d'", s.systemID)
	c := s.getConfig()

	flog.Info("Init Datasource config", "config", c)

	var dialector gorm.Dialector
	switch c.Dialect {
	case "postgres":
		dialector = postgres.Open(c.Dsn)
	case "mysql":
		dialector = mysql.Open(c.Dsn)
	case "sqlite":
		dialector = sqlite.Open(c.Dsn)
	default:
		panic("Unknow Dialect: " + c.Dialect)
	}

	logLevel := logger.Info
	if c.EnableLog {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   c.Prefix,
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		flog.Error(err, " gorm.Open(dialector)", "cfg", c)
		panic(fmt.Sprintf("Failed to init gorm: %+v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	d, err := time.ParseDuration(c.IdleTimeout)
	if err != nil {
		flog.Warn("parse c.IdleTimeout error. set to default '10s'")
		d = 10 * time.Second
	}

	sqlDB.SetMaxIdleConns(c.MinPoolSize)
	sqlDB.SetMaxOpenConns(c.MaxPoolSize)
	sqlDB.SetConnMaxLifetime(d)

	return db
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

func defaultConfig() Config {
	return Config{
		Dialect:       "mysql",
		Debug:         true,
		EnableLog:     true,
		Prefix:        "",
		MinPoolSize:   2,
		MaxPoolSize:   10,
		IdleTimeout:   "10s",
		QueryTimeout:  "2s",
		ExecTimeout:   "2s",
		TranTimeout:   "2s",
		SingularTable: true,
	}
}

func (s *DataSource) getConfig() *Config {
	c := defaultConfig()
	path := filepath.Join(confkey.FwDatasource, strconv.Itoa(s.systemID))
	if err := s.conf.GetJson(path, &c); err != nil {
		flog.Error(err, "GetJson() err:")
		panic(err)
	}
	return &c
}
