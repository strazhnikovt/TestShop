package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/strazhnikovt/TestShop/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	go func() {
		if err := application.Run(); err != nil {
			application.Logger.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	application.Logger.Printf("Shutting down server...")
	if err := application.Shutdown(); err != nil {
		application.Logger.Fatalf("Server shutdown error: %v", err)
	}
	application.Logger.Printf("Server exited")
}
