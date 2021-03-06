package consul

import (
	"errors"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
)

//Register 注册中心
type Register struct {
	client *consulapi.Client
}

//registerManage 管理注册中心
type registerManage struct {
	mu        sync.RWMutex
	Registers map[string]*Register
}

type ServerBase struct {
	IPAddress string //服务端地址
	Port      int    //服务端口
	ID        string //服务id
	Name      string //服务名称
}

var rManage registerManage

func init() {
	rManage = registerManage{Registers: make(map[string]*Register)}
}

//NewConsulRegister 初始化 consul注册中心
func NewConsulRegister(Address string) (r *Register, e error) {
	//允许大量读
	rManage.mu.RLock()

	//返回注册对象
	register, ok := rManage.Registers[Address]
	if ok {
		rManage.mu.RUnlock()
		return register, nil
	}

	rManage.mu.RUnlock()

	rManage.mu.Lock()
	defer rManage.mu.Unlock()

	//多个读锁等待写， 只有一个会启用， 其他需再次检查
	register, ok = rManage.Registers[Address]
	if ok {
		rManage.mu.RUnlock()
		return register, nil
	}

	//创建
	config := consulapi.DefaultConfig()
	config.Address = Address
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	register = &Register{client: client}
	rManage.Registers[Address] = register

	return register, nil
}

//DoRegister 注册服务到 consul
func (r *Register) DoRegister(sr *consulapi.AgentServiceRegistration) error {
	return r.client.Agent().ServiceRegister(sr)
}

//Deregister 到 consul 删除服务
func (r *Register) Deregister(ID string) error {
	return r.client.Agent().ServiceDeregister(ID)
}

//GetServices 获取consul 某个服务
func (r *Register) GetService(ID string) (*consulapi.AgentService, error) {
	services, err := r.GetAllService()
	if err != nil {
		return nil, err
	}

	service, ok := services[ID]

	if !ok {
		return nil, errors.New("Service is not exists")
	}

	return service, nil
}

//GetAllService 获取consul全部服务
func (r *Register) GetAllService() (map[string]*consulapi.AgentService, error) {
	return r.client.Agent().Services()
}

//RegisterService 注册服务
func RegisterService(consulIp string, id string, name string, ipAddress string, port int) {

	//获取consul对象(当前的consul地址是测试的， 正常的开发应该走 dns或者 以 虚拟Ip 来作为consul地址)
	register, err := NewConsulRegister(consulIp)

	if err != nil {
		panic(err)
	}

	//创建consul服务
	registration := &consulapi.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Port:    port,
		Tags:    []string{""},
		Address: ipAddress,
	}

	//注册该服务
	err = register.DoRegister(registration)
	if err != nil {
		panic(err)
	}

}
