package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/inbound/http"
	dbout "github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/payment-service/internal/app"
	"github.com/wnmay/horo/shared/db"
)

func main() {
	port := getenv("REST_PORT", "3001")

	gormDB := db.MustOpen()
	repo := dbout.NewGormPersonRepository(gormDB)
	svc  := app.NewService(repo)

	appFiber := fiber.New()
	http.NewHandler(svc).Register(appFiber)

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

func getenv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }
