package main

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"

	//"github.com/gomodule/redigo/redis"
	"github.com/segmentio/feature-flag/pkg/proto"
	redis_pkg "github.com/segmentio/feature-flag/pkg/redis"
	"github.com/segmentio/feature-flag/pkg/service"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

//configs ,redis ADddr, connections,timeout
func main() {
	const maxConnections = 10
	timeout := 1 * time.Minute
	redisAddr := fmt.Sprintf("%s:%s", "34.125.145.148", "6379")

	redisClient, _ := redis_pkg.GetRedisClient(redisAddr, maxConnections, timeout)
	var wg sync.WaitGroup
	wg.Add(1)

	redisInterface := redis_pkg.NewPool(redisClient)
	go startGrpc(redisInterface)
	go startHTTP()
	defer redisClient.Close()
	wg.Wait()

}
func startGrpc(redisInterface redis_pkg.PoolInterface) {
	lis, err := net.Listen("tcp", "localhost:5566")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	//can pass this as flag
	mgr := service.NewConcurrentManager(3, 5, 2, 30*time.Second)

	srv := service.NewFlagoServer(redisInterface, mgr)

	proto.RegisterFlagoServiceServer(grpcServer, srv)
	log.Info("gRPC server ready...")
	grpcServer.Serve(lis)
}
func startHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	conn, err := grpc.Dial("localhost:5566", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	rmux := runtime.NewServeMux()
	client := proto.NewFlagoServiceClient(conn)
	err = proto.RegisterFlagoServiceHandlerClient(ctx, rmux, client)
	if err != nil {
		log.Fatalf("failed to register http client to grpc server", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	err = http.ListenAndServe("localhost:9080", logRequest(mux))
	if err != nil {
		log.Fatalf("failed to listen and serve %v", err)
	}
	log.Info("http server ready\n")

}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// All flags
// worker count
//redis max connections
//redis timeout
//retry factor
//jitter
//retry count
