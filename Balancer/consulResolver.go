package balancer

type consulResolver struct {
	notify chan serviceNotify
}

func (cr *consulResolver) Init() {
	//初始化consul

	//初始化 notify
}

func (cr *consulResolver) GetEndPoint(service string) (addrs []string) {
	//根据服务名称获取到服务地址
	addrs = []string{"127.0.0.0:8100"}

	//开始监控
	go func() {
		for {
			newAddrs := cr.Watches(service)

			//有更改， 通知
			if newAddrs != nil {
				cr.notify <- serviceNotify{
					serviceName: service,
					addresses:   newAddrs,
				}
			}
		}
	}()

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