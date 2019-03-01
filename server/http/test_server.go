package main

import (
	"flag"
	gw "lincoln/smartgateway/proto/test"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	echoEndPoint = flag.String("echo_endpoint", "localhost:8001", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterBasicServiceHandlerFromEndpoint(ctx, mux, *echoEndPoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":50052", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
