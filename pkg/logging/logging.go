package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	level := os.Getenv("BGS_LOG_LEVEL")
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info", "":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
