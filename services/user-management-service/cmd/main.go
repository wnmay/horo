package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/wnmay/horo/services/user-management-service/internal/adapters/db"
	grpcadapter "github.com/wnmay/horo/services/user-management-service/internal/adapters/grpc"
	"github.com/wnmay/horo/services/user-management-service/internal/app"
	"github.com/wnmay/horo/services/user-management-service/internal/config"
	"github.com/wnmay/horo/shared/env"
	proto "github.com/wnmay/horo/shared/proto/user-management"
	"google.golang.org/grpc"
)

const (
	service_name = "user-management-service"
)

func main() {
	_ = env.LoadEnv(service_name)
	cfg := config.LoadConfig()

	repo, err := db.NewMongoUserRepository(cfg.MongoURI, cfg.MongoDBName, cfg.UserCollectionName)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	ctx := context.Background()
	userManagementService := app.NewUserManagementService(ctx, repo, cfg)
	grpcServer := grpcadapter.NewGRPCServer(userManagementService)
	server := grpc.NewServer()
	proto.RegisterUserManagementServiceServer(server, grpcServer)

	// Start listening
	address := fmt.Sprintf(":%s", cfg.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("gRPC server started on port %s", cfg.GRPCPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
