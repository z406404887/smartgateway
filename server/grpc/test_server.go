package grpcserver

import (
	"context"
	"lincoln/smartgateway/proto/test"
)

//服务相关信息
var (
	name = "/test/test_grpc"
	id   = "C66F55B4-81CB-4496-A119-523A7C1E8E11" //此ID应由外部来生成保存
)

//TestServer 用于实现BasicServiceServer
type TestServer struct {
	*ServerBase
}

//Say 定义Say方法，用于实现BasicServiceServer里面的Login方法
func (s *TestServer) Say(ctx context.Context, req *test.Request) (*test.Response, error) {
	return &test.Response{Returnmsg: "You name is " + req.Username}, nil
}

//NewTestServer 初始化服务
func NewTestServer(ipAddress string, port int) *TestServer {
	sb, err := newServerBase(ipAddress, port, name, id)

	if err != nil {
		panic(err)
	}

	//test 服务
	testserver := &TestServer{
		ServerBase: sb,
	}

	return testserver
}
