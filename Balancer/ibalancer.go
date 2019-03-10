package balancer

//Base  负载均衡基础接口
type Base interface {
	NewResolver(resolver)                          //初始化 resolver , 从resolver 获取地址
	GetAddrByAlgorithms(name string) (addr string) //获取负载算法获取 服务最终算法
	Up(resolver)                                   //更新地址
}

//resolver 服务地址管理对象
type resolver interface {
	Build()                                      //初始化
	GetEndPoint(service string) (addrs []string) //返回 所有可选的地址
	Watches(serviceName string) (addrs []string) //监控服务地址变化
	Notify() <-chan serviceNotify                //通知服务地址变化
}

type endPoints struct {
	addresses []string
	next      int
}

type serviceNotify struct {
	serviceName string
	addresses   []string
}
