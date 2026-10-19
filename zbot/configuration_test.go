package zbot

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validConfig = `
zbot:
  token: <TELEGRAM_TOKEN>
  ignore_duration: 300
  ignore: true
  level: false
  admins:
  - admin_username_here
db:
  engine: sqlite
  name: db_name
  file: path_to_sqlite_file.db
  host: 127.0.0.1
  port: 3306
  username: db_username
  password: db_password
modules:
  path: ./modules/
  list:
    - key: crypto
      file: cypto
      description: get some crypto data
    - key: test
      file: test
      description: test module
`

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "zbot.conf")
	err := os.WriteFile(path, []byte(content), 0o600)
	assert.NoError(t, err)

	return path
}

func TestLoadConfigurationValid(t *testing.T) {
	path := writeTempConfig(t, validConfig)

	conf, err := LoadConfiguration(path)
	assert.NoError(t, err)

	assert.Equal(t, "<TELEGRAM_TOKEN>", conf.Zbot.Token)
	assert.Equal(t, 300, conf.Zbot.IgnoreDuration)
	assert.True(t, conf.Zbot.Ignore)
	assert.False(t, conf.Zbot.Level)

	assert.Equal(t, "sqlite", conf.Db.Engine)
	assert.Equal(t, "db_name", conf.Db.Name)
	assert.Equal(t, "path_to_sqlite_file.db", conf.Db.File)
	assert.Equal(t, "127.0.0.1", conf.Db.Host)
	assert.Equal(t, 3306, conf.Db.Port)
	assert.Equal(t, "db_username", conf.Db.Username)
	assert.Equal(t, "db_password", conf.Db.Password)

	assert.Equal(t, []string{"admin_username_here"}, conf.Zbot.Admins)

	assert.Equal(t, "./modules/", conf.Modules.Path)
	assert.Equal(t, []ConfigurationModule{
		{Key: "crypto", File: "cypto", Description: "get some crypto data"},
		{Key: "test", File: "test", Description: "test module"},
	}, conf.Modules.List)
}

func TestLoadConfigurationFileNotFound(t *testing.T) {
	conf, err := LoadConfiguration("/nonexistent/path/zbot.conf")
	assert.Nil(t, conf)
	assert.Error(t, err)
}

func TestLoadConfigurationInvalidYAML(t *testing.T) {
	path := writeTempConfig(t, "not: valid: yaml: ::::\n")

	conf, err := LoadConfiguration(path)
	assert.Nil(t, conf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), path)
}

func TestLoadConfigurationEmptyFile(t *testing.T) {
	path := writeTempConfig(t, "")

	conf, err := LoadConfiguration(path)
	assert.NoError(t, err)
	assert.Equal(t, &Configuration{}, conf)
}
