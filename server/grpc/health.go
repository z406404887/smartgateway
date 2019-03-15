package grpcserver

import (
	"context"
	"lincoln/smartgateway/proto/health"
	"lincoln/smartgateway/proto/test"
)

type HealthServer struct {
}

func (h *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	 
}
