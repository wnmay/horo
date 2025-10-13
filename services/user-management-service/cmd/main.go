package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/wnmay/horo/services/user-management-service/internal/adapters/db"
	"github.com/wnmay/horo/services/user-management-service/internal/adapters/firebase"
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

	userRepo, err := db.NewMongoUserRepository(cfg.MongoURI, cfg.MongoDBName, cfg.UserCollectionName)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	ctx := context.Background()
	firebaseClient := firebase.InitFirebase(ctx, cfg.FirebaseAccountKeyFile)
	firebaseAdapter := firebase.NewAuthAdapter(firebaseClient)

	userApp := app.NewUserManagementService(firebaseAdapter, userRepo)
	authApp := app.NewAuthService(firebaseAdapter)

	userServer := grpcadapter.NewUserManagementServer(userApp)
	authServer := grpcadapter.NewAuthServer(authApp)


	grpcServer := grpc.NewServer()

	// Register both services on the same server
	proto.RegisterUserManagementServiceServer(grpcServer, userServer)
	proto.RegisterAuthServiceServer(grpcServer, authServer)

	// Start listening
	address := fmt.Sprintf(":%s", cfg.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("gRPC server started on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
