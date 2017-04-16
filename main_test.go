package main

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/ssalvatori/zbot-telegram-go/zbot"
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