package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

var slog = log.New()

type setupCallback func(e *echo.Echo)

var SwagHandler echo.HandlerFunc

type config struct {
	Addr        string `json:"addr"`
	EnableGzip  bool   `json:"enableGzip"`
	MetricsAddr string `json:"metricsAddr"`
	SvcName     string `json:"svcName"`
}

func App(conf configuration.Configuration, systemID int, cb setupCallback) error {
	if conf == nil {
		c, err := configuration.Default()
		if err != nil {
			return err
		}
		conf = c
	}

	srv := echo.New()
	srv.HideBanner = true
	srv.Validator = formValidator
	//srv.Renderer = &Template{}
	c := getWebConfig(conf, systemID)
	svcName := c.SvcName
	if svcName == "" {
		svcName = strconv.Itoa(systemID)
	}

	srv.Use(middleware.Recover(), otelecho.Middleware(svcName), middleware.Logger())

	cb(srv)

	srv.GET("/healthy", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK!")
	})

	if SwagHandler != nil {
		srv.GET("/doc/*", SwagHandler)
	}

	if c.EnableGzip {
		srv.Use(middleware.Gzip())
	}

	if c.MetricsAddr != "" {
		go startMetricSrv(c.MetricsAddr)
	}

	go func() {
		if err := srv.Start(c.Addr); err != nil {
			slog.Errorf("Echo srv err: %v", err)
		}
	}()

	// waiting Siganl to exit server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Exit Echo server")
	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Error("Exit Echo server err: ", err)
	}

	return nil
}

func getWebConfig(conf configuration.Configuration, systemID int) config {
	c := config{}
	conf.GetJson("/microservice/framwork/web/"+strconv.Itoa(systemID), &c)
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.SvcName == "" {
		c.SvcName = strconv.Itoa(systemID)
	}
	return c
}

func startMetricSrv(addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.Handler()
		h.ServeHTTP(w, r)
	})
	if addr == "" {
		addr = ":7070"
	}
	slog.Infof("starting metric srv: %s\n" + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Infof("metric srv err: %+v\n", err)
	}
}
