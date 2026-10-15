package main

import (
	"log/slog"
	"os"
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/zbot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetUp(t *testing.T) {

	bot, err := setup("./zbot.conf")
	assert.NoError(t, err)
	assert.NotNil(t, bot)

}

func TestConfigFromFile(t *testing.T) {
	conf, err := zbot.LoadConfiguration("./zbot.conf")
	require.NoError(t, err)
	cfg, err := zbot.LoadConfig(conf)
	assert.NoError(t, err)
	assert.Equal(t, "./modules/", cfg.ModulesPath, "Setting module path")
	assert.Equal(t, "<TELEGRAM_TOKEN>", cfg.Token, "API TOKEN")
	assert.Equal(t, 300, cfg.IgnoreDuration, "IgnoreDuration")
	assert.Equal(t, true, cfg.Flags.Ignore, "Ignore Flags")
	assert.Equal(t, false, cfg.Flags.Level, "Level Flags")
}

func TestSetUpConfigurationNotFound(t *testing.T) {
	bot, err := setup("/nonexistent/path/zbot.conf")
	assert.Error(t, err)
	assert.Nil(t, bot)
}

func TestSetUpNoDatabaseEngine(t *testing.T) {
	f, err := os.CreateTemp("", "zbot_nodb_*.conf")
	require.NoError(t, err)
	_, err = f.WriteString("zbot:\n  token: x\ndb:\n  engine: \"\"\n")
	require.NoError(t, err)
	defer func() { _ = os.Remove(f.Name()) }()
	_ = f.Close()

	bot, err := setup(f.Name())
	assert.ErrorIs(t, err, db.ErrNoEngine)
	assert.Nil(t, bot)
}

func TestSetupLog(t *testing.T) {

	levels := map[string]slog.Level{
		"info":  slog.LevelInfo,
		"debug": slog.LevelDebug,
		"error": slog.LevelError,
		"":      slog.LevelInfo,
	}
	for key, value := range levels {
		setupLog(key)
		assert.Equal(t, logLevel.Level(), value, "Setup log levels for key: "+key)
	}

}

func TestReadConfigurationNotFound(t *testing.T) {
	_, err := zbot.LoadConfiguration("/nonexistent/path/zbot.conf")
	assert.Error(t, err)
}

func TestReadConfigurationInvalidYAML(t *testing.T) {
	f, err := os.CreateTemp("", "zbot_bad_*.conf")
	require.NoError(t, err)
	_, err = f.WriteString("not: valid: yaml: ::::\n")
	require.NoError(t, err)
	defer func() { _ = os.Remove(f.Name()) }()
	_ = f.Close()
	_, err = zbot.LoadConfiguration(f.Name())
	assert.Error(t, err)
}
