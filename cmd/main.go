package main

import (
	"fmt"
	"log/slog"
	"os"
	todoapp "todo-app"
	"todo-app/pkg/handler"
)

func main() {
	log := setupLogger(os.Getenv("TYPE_OF_LOG"), os.Getenv("LOG_LEVEL"))
	handlers := new(handler.Handler)
	srv := new(todoapp.Server)
	log.Info(fmt.Sprintf("server starting on %s:%s", os.Getenv("ADDR"), os.Getenv("PORT")))
	if err := srv.Run(os.Getenv("ADDR"), os.Getenv("PORT"), handlers.InitRouts()); err != nil {
		log.Error("error while starting the HTTP server", err)
	}

}
func setupLogger(typeOfLog, level string) *slog.Logger {
	var log *slog.Logger
	opts := slog.HandlerOptions{}

	switch level {
	case "DEBUG":
		opts = slog.HandlerOptions{Level: slog.LevelDebug}
	case "INFO":
		opts = slog.HandlerOptions{Level: slog.LevelInfo}
	case "WARN":
		opts = slog.HandlerOptions{Level: slog.LevelWarn}
	case "ERROR":
		opts = slog.HandlerOptions{Level: slog.LevelError}
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
