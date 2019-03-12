package balancer

import "sync"

//RoundRobin 负载均衡： 轮询
type RoundRobin struct {
	mu        sync.Mutex
	Resolver  resolver
	EndPoints map[string]*endPoints
}

//NewRoundRobin 返回 接口对象
func NewRoundRobin(consulIP string) Base {

	//初始化consul 对象管理
	consulResolver, err := NewConsulResolver(consulIP)
	if err != nil {
		return &RoundRobin{}
	}

	//初始化轮询对象
	round := &RoundRobin{
		EndPoints: make(map[string]*endPoints),
	}
	round.NewResolver(consulResolver)

	return round
}

//NewResolver 初始化 resolver
func (r *RoundRobin) NewResolver(rs resolver) {
	r.Resolver = rs
}

//GetAddrByAlgorithms 获取 负载算法后的地址
func (r *RoundRobin) GetAddrByAlgorithms(serviceName string) (addr string) {
	//获取服务的一组地址
	ep, ok := r.EndPoints[serviceName]

	//还没有该服务地址
	if !ok {
		//到对应的 resolver 里 获取服务地址
		arrayAddress := r.Resolver.GetEndPoint(serviceName)

		//添加到RoundRobin中
		ep = &endPoints{
			addresses: arrayAddress,
			next:      0,
		}
		r.EndPoints[serviceName] = ep

		//开始监控该服务的变化
		r.WatchAndUp(serviceName)
	}

	//开始对该组地址使用轮询算法
	r.mu.Lock()
	ep.next = (ep.next + 1) % len(ep.addresses)
	addr = ep.addresses[ep.next]
	r.mu.Unlock()

	return addr
}

//WatchAndUp RoundRobin让Resolver 监控服务
func (r *RoundRobin) WatchAndUp(service string) {

	//交给Resolver,  开始监控该服务
	go r.Resolver.Watches(service)

	//Resolver 服务通知处理
	go func() {
		for {
			select {
			case changeNotify := <-r.Resolver.Notify(): //Resolver有服务更改通知过来
				if _, ok := r.EndPoints[changeNotify.serviceName]; ok {
					r.EndPoints[changeNotify.serviceName] = &endPoints{
						addresses: changeNotify.addresses,
						next:      0,
					}
				}
			}

		}
	}()
}
