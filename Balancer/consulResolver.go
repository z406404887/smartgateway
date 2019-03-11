package balancer

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
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
	// services, metainfo, err := w.client.Health().Service(w.service, "", true, &api.QueryOptions{
	// 	WaitIndex: w.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	// })
	// if err != nil {
	// 	log.Fatal("error retrieving instances from Consul: %v", err)
	// }
	// w.lastIndex = metainfo.LastIndex
	return nil
}

func (cr *consulResolver) Notify() (notify <-chan serviceNotify) {
	return cr.notify
}
