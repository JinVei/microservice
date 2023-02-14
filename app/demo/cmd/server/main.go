package main

import (
	"context"
	"log"
	"sync"

	"github.com/jinvei/microservice/base/framework/rpc"
	"github.com/jinvei/microservice/base/framework/web"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const systemID = 11001

type server struct{}

func (s *server) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *server) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return nil
}

func main() {
	var waitSrv sync.WaitGroup
	waitSrv.Add(2)
	go func() {
		defer waitSrv.Done()
		err := rpc.Serve(nil, systemID, func(srv *grpc.Server) {
			grpc_health_v1.RegisterHealthServer(srv, &server{})
		})
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer waitSrv.Done()
		web.App(nil, systemID, func(e *echo.Echo) {
			e.GET("/demo/test", func(c echo.Context) error {
				c.JSON(200, "Well Done!")
				return nil
			})
		})
	}()

	waitSrv.Wait()
}
