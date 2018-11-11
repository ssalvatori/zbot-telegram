package main

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/zbot"
	"github.com/stretchr/testify/assert"
)

func TestSetUp(t *testing.T) {

	os.Setenv("ZBOT_TOKEN", "test:test")
	os.Setenv("ZBOT_SQLITE_DATABASE", "new_database.sql")

	dir, _ := os.Getwd()

	os.Setenv("ZBOT_MODULES_PATH", dir)

	setUp()
	assert.Equal(t, dir+"/", zbot.ModulesPath, "Setting module path")

	os.Setenv("ZBOT_MODULES_PATH", "/tmp")
	setUp()
	assert.Equal(t, "/tmp/", zbot.ModulesPath, "Setting module path")
}

func TestSetUpLog(t *testing.T) {

	levels := map[string]string{
		"info":  "info",
		"debug": "debug",
		"panic": "panic",
		"error": "error",
		"warn":  "warning",
		"fatal": "fatal",
		"":      "info",
	}
	for key, value := range levels {
		os.Setenv("ZBOT_LOG_LEVEL", key)
		setUpLog()
		assert.Equal(t, log.GetLevel().String(), value, key+"OK")
	}

}

func TestSetupNoDatabase(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "", "DataBaseType empty OK")
}

func TestSetupDatabaseSqLite(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "sqlite")
	os.Setenv("ZBOT_SQLITE_DATABASE", "hola.sql")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "sqlite", "DataBaseType sqlite OK")
	assert.Equal(t, "DB: hola.sql", zbot.Db.GetConnectionInfo(), "DataBaseType sqlite OK")

}

func TestSetupDatabaseMysql(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "mysql")
	os.Setenv("ZBOT_MYSQL_HOSTNAME", "localhost")
	os.Setenv("ZBOT_MYSQL_USERNAME", "root")
	os.Setenv("ZBOT_MYSQL_PASSWORD", "pass")
	os.Setenv("ZBOT_MYSQL_DATABASE", "test")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "mysql", "DataBaseType mysql OK")
	assert.Equal(t, "root:pass@localhost:3306/test", zbot.Db.GetConnectionInfo(), "DataBaseType mysql OK")

}

func TestGetDisabledCommandsError(t *testing.T) {
	assert.Error(t, getDisabledCommands("lala.json"), "", "Disabled Command Error")
}

func TestGetDisabledCommands(t *testing.T) {
	assert.Equal(t, getDisabledCommands("main_test.go"), nil, "Set Disabled Commands")
}
