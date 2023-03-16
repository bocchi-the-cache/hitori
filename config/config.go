package config

import (
	"github.com/bocchi-the-cache/hitori/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var Cfg Config

// Init Read config from file and parse. If file not exist, create it.
func Init(configFile string) {
	if configFile == "" {
		configFile = "config.yaml"
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Panic("Config file not exist, exit...")
	}

	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Panic("Read config file failed: ", err)
	}

	err = yaml.Unmarshal(f, &Cfg)
	if err != nil {
		logger.Panic("Parse config file failed: ", err)
	}

	return
}
