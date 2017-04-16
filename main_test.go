package main

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/ssalvatori/zbot-telegram-go/zbot"
	log "github.com/Sirupsen/logrus"
)

func TestSetUp(t *testing.T) {

	os.Setenv("ZBOT_TOKEN", "test:test")
	os.Setenv("ZBOT_DATABASE", "new_database.sql")

	dir,_ := os.Getwd()

	os.Setenv("ZBOT_MODULES_PATH", dir)

	setUp()
	assert.Equal(t,"new_database.sql",zbot.Database, "Setting datbase")
	assert.Equal(t,dir+"/",zbot.ModulesPath, "Setting module path")

	os.Setenv("ZBOT_MODULES_PATH", "/tmp")
	setUp()
	assert.Equal(t,"/tmp/", zbot.ModulesPath, "Setting module path")
}

func TestSetUpLog(t *testing.T) {

	levels := map[string]string {
		"info": "info",
		"debug": "debug",
		"panic": "panic",
		"error": "error",
		"warn": "warning",
		"": "info",
	}
	for key, value := range levels {
		os.Setenv("ZBOT_LOG_LEVEL", key)
		setUpLog()
		assert.Equal(t, log.GetLevel().String(), value, key + "OK")
	}

}