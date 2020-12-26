package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

/*
type Configuration struct {
	Zbot struct {
		Token          string `yaml:"token"`
		IgnoreDuration int    `yaml:"ignore_duration"`
		Ignore         bool   `yaml:"ignore"`
		Level          bool   `yaml:"level"`
	} `yaml:"zbot"`
	Db struct {
		Engine   string `yaml:"engine"`
		Name     string `yaml:"name"`
		File     string `yaml:"file"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"db"`
	Commands struct {
		Learn struct {
			Disabled []string `yaml:"disabled"`
		} `yaml:"learn"`
		Disabled []string `yaml:"disabled"`
	} `yaml:"commands"`
	Modules struct {
		Path string `yaml:"path"`
		List []struct {
			Key         string `yaml:"key"`
			File        string `yaml:"file"`
			Description string `yaml:"description"`
		} `yaml:"list"`
	} `yaml:"modules"`
}

*/

//Configuration bot configuration
type Configuration struct {
	Zbot     configurationZbot     `yaml:"zbot"`
	Db       configurationDb       `yaml:"db"`
	Commands configurationCommands `yaml:"commands"`
	Modules  configurationModules  `yaml:"modules"`
}

type configurationCommands struct {
	Learn    configurationLearn `yaml:"learn"`
	Disabled []string           `yaml:"disabled"`
}
type configurationDisabledList struct {
	List []string `yaml:"disabled"`
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
	Path string `yaml:"path"`
	List []struct {
		Key         string `yaml:"key"`
		File        string `yaml:"file"`
		Description string `yaml:"description"`
	} `yaml:"list"`
}

type configurationModule struct {
	Key         string `yaml:"key"`
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

func readConfiguration(filename string) (*Configuration, error) {

	log.Info("Reading file " + filename)
	buf, err := ioutil.ReadFile(filename)
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
