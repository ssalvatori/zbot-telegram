package configuration

import (
	"fmt"
	"log/slog"
	"os"

	"go.yaml.in/yaml/v4"
)

// Configuration bot configuration
type Configuration struct {
	Zbot     ConfigurationZbot     `yaml:"zbot"`
	Db       ConfigurationDb       `yaml:"db"`
	Commands ConfigurationCommands `yaml:"commands"`
	Modules  ConfigurationModules  `yaml:"modules"`
}

type ConfigurationCommands struct {
	Learn    ConfigurationLearn `yaml:"learn"`
	Disabled []string           `yaml:"disabled"`
}

type ConfigurationZbot struct {
	Token          string `yaml:"token"`
	IgnoreDuration int    `yaml:"ignore_duration"`
	Ignore         bool   `yaml:"ignore"`
	Level          bool   `yaml:"level"`
}

type ConfigurationDb struct {
	Engine   string `yaml:"engine"`
	Name     string `yaml:"name"`
	File     string `yaml:"file"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ConfigurationLearn struct {
	Disabled []string `yaml:"disabled"`
}

type ConfigurationModules struct {
	Path string                `yaml:"path"`
	List []ConfigurationModule `yaml:"list"`
}

type ConfigurationModule struct {
	Key         string `yaml:"key"`
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

func LoadConfiguration(filename string) (*Configuration, error) {

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
