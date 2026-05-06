package main

import (
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/zbot"
	"github.com/stretchr/testify/assert"
)

var confTest = Configuration{
	Zbot: configurationZbot{
		Token:          "",
		IgnoreDuration: 10,
		Ignore:         false,
		Level:          false,
	},
	Db: configurationDb{
		Engine: "sqlite",
		File:   "",
	},
	Commands: configurationCommands{
		Learn: configurationLearn{
			Disabled: []string{},
		},
		Disabled: []string{},
	},
	Modules: configurationModules{
		Path: "./module",
		List: []configurationModule{
			{
				Key:         "cmd1",
				File:        "cmdFile1",
				Description: "description 1",
			},
			{
				Key:         "cmd2",
				File:        "cmdFile2",
				Description: "description 2",
			},
		},
	},
}

func TestSetUp(t *testing.T) {

	os.Setenv("ZBOT_CONFIG_FILE", "./zbot.conf")

	setup()
	assert.Equal(t, "./modules/", zbot.ModulesPath, "Setting module path")
	assert.Equal(t, "<TELEGRAM_TOKEN>", zbot.APIToken, "API TOKEN")
	assert.Equal(t, 300, zbot.IgnoreDuration, "IgnoreDuration")
	assert.Equal(t, true, zbot.Flags.Ignore, "Ignore Flags")
	assert.Equal(t, false, zbot.Flags.Level, "Level Flags")

}

func TestSetupLog(t *testing.T) {

	levels := map[string]slog.Level{
		"info":  slog.LevelInfo,
		"debug": slog.LevelDebug,
		"error": slog.LevelError,
		"":      slog.LevelInfo,
	}
	for key, value := range levels {
		os.Setenv("ZBOT_LOG_LEVEL", key)
		setupLog()
		assert.Equal(t, logLevel.Level(), value, "Setup log levels for key: "+key)
	}

}

func TestSetupNoDatabase(t *testing.T) {

	// Run the crashing code when FLAG is set
	if os.Getenv("FLAG") == "1" {
		confTest.Db.Engine = ""
		_ = setupDatabase(&confTest)
		return
	}

	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestSetupNoDatabase")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)

	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
}

func TestSetupDatabaseSqLite(t *testing.T) {

	os.Setenv("ZBOT_DATABASE_TYPE", "sqlite")
	os.Setenv("ZBOT_SQLITE_DATABASE", "hola.sql")
	dbInstance := setupDatabase(&confTest)
	assert.Equal(t, zbot.DatabaseType, "sqlite", "DataBaseType sqlite OK")
	assert.IsType(t, &db.ZbotDatabaseSqlite{}, dbInstance)
}

func TestReadConfigurationNotFound(t *testing.T) {
	_, err := readConfiguration("/nonexistent/path/zbot.conf")
	assert.Error(t, err)
}

func TestReadConfigurationInvalidYAML(t *testing.T) {
	f, err := os.CreateTemp("", "zbot_bad_*.conf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("not: valid: yaml: ::::\n")
	f.Close()
	_, err = readConfiguration(f.Name())
	assert.Error(t, err)
}
