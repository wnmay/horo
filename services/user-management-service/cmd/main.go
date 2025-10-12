package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wnmay/horo/shared/config"
)

func main() {
	_ = config.LoadEnv("user-management-service")
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
