package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
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

	servCfg, err := todoapp.NewConfig(log)
	if err != nil {
		log.Fatalf("failed to get App config: %s", err.Error())
	}
	postgresCfg, err := todoapp.NewConfigPostgres()
	if err != nil {
		log.Fatalf("failed to get DB config: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(postgresCfg)
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
