package grpcserver

import (
	"fmt"
	"lincoln/smartgateway/ratelimit"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	// grpc 响应状态码
	// grpc metadata包
)

// Interceptor 拦截器
func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	err := rateVerify(info.FullMethod)
	if err != nil {
		return nil, err
	}

	// 继续处理请求
	return handler(ctx, req)
}

//请求的接口流量验证
func rateVerify(FullMethod string) error {
	bucket := ratelimit.NewTokenBucket(FullMethod, nil)
	ok := bucket.Take(100000)

	if !ok {
		return grpc.Errorf(codes.ResourceExhausted, "已达到最大请求数")
	}
	fmt.Println("请求通过")
	return nil
}
