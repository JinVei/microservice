package rpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var slog = log.New()

type setupCallback func(e *grpc.Server)

type config struct {
	Addr    string `json:"addr"`
	SvcName string `json:"svcName"`
}

func Serve(conf configuration.Configuration, systemID int, cb setupCallback) error {
	if conf == nil {
		c, err := configuration.Default()
		if err != nil {
			return err
		}
		conf = c
	}

	sconfig := config{}

	if err := conf.GetJson(filepath.Join("/microservice/framework/rpc/", strconv.Itoa(systemID)), &sconfig); err != nil {
		return err
	}
	if sconfig.Addr == "" {
		sconfig.Addr = ":9090"
	}

	srv := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			//tags.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
			//logging.StreamServerInterceptor(newLog()),
			//auth.StreamServerInterceptor(myAuthFunction),
			recovery.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			//tags.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
			//logging.UnaryServerInterceptor(newLog()),
			//auth.UnaryServerInterceptor(myAuthFunction),
			recovery.UnaryServerInterceptor(),
		),
	)

	cb(srv)

	lis, err := net.Listen("tcp", sconfig.Addr)
	if err != nil {
		slog.Error("failed to listen: ", err)
		return err
	}

	slog.Infof("Run rpc server. SystemID=%d, ListenAt='%s'", systemID, sconfig.Addr)

	go func() {
		if err := srv.Serve(lis); err != nil {
			slog.Error("failed to serve: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Exit rpc srv")
	srv.GracefulStop()

	return nil
}

func NewClientConn(conf configuration.Configuration, systemID int) (*grpc.ClientConn, error) {
	if conf == nil {
		c, err := configuration.Default()
		if err != nil {
			return nil, err
		}
		conf = c
	}

	sconfig := config{}
	confpath := filepath.Join("/microservice/framework/rpc/", strconv.Itoa(systemID))
	if err := conf.GetJson(confpath, &sconfig); err != nil {
		return nil, err
	}

	if sconfig.Addr == "" {
		return nil, fmt.Errorf("addr field is empty. path='%s'", sconfig.Addr)
	}

	conn, err := grpc.Dial(sconfig.Addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		slog.Errorf("Failed to connect: %v\n", err)
		return nil, err
	}
	return conn, err
}
