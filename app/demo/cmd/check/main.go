package main

import (
	"context"
	"fmt"

	"github.com/jinvei/microservice/base/framework/rpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const DemoSystemID = 11001

func main() {

	conn, err := rpc.NewClientConn(nil, DemoSystemID)

	client := healthpb.NewHealthClient(conn)
	response, err := client.Check(context.Background(), &healthpb.HealthCheckRequest{Service: "service_name"})
	if err != nil {
		fmt.Printf("Failed to perform health check: %v\n", err)
		return
	}

	if response.GetStatus() == healthpb.HealthCheckResponse_SERVING {
		fmt.Println("Service is serving.")
	} else {
		fmt.Println("Service is NOT serving.")
	}
}
