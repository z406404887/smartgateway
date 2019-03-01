package grpcserver

//ServerBase 基础服务
type ServerBase struct {
	IPAddress string //服务端地址
	Port      int    //服务端口
	ID        string //服务id
	Name      string //服务名称
}

//newServerBase 初始化 基础服务
func newServerBase(ipAddress string, port int, name, id string) (*ServerBase, error) {

	sb := &ServerBase{
		IPAddress: ipAddress,
		Port:      port,
		ID:        id,
		Name:      name,
	}

	return sb, nil
}
