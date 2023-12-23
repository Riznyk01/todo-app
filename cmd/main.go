package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	todoapp "todo-app"
	"todo-app/pkg/handler"
	"todo-app/pkg/logger"
	"todo-app/pkg/repository"
	"todo-app/pkg/service"
)

func main() {
	logger.SetupLogger()

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

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todoapp.Server)
	log.Infof("server starting on %s:%s", os.Getenv("APP_ADDR"), os.Getenv("APP_PORT"))
	if err := srv.Run(os.Getenv("APP_ADDR"), os.Getenv("APP_PORT"), handlers.InitRouts()); err != nil {
		log.Fatalf("error occured while starting the HTTP server: %s", err)
	}
}
