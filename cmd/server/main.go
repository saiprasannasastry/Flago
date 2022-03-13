package main

import (
	"fmt"
	"github.com/apex/log"
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

// All flags
// worker count
//redis max connections
//redis timeout
//retry factor
//jitter
//retry count
