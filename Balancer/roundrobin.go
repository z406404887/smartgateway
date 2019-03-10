package balancer

import "sync"

//RoundRobin 负载均衡： 轮询
type RoundRobin struct {
	mu        sync.Mutex
	Resolver  resolver
	EndPoints map[string]endPoints
}

//NewResolver 初始化 resolver
func (r *RoundRobin) NewResolver(rs resolver) {
	r.Resolver = rs

	//监控通知
	go r.Up()
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
		ep = endPoints{
			addresses: arrayAddress,
			next:      0,
		}
		r.EndPoints[serviceName] = ep
	}

	//开始对该组地址使用轮询算法
	r.mu.Lock()
	ep.next = (ep.next + 1) % len(ep.addresses)
	addr = ep.addresses[ep.next]
	r.mu.Unlock()

	return addr
}

//Up 更新服务地址
func (r *RoundRobin) Up() {
	go func() {
		select {
		case newService := <-r.Resolver.Notify():
			if _, ok := r.EndPoints[newService.serviceName]; ok {
				r.mu.Lock()
				r.EndPoints[newService.serviceName] = endPoints{
					addresses: newService.addresses,
					next:      0,
				}
				r.mu.Unlock()
			}
		}
	}()
}
