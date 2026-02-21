package main

import (
	"log"
	"log/slog"
	"os"
	"queue/internal/app"
	"queue/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("service", "my-service")
	slog.SetDefault(logger)
	slog.Info("Начало работы")
	cfg := config.New()
	err = app.Start(cfg)
	if err != nil {
		log.Fatal(err)
	}

}
