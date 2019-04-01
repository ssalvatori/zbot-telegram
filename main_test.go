package main

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/zbot"
	"github.com/stretchr/testify/assert"
)

func Test_setUp(t *testing.T) {

	os.Setenv("ZBOT_TOKEN", "test:test")
	os.Setenv("ZBOT_SQLITE_DATABASE", "new_database.sql")

	dir, _ := os.Getwd()

	os.Setenv("ZBOT_MODULES_PATH", dir)

	setUp()
	assert.Equal(t, dir+"/", zbot.ModulesPath, "Setting module path")

	os.Setenv("ZBOT_MODULES_PATH", "/tmp")
	setUp()
	assert.Equal(t, "/tmp/", zbot.ModulesPath, "Setting module path")

	os.Setenv("ZBOT_DISABLED_COMMANDS", "lala.json")
	setUp()
	assert.Equal(t, []string(nil), zbot.GetDisabledCommands(), "Setting DisableCommands")
}

func Test_setUpLog(t *testing.T) {

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

func Test_setupNoDatabase(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "", "DataBaseType empty OK")
}

func Test_setupDatabaseSqLite(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "sqlite")
	os.Setenv("ZBOT_SQLITE_DATABASE", "hola.sql")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "sqlite", "DataBaseType sqlite OK")
	assert.Equal(t, "DB: hola.sql", zbot.Db.GetConnectionInfo(), "DataBaseType sqlite OK")

}

func Test_setupDatabaseMysql(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "mysql")
	os.Setenv("ZBOT_MYSQL_HOSTNAME", "localhost")
	os.Setenv("ZBOT_MYSQL_USERNAME", "root")
	os.Setenv("ZBOT_MYSQL_PASSWORD", "pass")
	os.Setenv("ZBOT_MYSQL_DATABASE", "test")
	setupDatabase()
	assert.Equal(t, zbot.DatabaseType, "mysql", "DataBaseType mysql OK")
	assert.Equal(t, "root:pass@localhost:3306/test", zbot.Db.GetConnectionInfo(), "DataBaseType mysql OK")

}

func Test_setDisabledCommands(t *testing.T) {
	assert.Equal(t, setDisabledCommands("main_test.go"), nil, "Set Disabled Commands")

	assert.Error(t, setDisabledCommands("lala.json"), "", "Disabled Command Error")
}

func Test_setupFlags(t *testing.T) {
	os.Unsetenv("ZBOT_FLAG_IGNORE")
	setupFlags()
	assert.Equal(t, zbot.Flags.Ignore, false, "Ignore Off")

	os.Setenv("ZBOT_FLAG_IGNORE", "true")
	setupFlags()
	assert.Equal(t, zbot.Flags.Ignore, true, "Ignore ON")

	os.Unsetenv("ZBOT_FLAG_LEVEL")
	setupFlags()
	assert.Equal(t, zbot.Flags.Level, false, "Level Off")

	os.Setenv("ZBOT_FLAG_LEVEL", "true")
	setupFlags()
	assert.Equal(t, zbot.Flags.Level, true, "Level ON")
}
