package balancer

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

var (
	err    error
	params map[string]interface{}
	plan   *watch.Plan
	ch     chan int
)

type consulResolver struct {
	allService map[string]*consulapi.AgentService
	client     *consulapi.Client
	notify     chan serviceNotify
}

//NewConsulResolver 返回consulResolver
func NewConsulResolver(consulAddr string) (resolver, error) {

	//consul 配置
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)

	if err != nil {
		return nil, err
	}

	cr := &consulResolver{
		client: client,
		notify: make(chan serviceNotify),
	}

	//获取所有地址
	cr.allService, err = cr.client.Agent().Services()
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (cr *consulResolver) GetEndPoint(service string) (addrs []string) {
	//遍历获取
	for _, agentService := range cr.allService {
		if agentService.Service == service {

			addr := fmt.Sprintf("%s:%d", agentService.Address, agentService.Port)
			addrs = append(addrs, addr)
		}
	}

	return addrs
}

func (cr *consulResolver) Watches(service string) (addrs []string) {

	//暂时不针对某个service 监控， 直接对 allService 定时更新
	// ch = make(chan int, 1)

	// params = make(map[string]interface{})
	// params["type"] = "services"
	// params["passingonly"] = false
	// plan, err = watch.Parse(params)
	// if err != nil {
	// 	panic(err)
	// }
	// plan.Handler = func(index uint64, result interface{}) {
	// 	if entries, ok := result.([]*consulapi.ServiceEntry); ok {
	// 		fmt.Printf("serviceEntries:%v", entries)
	// 		// your code
	// 		ch <- 1
	// 	}
	// }
	// go func() {
	// 	// your consul agent addr
	// 	if err = plan.Run("192.168.1.105:8500"); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// go http.ListenAndServe(":8080", nil)

	// for {
	// 	<-ch
	// 	fmt.Printf("get change")
	// }
	return nil
}

func (cr *consulResolver) Notify() (notify <-chan serviceNotify) {
	return cr.notify
}
