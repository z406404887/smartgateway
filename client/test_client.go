package main

import (
	"context"
	"fmt"
	"lincoln/gohelper"
	"lincoln/smartgateway/proto/test"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	consulIPPort = "192.168.1.105:8500"
	md5Key       = "378799e6bcc25ffd1f3a51b"
)

// func main() {
// 	//获取consul对象(当前的consul地址是测试的， 正常的开发应该走 dns或者 以 虚拟Ip 来作为consul地址)
// 	register, err := consul.NewConsulRegister(consulIPPort)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	//到consul 获取服务
// 	service, err := register.GetService("38600a02-b49d-476a-a374-5e466b68bf52")
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	//初始化grpc
// 	conn, err := grpc.Dial(service.Address+fmt.Sprintf(":%d", service.Port), grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("fail to connect service, err :%v\n", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// 初始化服务客户端
// 	client := test.NewBasicServiceClient(conn)

// 	// 调用Say服务
// 	for i := 0; i < 500; i++ {
// 		fmt.Println("开始调用")
// 		resp, err := client.Say(context.Background(), &test.Request{Username: "123"})
// 		if err != nil {
// 			log.Fatalf("Say err:%v\n", err)
// 		}

// 		fmt.Printf("Say response:%v\n", resp)
// 	}
// }
func main() {

	//初始化grpc
	conn, err := grpc.Dial("127.0.0.1:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect service, err :%v\n", err)
		return
	}
	defer conn.Close()

	// 初始化服务客户端
	client := test.NewBasicServiceClient(conn)

	//新建一个有 metadata 的 context
	appid := "456"
	appkey := "456"
	crdate := gohelper.GetTimeStr(time.Now().Add(time.Minute * -28))
	fmt.Println("时间：" + crdate)
	//生成token
	content := appid + "&" + appkey + "&" + crdate + "&" + md5Key
	token := gohelper.MD5Encode(content)

	md := metadata.Pairs("appid", appid, "appkey", appkey, "crdate", crdate, "token", token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	for i := 0; i < 200; i++ {
		go func() {
			fmt.Println("开始调用")
			resp, err := client.Say(ctx, &test.Request{Username: "123"})
			if err != nil {
				log.Fatalf("Say err:%v\n", err)
			}

			fmt.Printf("Say response:%v\n", resp)
		}()
	}

	time.Sleep(20 * time.Second)
}
