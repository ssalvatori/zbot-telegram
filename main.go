package main

import (
	"fmt"
	"log/slog"
	"os"

	env "github.com/caarlos0/env/v11"
	"github.com/go-viper/mapstructure/v2"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/zbot"
)

// logLevel holds the current log level and can be inspected in tests
var logLevel = new(slog.LevelVar)

// setupLog configures the default slog logger based on ZBOT_LOG_LEVEL
func setupLog() {
	switch os.Getenv("ZBOT_LOG_LEVEL") {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		logLevel.Set(slog.LevelInfo)
	}
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(h))
}

func setupDatabase(conf *Configuration) db.ZbotDatabase {

	var db db.ZbotDatabase

	switch conf.Db.Engine {
	case "mysql":
		slog.Info("Setting up mysql connections")
		slog.Error("mysql not implemented")
		os.Exit(1)
	case "sqlite":
		slog.Info("Setting up sqlite connections")
		db = setupDatabaseSqlite(conf)
	default:
		slog.Error("no database type selected")
		os.Exit(1)
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
		slog.Error("env parse error", "err", err)
		os.Exit(1)
	}

	slog.Debug("config", "cfg", fmt.Sprintf("%+v", cfg))

	configuration, err := readConfiguration(cfg.ConfigurationFile)
	if err != nil {
		slog.Error("read configuration error", "err", err)
		os.Exit(1)
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
		slog.Error("decode modules error", "err", err)
		os.Exit(1)
	}

}

func init() {
	setupLog()
}

func main() {
	setup()
	zbot.Execute()
}
