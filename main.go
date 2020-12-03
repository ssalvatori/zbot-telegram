package main

import (
	"os"

	"fmt"

	env "github.com/caarlos0/env/v6"
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
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
		break
	}
}

func setupDatabase(conf *Configuration) db.ZbotDatabase {

	var db db.ZbotDatabase

	switch conf.Db.Engine {
	case "mysql":
		log.Info("Setting up mysql connections")
		//db = setupDatabaseMysql()
		break
	case "sqlite":
		log.Info("Setting up sqlite connections")
		db = setupDatabaseSqlite(conf)
		break
	default:
		log.Fatal("Select a database type (mysql o sqlite)")
		break
	}
	return db

}

func setupDatabaseSqlite(conf *Configuration) db.ZbotDatabase {
	zbot.DatabaseType = "sqlite"
	database := new(db.ZbotDatabaseSqlite)
	database.File = conf.Db.File
	return database
}

// func setupDatabaseMysql() db.ZbotDatabase {
// 	zbot.DatabaseType = "mysql"

// 	type MysqlConnection struct {
// 		Username     string `env:"ZBOT_MYSQL_USERNAME,required"`
// 		Password     string `env:"ZBOT_MYSQL_PASSWORD,required"`
// 		DatabaseName string `env:"ZBOT_MYSQL_DATABASE,required"`
// 		HostName     string `env:"ZBOT_MYSQL_HOSTNAME,required"`
// 		Protocol     string `env:"ZBOT_MYSQL_PROTOCOL" envDefault:"tcp"`
// 		Port         int    `env:"ZBOT_MYSQL_PORT" envDefault:"3306"`
// 	}

// 	cfg := MysqlConnection{}
// 	if err := env.Parse(&cfg); err != nil {
// 		log.Fatal(fmt.Printf("%+v\n", err))
// 	}

// 	database := new(db.ZbotMysqlDatabase)
// 	database.Connection = db.MysqlConnection(cfg)

// 	return database
// }

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
	zbot.ModulesPath = configuration.Modules.Path + "/"
	zbot.IgnoreDuration = configuration.Zbot.IgnoreDuration
	zbot.Flags.Ignore = configuration.Zbot.Ignore
	zbot.Flags.Level = configuration.Zbot.Level

	zbot.SetDisabledLearnChannels(configuration.Commands.Learn.Disabled)

	zbot.Db = setupDatabase(configuration)
	zbot.ExternalModules = zbot.ExternalModulesList(configuration.Modules.List)
}

func init() {
	setupLog()
}

func main() {
	setup()
	zbot.Execute()
}
