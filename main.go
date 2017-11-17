package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/zbot"
)

// setUp
func setUp() {

	if os.Getenv("ZBOT_TOKEN") == "" {
		log.Fatal("You must set the ZBOT_TOKEN environment variable first")
		os.Exit(1)
	}

	if os.Getenv("ZBOT_MODULES_PATH") != "" {
		log.Debug("Module path using: " + os.Getenv("ZBOT_MODULES_PATH"))
		zbot.ModulesPath = os.Getenv("ZBOT_MODULES_PATH") + "/"
	}

	if os.Getenv("ZBOT_DISABLE_COMMANDS") != "" {
		log.Info("Disable modules configuration = ", os.Getenv("ZBOT_DISABLE_COMMANDS"))
		zbot.GetDisabledCommands(os.Getenv("ZBOT_DISABLE_COMMANDS"))
	}

	zbot.ApiToken = os.Getenv("ZBOT_TOKEN")

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

	database := new(db.ZbotSqliteDatabase)

	if os.Getenv("ZBOT_DATABASE") != "" {
		database.File = os.Getenv("ZBOT_DATABASE")
	} else {
		log.Error("Insert the sqlite file name")
	}

	return database
}

func setupDatabaseMysql() db.ZbotDatabase {
	zbot.DatabaseType = "mysql"

	database := new(db.ZbotMysqlDatabase)

	if os.Getenv("ZBOT_MYSQL_HOSTNAME") != "" {
		database.Connection.HostName = os.Getenv("ZBOT_MYSQL_HOSTNAME")
	} else {
		log.Error("Insert the mysql hostname")
	}

	if os.Getenv("ZBOT_MYSQL_USERNAME") != "" {
		database.Connection.Username = os.Getenv("ZBOT_MYSQL_USERNAME")
	} else {
		log.Error("Insert the mysql username")
	}

	if os.Getenv("ZBOT_MYSQL_PASSWORD") != "" {
		database.Connection.Password = os.Getenv("ZBOT_MYSQL_PASSWORD")
	} else {
		log.Error("Insert the mysql password")
	}

	if os.Getenv("ZBOT_MYSQL_DATABASE") != "" {
		database.Connection.DatabaseName = os.Getenv("ZBOT_MYSQL_DATABASE")
	} else {
		log.Error("Insert mysql database name")
	}

	return database
}

func main() {
	setUpLog()
	setUp()

	zbot.Execute()
}
