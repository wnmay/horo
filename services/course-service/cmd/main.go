package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	grpcadapter "github.com/wnmay/horo/services/course-service/internal/adapters/inbound/grpc"
	httpadapter "github.com/wnmay/horo/services/course-service/internal/adapters/inbound/http"
	dbout "github.com/wnmay/horo/services/course-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/course-service/internal/app"
	pb "github.com/wnmay/horo/shared/course/proto"
	"github.com/wnmay/horo/shared/db"
	"github.com/wnmay/horo/shared/env"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// === 1. Load configuration ===
	_ = env.LoadEnv("course-service")
	restPort := env.GetString("REST_PORT", "3002")
	grpcPort := env.GetString("GRPC_PORT", "50052")

	dbName := env.GetString("DB_NAME", "coursedb")
	cfg := db.NewMongoDefaultConfig(dbName)
	client, err := db.NewMongoClient(context.Background(), cfg)
	if err != nil {
		log.Fatalf("❌ mongo connect error: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	database := db.GetDatabase(client, cfg)

	// === 2. Setup domain & service ===
	repo := dbout.NewMongoPersonRepository(database)
	svc := app.NewService(repo)

	// === 3. Setup Fiber (REST API) ===
	appFiber := fiber.New()
	httpadapter.NewHandler(svc).Register(appFiber)

	// === 4. Setup gRPC server ===
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("❌ failed to listen on %s: %v", grpcPort, err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterCourseServiceServer(grpcServer, grpcadapter.NewCourseGRPCServer(svc))
		reflection.Register(grpcServer) // enable reflection for grpcurl testing

		log.Printf("🚀 gRPC CourseService running on :%s\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ gRPC server stopped: %v", err)
		}
	}()

	// === 5. Start Fiber server ===
	go func() {
		log.Printf("🌐 REST listening on :%s\n", restPort)
		if err := appFiber.Listen(":" + restPort); err != nil {
			log.Println("REST server stopped:", err)
		}
	}()

	// === 6. Wait for shutdown signal ===
	waitForSignal()
	log.Println("🛑 shutting down servers...")
	_ = appFiber.Shutdown()
}

// graceful shutdown
func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
