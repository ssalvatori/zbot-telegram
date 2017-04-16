package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"github.com/ssalvatori/zbot-telegram-go/zbot"
)

func setUp() {
	if os.Getenv("ZBOT_DATABASE") != "" {
		zbot.Database = os.Getenv("ZBOT_DATABASE")
	}

	if os.Getenv("ZBOT_TOKEN") == "" {
		log.Fatal("You must set the ZBOT_TOKEN environment variable first")
		os.Exit(1)
	}

	if os.Getenv("ZBOT_MODULES_PATH") != "" {
		log.Debug("Module path using: " + os.Getenv("ZBOT_MODULES_PATH"))
		zbot.ModulesPath = os.Getenv("ZBOT_MODULES_PATH") + "/"
	}

	zbot.ApiToken = os.Getenv("ZBOT_TOKEN")


}

func setUpLog() {

	switch os.Getenv("ZBOT_LOG_LEVEL") {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
		break
	}

	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	setUp()
	setUpLog()

	zbot.Execute()
}
