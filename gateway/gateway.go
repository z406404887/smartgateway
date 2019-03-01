package gateway

import (
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcIP   = "127.0.0.1"
	grpcPort = 9001
)

//Run 网关入口
func Run() {
	//编码
	rcServerOption := grpc.CustomCodec(NewRawCodec())

	//Handle
	handle := &UnknowServerHandler{}

	dcServerOption := grpc.UnknownServiceHandler(handle.Handler)

	//初始化grpc
	s := grpc.NewServer(rcServerOption, dcServerOption)

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
