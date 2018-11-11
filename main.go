package main

import (
	"io/ioutil"
	"os"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/caarlos0/env"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/zbot"
)

// setUp
func setUp() {

	type EnvironmentVariables struct {
		Token            string `env:"ZBOT_TOKEN,required"`
		ModulesPath      string `env:"ZBOT_MODULES_PATH" envDefault:"."`
		DisabledCommands string `env:"ZBOT_DISABLED_COMMANDS" `
	}

	cfg := EnvironmentVariables{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v\n", err))
	}

	log.Debug(fmt.Sprintf("%+v\n", cfg))

	zbot.APIToken = cfg.Token
	zbot.ModulesPath = cfg.ModulesPath + "/"

	if os.Getenv("ZBOT_DISABLED_COMMANDS") != "" {
		log.Info("Disabled modules configuration = ", os.Getenv("ZBOT_DISABLED_COMMANDS"))
		setDisabledCommands(os.Getenv("ZBOT_DISABLED_COMMANDS"))
	}

}

func setDisabledCommands(file string) error {

	log.Debug("Reading file ", file)
	raw, err := ioutil.ReadFile(file)

	if err != nil {
		log.Error(fmt.Sprintf("ccould not file: %v", err))
		return fmt.Errorf("could not file: %v", err)
	}

	zbot.SetDisabledCommands(raw)
	return nil
}

// setUpLog setup log level using environment variables
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

}

// setupDatabase this function will get the data
func setupDatabase() {

	switch os.Getenv("ZBOT_DATABASE_TYPE") {
	case "mysql":
		log.Info("Setting up mysql connections")
		zbot.Db = setupDatabaseMysql()
		break
	case "sqlite":
		log.Info("Setting up sqlite connections")
		zbot.Db = setupDatabaseSqlite()
		break
	default:
		log.Error("Select a database type (mysql o sqlite)")
		break
	}

}

func setupDatabaseSqlite() db.ZbotDatabase {
	zbot.DatabaseType = "sqlite"

	type SqliteEnvironmentConfig struct {
		File string `env:"ZBOT_SQLITE_DATABASE"`
	}

	cfg := SqliteEnvironmentConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(fmt.Printf("%+v\n", err))
	}

	database := new(db.ZbotSqliteDatabase)
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
		Port         int    `env:"ZBOT_MYSQL_PORT" envDefault:"3306"`
	}

	cfg := MysqlConnection{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(fmt.Printf("%+v\n", err))
	}

	database := new(db.ZbotMysqlDatabase)
	database.Connection = db.MysqlConnection(cfg)

	return database
}

func main() {
	setUpLog()
	setUp()
	setupDatabase()

	zbot.Execute()
}
