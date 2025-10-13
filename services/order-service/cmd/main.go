package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	httpadapter "github.com/wnmay/horo/services/order-service/internal/adapters/inbound/http"
	dbout "github.com/wnmay/horo/services/order-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/order-service/internal/app"
	"github.com/wnmay/horo/shared/config"
	"github.com/wnmay/horo/shared/db"
	"github.com/wnmay/horo/shared/env"
)

func main() {
	_ = config.LoadEnv("order-service")
	port := env.GetString("REST_PORT", "3002")

	dbName := env.GetString("DB_NAME", "orderdb")
	cfg := db.NewMongoDefaultConfig(dbName)
	client, err := db.NewMongoClient(context.Background(), cfg)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	database := db.GetDatabase(client, cfg)

	repo := dbout.NewMongoPersonRepository(database)
	svc := app.NewService(repo)

	appFiber := fiber.New()
	httpadapter.NewHandler(svc).Register(appFiber)

	go func() {
		log.Println("REST listening on :" + port)
		if err := appFiber.Listen(":" + port); err != nil {
			log.Println("server stopped:", err)
		}
	}()

	waitForSignal()
	_ = appFiber.Shutdown()
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
