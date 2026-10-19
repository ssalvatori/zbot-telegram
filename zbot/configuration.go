package zbot

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ssalvatori/zbot-telegram/db"
	"go.yaml.in/yaml/v3"
)

// Configuration contains the settings loaded from the bot configuration file.
type Configuration struct {
	Zbot    ConfigurationZbot    `yaml:"zbot"`
	Db      db.Config            `yaml:"db"`
	Modules ConfigurationModules `yaml:"modules"`
}

// ConfigurationZbot contains the bot's token and feature settings.
type ConfigurationZbot struct {
	Token          string   `yaml:"token"`
	IgnoreDuration int      `yaml:"ignore_duration"`
	Ignore         bool     `yaml:"ignore"`
	Level          bool     `yaml:"level"`
	Admins         []string `yaml:"admins"`
}

// ConfigurationModules contains external module settings.
type ConfigurationModules struct {
	Path string                `yaml:"path"`
	List []ConfigurationModule `yaml:"list"`
}

// ConfigurationModule describes an external module.
type ConfigurationModule struct {
	Key         string `yaml:"key"`
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

// LoadConfiguration reads and parses a bot configuration file.
func LoadConfiguration(filename string) (*Configuration, error) {
	slog.Info("Reading file", "filename", filename)
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Configuration{}
	if err := yaml.Unmarshal(buf, c); err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
