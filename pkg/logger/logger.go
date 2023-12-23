package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupLogger() {
	//panic, fatal, error, warn, warning, info, debug, trace
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)
	setLogType()
}
func setLogType() {
	switch os.Getenv("TYPE_OF_LOG") {
	case "TEXTLOG":
		log.SetFormatter(&log.TextFormatter{})
	case "JSONLOG":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.JSONFormatter{})
	}
}
