package main

import (
	"log"
	"log/slog"
	"os"
	"queue/internal/app"
	"queue/internal/config"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}),
	).With(
		"service", "my-service")
	slog.SetDefault(logger)
	slog.Info("Начало работы")
	cfg := config.New()
	err := app.Start(cfg)
	if err != nil {
		log.Fatal(err)
	}

}
