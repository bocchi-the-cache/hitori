package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/bocchi-the-cache/hitori/pkg/logger"
)

var Cfg Config

// Init Read config from file and parse. If file not exist, create it.
func Init(configFile string) {
	if configFile == "" {
		configFile = "hitori.yaml"
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Panic("Config file not exist, exit...")
	}

	f, err := os.ReadFile(configFile)
	if err != nil {
		logger.Panic("Read config file failed: ", err)
	}

	if err = yaml.Unmarshal(f, &Cfg); err != nil {
		logger.Panic("Parse config file failed: ", err)
	}

	return
}
