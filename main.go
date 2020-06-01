package main

import (
	"os"
	"strings"

	"fmt"

	env "github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/zbot"
)

func setDisabledCommands(cmds string) []string {

	if len(cmds) == 0 {
		return []string{}
	}

	log.Debug("Creating disabled commands list")
	cmdList := strings.Split(cmds, ",")

	for i := range cmdList {
		cmdList[i] = strings.TrimSpace(cmdList[i])
	}

	return cmdList
}

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

// setupDatabase this function will get the data
func setupDatabase() db.ZbotDatabase {

	var db db.ZbotDatabase

	switch os.Getenv("ZBOT_DATABASE_TYPE") {
	case "mysql":
		log.Info("Setting up mysql connections")
		db = setupDatabaseMysql()
		break
	case "sqlite3":
		log.Info("Setting up sqlite connections")
		db = setupDatabaseSqlite3()
		break
	default:
		log.Fatal("Select a database type (mysql o sqlite3)")
		break
	}
	return db

}

func setupDatabaseSqlite3() db.ZbotDatabase {
	zbot.DatabaseType = "sqlite3"

	type SqliteEnvironmentConfig struct {
		File string `env:"ZBOT_SQLITE3_DATABASE"`
	}

	cfg := SqliteEnvironmentConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(fmt.Printf("%+v\n", err))
	}

	database := new(db.ZbotSqlite3Database)
	database.File = cfg.File
	return database
}

func setupDatabaseMysql() db.ZbotDatabase {
	zbot.DatabaseType = "mysql"

	type MysqlConnection struct {
		Username     string `env:"ZBOT_MYSQL_USERNAME,required"`
		Password     string `env:"ZBOT_MYSQL_PASSWORD,required"`
		DatabaseName string `env:"ZBOT_MYSQL_DATABASE,required"`
		HostName     string `env:"ZBOT_MYSQL_HOSTNAME,required"`
		Protocol     string `env:"ZBOT_MYSQL_PROTOCOL" envDefault:"tcp"`
		Port         int    `env:"ZBOT_MYSQL_PORT" envDefault:"3306"`
	}

	cfg := MysqlConnection{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(fmt.Printf("%+v\n", err))
	}

	database := new(db.ZbotMysqlDatabase)
	database.Connection = db.MysqlConnection(cfg)

	return database
}

//setupFlags Set flags configurations
func setupFlags() {
	var ok bool
	_, ok = os.LookupEnv("ZBOT_FLAG_ACTIVATE_IGNORE")

	if ok {
		zbot.Flags.Ignore = true
	}

	_, ok = os.LookupEnv("ZBOT_FLAG_ACTIVATE_LEVELS")
	if ok {
		zbot.Flags.Level = true
	}

}

func setup() {

	type EnvironmentVariables struct {
		Token            string `env:"ZBOT_TOKEN,required"`
		ModulesPath      string `env:"ZBOT_MODULES_PATH" envDefault:"."`
		DisabledCommands string `env:"ZBOT_DISABLED_COMMANDS" `
	}

	cfg := EnvironmentVariables{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(fmt.Sprintf("%+v\n", err))
	}

	log.Debug(fmt.Sprintf("%+v\n", cfg))

	zbot.APIToken = cfg.Token
	zbot.ModulesPath = cfg.ModulesPath + "/"

	if os.Getenv("ZBOT_DISABLED_COMMANDS") != "" {
		zbot.SetDisabledCommands(setDisabledCommands(os.Getenv("ZBOT_DISABLED_COMMANDS")))
	}

}

func init() {
	setupLog()
	setupFlags()
}

func main() {
	setup()
	zbot.Db = setupDatabase()
	zbot.Execute()
}
