package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	todoapp "todo-app"
	"todo-app/pkg/handler"
	"todo-app/pkg/repository"
	"todo-app/pkg/service"
)

// @title Todo app API
// @version 1.0
// @description API Server for Todolist application

// @host localhost:8080
// @BasePath/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	log := setupLogger()

	servCfg := todoapp.NewConfig()
	servCfg.AccessTokenTtl = setTokenTTL(log, os.Getenv("ACCESS_TTL"), "30m")
	servCfg.RefreshTokenTtl = setTokenTTL(log, os.Getenv("REFRESH_TTL"), "720h")

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_Host"),
		Port:     os.Getenv("DB_Port"),
		Username: os.Getenv("DB_Username"),
		Password: os.Getenv("DB_Password"),
		DBName:   os.Getenv("DB_Name"),
		SSLMode:  os.Getenv("DB_SSLMode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(log, db)
	services := service.NewService(log, servCfg, repos)
	handlers := handler.NewHandler(log, services)

	srv := new(todoapp.Server)
	log.Infof("server starting on %s:%s", os.Getenv("APP_ADDR"), os.Getenv("APP_PORT"))
	if err := srv.Run(os.Getenv("APP_ADDR"), os.Getenv("APP_PORT"), handlers.InitRouts()); err != nil {
		log.Fatalf("error occured while starting the HTTP server: %s", err)
	}
}
func setupLogger() *logrus.Logger {
	//panic, fatal, error, warn, warning, info, debug, trace
	log := logrus.New()

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)
	setLogType(log)
	return log
}
func setLogType(log *logrus.Logger) {
	switch os.Getenv("TYPE_OF_LOG") {
	case "TEXTLOG":
		log.SetFormatter(&logrus.TextFormatter{})
	case "JSONLOG":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.JSONFormatter{})
	}
}
func setTokenTTL(log *logrus.Logger, t, def string) time.Duration {
	TokenTTL, err := time.ParseDuration(t)
	if err != nil {
		TokenTTL, _ = time.ParseDuration(def)
		log.Errorf("Failed to set access/refresh token expiration duration. Using default value - %v", TokenTTL)
		return TokenTTL
	}
	return TokenTTL
}
