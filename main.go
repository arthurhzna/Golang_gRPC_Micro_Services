package main

import (
	"log"
	"net"
	"os"

	"github.com/arthurhzna/Golang_gRPC/internal/handler"
	"github.com/arthurhzna/Golang_gRPC/pb/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	if os.Getenv("ENVIRONMENT") == "DEV" {
		reflection.Register(grpcServer)
	}

	serviceHandler := handler.NewServiceHandler()

	grpcServer.Serve(lis)

	service.RegisterHelloWorldServiceServer(grpcServer, serviceHandler)

}
