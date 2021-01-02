package main

import (
	"os"

	"fmt"

	env "github.com/caarlos0/env/v6"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/zbot"
)

// setUpLog setup log level using environment variables
func setupLog() {
	log.SetOutput(os.Stdout)

	switch os.Getenv("ZBOT_LOG_LEVEL") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func setupDatabase(conf *Configuration) db.ZbotDatabase {

	var db db.ZbotDatabase

	switch conf.Db.Engine {
	case "mysql":
		log.Info("Setting up mysql connections")
		log.Fatal("Not implemented")
	case "sqlite":
		log.Info("Setting up sqlite connections")
		db = setupDatabaseSqlite(conf)
	default:
		log.Fatal("Select a database type (mysql o sqlite)")
	}
	return db

}

func setupDatabaseSqlite(conf *Configuration) db.ZbotDatabase {
	zbot.DatabaseType = "sqlite"
	database := new(db.ZbotDatabaseSqlite)
	database.File = conf.Db.File
	return database
}

func setup() {

	type EnvironmentVariables struct {
		ConfigurationFile string `env:"ZBOT_CONFIG_FILE" envDefault:"./zbot.conf"`
	}

	cfg := EnvironmentVariables{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(fmt.Sprintf("%+v\n", err))
	}

	log.Debug(fmt.Sprintf("%+v\n", cfg))

	configuration, err := readConfiguration(cfg.ConfigurationFile)
	if err != nil {
		log.Fatal(err)
	}

	zbot.APIToken = configuration.Zbot.Token
	zbot.ModulesPath = configuration.Modules.Path
	zbot.IgnoreDuration = configuration.Zbot.IgnoreDuration
	zbot.Flags.Ignore = configuration.Zbot.Ignore
	zbot.Flags.Level = configuration.Zbot.Level

	zbot.SetDisabledLearnChannels(configuration.Commands.Learn.Disabled)

	zbot.Db = setupDatabase(configuration)

	zbot.ExternalModules = []zbot.ExternalModule{}
	err = mapstructure.Decode(configuration.Modules.List, &zbot.ExternalModules)
	if err != nil {
		log.Fatal(err)
	}

	// zbot.ExternalModules = []zbot.ExternalModule(configuration.Modules.List)
	zbot.Channels = setupChannels(configuration.Webhook.Auth)

	if configuration.Webhook.Disable {
		log.Info("WebServer: disable")
		zbot.Webhook.Enable = false
	} else {
		zbot.Webhook.Enable = true
		log.Info("WebServer: enable")

		if len(configuration.Webhook.Auth) == 0 {
			log.Fatal("No Webhook.Auth present, exiting now!!")
		}

		if configuration.Webhook.Port != 0 {
			zbot.Webhook.Port = configuration.Webhook.Port
		}
		log.Info(fmt.Sprintf("WebServer Port: %d", zbot.Webhook.Port))

	}

}

func setupChannels(configuration []channel) []zbot.Channel {
	var channels = []zbot.Channel{}

	for i := range configuration {
		channels = append(channels, zbot.Channel{
			ID:        configuration[i].ID,
			Title:     configuration[i].Channel,
			AuthToken: configuration[i].Token,
		})
	}

	return channels
}

func init() {
	setupLog()
}

func main() {
	setup()
	zbot.Execute()
}
