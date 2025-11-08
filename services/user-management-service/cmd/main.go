package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	firebase "github.com/wnmay/horo/services/user-management-service/internal/adapters/auth"
	"github.com/wnmay/horo/services/user-management-service/internal/adapters/db"
	grpcadapter "github.com/wnmay/horo/services/user-management-service/internal/adapters/grpc"
	httpadapter "github.com/wnmay/horo/services/user-management-service/internal/adapters/http"
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

	// Init db adapter
	userRepo, err := db.NewMongoUserRepository(cfg.MongoURI, cfg.MongoDBName, cfg.UserCollectionName)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	// Init firebase and adapter
	ctx := context.Background()
	firebaseClient := firebase.InitFirebase(ctx, cfg.FirebaseAccountKeyFile)
	firebaseAdapter := firebase.NewFirebaseAuthAdapter(firebaseClient)

	// Init core service
	userApp := app.NewUserManagementService(firebaseAdapter, userRepo)
	authApp := app.NewAuthService(firebaseAdapter)

	// Create grpc server
	authServer := grpcadapter.NewAuthServer(authApp)

	grpcServer := grpc.NewServer()

	// Register auth service on gRPC (user registration is now HTTP)
	proto.RegisterAuthServiceServer(grpcServer, authServer)

	// Create HTTP handler
	httpHandler := httpadapter.NewHTTPHandler(userApp, authApp)

	// Use WaitGroup to run both servers
	var wg sync.WaitGroup
	wg.Add(2)

	// Start gRPC server
	go func() {
		defer wg.Done()
		address := fmt.Sprintf(":%s", cfg.GRPCPort)
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port %s: %v", cfg.GRPCPort, err)
		}

		log.Printf("gRPC server started on port %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		defer wg.Done()
		if err := httpadapter.StartHTTPServer(httpHandler, cfg.HTTPPort); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	wg.Wait()
}
