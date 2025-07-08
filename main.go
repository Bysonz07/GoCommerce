package main

import (
	"log"
	"net"

	"ecommerce/api"
	"ecommerce/db/sqlc"
	"ecommerce/internal/service"
	"ecommerce/internal/worker"

	"github.com/redis/go-redis/v8"
	"google.golang.org/grpc"
)

func main() {
	// Initialize database and Redis
	db := sqlc.NewDB()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Start Redis worker
	go worker.StartRedisWorker(redisClient)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterECommerceServiceServer(grpcServer, service.NewECommerceService(db, redisClient))

	log.Println("Starting gRPC server on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
