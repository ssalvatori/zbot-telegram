package main

import (
	"fmt"
	"log/slog"
	"os"

	"go.yaml.in/yaml/v4"
)

// Configuration bot configuration
type Configuration struct {
	Zbot     configurationZbot     `yaml:"zbot"`
	Db       configurationDb       `yaml:"db"`
	Webhook  configurationWebhook  `yaml:"webhook"`
	Commands configurationCommands `yaml:"commands"`
	Modules  configurationModules  `yaml:"modules"`
}

type configurationWebhook struct {
	Disable bool      `yaml:"disable,omitempty"`
	Port    int       `yaml:"port"`
	Auth    []channel `yaml:"auth"`
}

type channel struct {
	Channel string `yaml:"channel"`
	ID      int64  `yaml:"id,omitempty"`
	Token   string `yaml:"token,omitempty"`
}

type configurationCommands struct {
	Learn    configurationLearn `yaml:"learn"`
	Disabled []string           `yaml:"disabled"`
}

type configurationZbot struct {
	Token          string `yaml:"token"`
	IgnoreDuration int    `yaml:"ignore_duration"`
	Ignore         bool   `yaml:"ignore"`
	Level          bool   `yaml:"level"`
}

type configurationDb struct {
	Engine   string `yaml:"engine"`
	Name     string `yaml:"name"`
	File     string `yaml:"file"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type configurationLearn struct {
	Disabled []string `yaml:"disabled"`
}

type configurationModules struct {
	Path string                `yaml:"path"`
	List []configurationModule `yaml:"list"`
}

type configurationModule struct {
	Key         string `yaml:"key"`
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

func readConfiguration(filename string) (*Configuration, error) {

	slog.Info("Reading file", "filename", filename)
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Configuration{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
