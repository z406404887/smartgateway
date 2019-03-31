package main

import (
	"lincoln/smartgateway/proto/test"
	"lincoln/smartgateway/server/grpc"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	consulIPPort = "192.168.1.105:8500"
	grpcIP       = "192.168.1.102"
	grpcPort     = 8001
)

// grpc + consul

func main() {

	//初始化grpc
	s := grpc.NewServer()

	//health 对外健康检查服务
	hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hsrv)

	//要提供给客户端的test 服务
	testserver := grpcserver.NewTestServer(grpcIP, grpcPort)
	test.RegisterBasicServiceServer(s, testserver)

	//注册grpc服务
	reflection.Register(s)

	//开grpc端口监听
	listen, err := net.Listen("tcp", grpcIP+":"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Fatalf("listen localhost:8001 fail, err :%v\n", err)
		return
	}

	if err := s.Serve(listen); err != nil {
		log.Fatalf("serve fail, err :%v\n", err)
	}
}
