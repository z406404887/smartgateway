package gateway

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

type UnknowServerHandler struct {
}

//Handler 该handler以gRPC server的模式来接受数据流，并将受到的数据转发到指定的connection中
func (s *UnknowServerHandler) Handler(srv interface{}, serverStream grpc.ServerStream) error {
	// 获取请求流的目的接口名称
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return grpc.Errorf(codes.Internal, "failed to get method from server stream")
	}

	outgoingCtx := serverStream.Context()
	endpoint := "127.0.0.1:8001" //根据fullMethodName 从consul获取

	// 中转 目的服务方
	backendConn, err := grpc.DialContext(outgoingCtx, endpoint, grpc.WithCodec(NewRawCodec()), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer backendConn.Close()

	// 封装为clientStream
	clientCtx, clientCancel := context.WithCancel(outgoingCtx)

	clientStream, err := grpc.NewClientStream(clientCtx, clientStreamDescForProxying, backendConn, fullMethodName)
	if err != nil {
		return err
	}

	// 启动流控，目的方->请求方
	s2cErrChan := s.forwardServerToClient(serverStream, clientStream)
	// 启动流控，请求方->目的方
	c2sErrChan := s.forwardClientToServer(clientStream, serverStream)

	// 数据流结束处理 & 错误处理
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan:
			if s2cErr == io.EOF {
				// 正常结束
				clientStream.CloseSend()
				break
			} else {
				// 错误处理 (如链接断开、读错误等)
				clientCancel()
				return grpc.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// 设置Trailer
			serverStream.SetTrailer(clientStream.Trailer())
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}
	}
	return grpc.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
}

func (s *UnknowServerHandler) forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		// 设置*bridge结构作为RecvMsg的参数，
		// *bridge即为我们自定义codec中使用到的数据结构
		f := &bridge{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err
				break
			}
			if i == 0 {
				// grpc中客户端到服务器的header只能在第一个客户端消息后才可以读取到，
				// 同时又必须在flush第一个msg之前写入到流中。
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (s *UnknowServerHandler) forwardServerToClient(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &bridge{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}