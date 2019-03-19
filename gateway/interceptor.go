package gateway

import (
	"fmt"
	"lincoln/gohelper"
	"lincoln/smartgateway/breaker"
	"lincoln/smartgateway/ratelimit"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata" // grpc metadata包
)

var (
	md5Key                = "378799e6bcc25ffd1f3a51b"
	defaultTriggerCircuit = triggerCircuitBreaker //触发熔断器条件
	maxRequest            = 30                    //熔断器半开时最大请求数
	beHalfOpenInterval    = time.Second * 30      //熔断器从Open 变到beHalfOpenInterval 的时间
	clearInterval         = time.Second * 60      //Close 时， 定时清理计数器
)

// Interceptor 拦截器
func Interceptor() grpc.StreamServerInterceptor {

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		//身份验证
		err := auth(ctx)
		if err != nil {
			return err
		}

		//流量验证
		err = rateVerify(info.FullMethod)
		if err != nil {
			return err
		}

		//断路器包装
		circuitBreaker := breaker.NewBreaker(info.FullMethod, defaultTriggerCircuit, maxRequest, beHalfOpenInterval, clearInterval)

		//要往下执行的
		doHandler := func(ctx context.Context) error { return handler(srv, ss) }

		//执行不下的处理
		failHandler := func(ctx context.Context, err error) error { return grpc.Errorf(codes.FailedPrecondition, err.Error()) }

		return circuitBreaker.Handle(ctx, doHandler, failHandler)
	}
}

//token验证
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	appid, ok := md["appid"]

	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	appkey, ok := md["appkey"]

	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	crdate, ok := md["crdate"]

	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	gettime, err := gohelper.GetTime(crdate[0])
	fmt.Printf("time:%v \r\n", gettime)
	if err != nil {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	//过期
	if gettime.Before(time.Now().Add(time.Minute * -30)) {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	token, ok := md["token"]

	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

	//开始验证
	data := appid[0] + "&" + appkey[0] + "&" + crdate[0] + "&" + md5Key

	ok = gohelper.MD5Check(data, token[0])

	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "验证失败")
	}

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

//触发熔断器
func triggerCircuitBreaker(c breaker.Counts) bool {
	if c.ContinuesFail > 5 {
		return true
	}

	return false
}
