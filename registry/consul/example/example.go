package main

import (
	"fmt"
	"lincoln/smartgateway/registry/consul"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	//注册服务
	TestRegistryService("web")

	//获取服务
	service := TestGetService("web")
	fmt.Printf("地址是:%s:%d", service.Address, service.Port)
}

func TestRegistryService(ID string) {
	register, err := consul.NewConsulRegister("192.168.1.105:8500")

	if err != nil {

		log.Fatal(err)
	}

	//创建一个新服务。
	registration := &consulapi.AgentServiceRegistration{
		ID:      ID,
		Name:    "web",
		Port:    8001,
		Tags:    []string{""},
		Address: "127.0.0.1",
	}

	//  //增加check。
	//  check := new(consulapi.AgentServiceCheck)
	//  check.HTTP = "ip:port/check"

	//  //设置超时 5s。
	//  check.Timeout = "5s"

	//  //设置间隔 5s。
	//  check.Interval = "5s"

	//  //注册check服务。
	//  registration.Check = check
	//  log.Println("get check.HTTP:", check)
	//
	//  err = client.Agent().ServiceRegister(registration)
	//
	//  if err != nil {
	//      log.Fatal("register server error : ", err)
	//  }

	//注册该服务
	err = register.DoRegister(registration)

	if err != nil {
		log.Fatal(err)
	}
}

func TestDelService(ID string) {
	register, err := consul.NewConsulRegister("192.168.1.105:8500")

	if err != nil {
		log.Fatal(err)
	}

	err = register.Deregister(ID)

	if err != nil {
		log.Fatal(err)
	}
}

func TestGetService(ID string) *consulapi.AgentService {
	register, err := consul.NewConsulRegister("192.168.1.105:8500")

	if err != nil {
		log.Fatal(err)
	}

	service, err := register.GetService(ID)

	if err != nil {
		log.Fatal(err)
	}

	return service
}
