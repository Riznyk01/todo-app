package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	todoapp "todo-app"
	"todo-app/pkg/handler"
	"todo-app/repository"
	"todo-app/service"
)

func main() {
	log := setupLogger(os.Getenv("TYPE_OF_LOG"), os.Getenv("LOG_LEVEL"))

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_Host"),
		Port:     os.Getenv("DB_Port"),
		Username: os.Getenv("DB_Username"),
		Password: os.Getenv("DB_Password"),
		DBName:   os.Getenv("DB_Name"),
		SSLMode:  os.Getenv("DB_SSLMode"),
	})
	if err != nil {
		log.Error(fmt.Sprintf("failed to initialize db: %s", err.Error()))
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todoapp.Server)
	log.Info(fmt.Sprintf("server starting on %s:%s", os.Getenv("APP_ADDR"), os.Getenv("APP_PORT")))
	if err := srv.Run(os.Getenv("APP_ADDR"), os.Getenv("APP_PORT"), handlers.InitRouts()); err != nil {
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
