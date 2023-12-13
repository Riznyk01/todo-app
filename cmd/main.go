package main

import (
	"fmt"
	"log/slog"
	"os"
	"todo-app/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(os.Getenv("TYPE_OF_LOG"), os.Getenv("LOG_LEVEL"))
	log.Info("starting app")
	log.Debug("debug messages are enabled")
	fmt.Println(cfg)
}
func setupLogger(typeOfLog, level string) *slog.Logger {
	var log *slog.Logger
	opts := slog.HandlerOptions{}

	switch level {
	case "DEBUG":
		opts = slog.HandlerOptions{Level: slog.LevelDebug}
	case "INFO":
		opts = slog.HandlerOptions{Level: slog.LevelInfo}
	default:
		opts = slog.HandlerOptions{Level: slog.LevelDebug}
	}

	switch typeOfLog {
	case "TEXTLOG":
		log = slog.New(slog.NewTextHandler(os.Stdout, &opts))
	case "JSONLOG":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	}
	return log
}
