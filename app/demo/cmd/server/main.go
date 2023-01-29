package main

import (
	"context"
	"log"

	"github.com/jinvei/microservice/base/framework/rpc"
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

	err := rpc.Serve(nil, systemID, func(srv *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(srv, &server{})
	})
	if err != nil {
		log.Println(err)
	}
}
