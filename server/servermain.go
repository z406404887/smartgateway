package main

import (
	"lincoln/smartgateway/proto/test"
	"lincoln/smartgateway/registry/consul"
	"lincoln/smartgateway/server/grpc"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	consulIPPort = "192.168.1.105:8500"
	grpcIP       = "127.0.0.1"
	grpcPort     = 8001
	opts         []grpc.ServerOption //拦截器
)

func main() {

	//初始化grpc
	opts = append(opts, grpc.UnaryInterceptor(grpcserver.Interceptor))
	s := grpc.NewServer(opts...)

	//要提供给客户端的test 服务
	testserver := grpcserver.NewTestServer(grpcIP, grpcPort)
	test.RegisterBasicServiceServer(s, testserver)

	// 在consul注册服务(若consul集群了， 也可类似的注册其他机子)
	consul.RegisterService(consulIPPort, testserver.ID, testserver.Name, testserver.IPAddress, testserver.Port)

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
