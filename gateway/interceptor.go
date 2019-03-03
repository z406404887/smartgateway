package gateway

import (
	"fmt"
	"lincoln/smartgateway/ratelimit"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata" // grpc metadata包
)

// Interceptor 拦截器
func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//身份验证
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	//流量验证
	err = rateVerify(info.FullMethod)
	if err != nil {
		return nil, err
	}

	// 继续处理请求
	return handler(ctx, req)
}

//token验证
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	authInfo, ok := md["auth"]

	if !ok || len(authInfo) < 2 {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	name := authInfo[0]
	token := authInfo[1]

	//开始验证

	return nil
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
