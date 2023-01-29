package grpc

import (
	"log"
	"net"

	"github.com/jinvei/microservice/base/framework/configuration"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type setupCallback func(e *grpc.Server)

func App(conf configuration.Configuration, systemID int, cb setupCallback) {

	srv := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			//tags.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
			logging.StreamServerInterceptor(newLog()),
			//auth.StreamServerInterceptor(myAuthFunction),
			recovery.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			//tags.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(newLog()),
			//auth.UnaryServerInterceptor(myAuthFunction),
			recovery.UnaryServerInterceptor(),
		),
	)

	cb(srv)

	lis, err := net.Listen("tcp", "9090") // todo
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
