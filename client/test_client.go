package main

import (
	"context"
	"fmt"
	"lincoln/smartgateway/proto/test"
	"lincoln/smartgateway/registry/consul"
	"log"

	"google.golang.org/grpc"
)

var (
	consulIPPort = "192.168.1.105:8500"
)

func main() {
	//获取consul对象(当前的consul地址是测试的， 正常的开发应该走 dns或者 以 虚拟Ip 来作为consul地址)
	register, err := consul.NewConsulRegister(consulIPPort)
	if err != nil {
		log.Fatal(err)
		return
	}

	//到consul 获取服务
	service, err := register.GetService("38600a02-b49d-476a-a374-5e466b68bf52")
	if err != nil {
		log.Fatal(err)
		return
	}

	//初始化grpc
	conn, err := grpc.Dial(service.Address+fmt.Sprintf(":%d", service.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect service, err :%v\n", err)
		return
	}
	defer conn.Close()

	// 初始化服务客户端
	client := test.NewBasicServiceClient(conn)

	// 调用Say服务
	for i := 0; i < 500; i++ {
		fmt.Println("开始调用")
		resp, err := client.Say(context.Background(), &test.Request{Username: "123"})
		if err != nil {
			log.Fatalf("Say err:%v\n", err)
		}

		fmt.Printf("Say response:%v\n", resp)
	}
}