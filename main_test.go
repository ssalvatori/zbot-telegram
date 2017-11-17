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
	os.Setenv("ZBOT_DATABASE", "new_database.sql")

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

func TestSetupDatabaseSqLite(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "sqlite")
	os.Setenv("ZBOT_DATABASE", "hola.sql")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "sqlite", "DataBaseType sqlite OK")
	assert.EqualValues(t, "DB: hola.sql", zbot.Db.GetConnectionInfo(), "DataBaseType sqlite OK")

}

func TestSetupDatabaseMysql(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "mysql")
	os.Setenv("ZBOT_MYSQL_HOSTNAME", "localhost")
	os.Setenv("ZBOT_MYSQL_USERNAME", "root")
	os.Setenv("ZBOT_MYSQL_PASSWORD", "pass")
	os.Setenv("ZBOT_MYSQL_DATABASE", "test")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "mysql", "DataBaseType mysql OK")
	assert.EqualValues(t, "root:pass@localhost/test", zbot.Db.GetConnectionInfo(), "DataBaseType mysql OK")

}
