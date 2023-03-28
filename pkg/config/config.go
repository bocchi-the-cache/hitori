package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var Cfg Config

// Init Read config from file and parse. If file not exist, create it.
func Init(configFile string) {
	if configFile == "" {
		configFile = "hitori.yaml"
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Panic("Config file not exist, exit...")
	}

	f, err := os.ReadFile(configFile)
	if err != nil {
		log.Panic("Read config file failed: ", err)
	}

	if err = yaml.Unmarshal(f, &Cfg); err != nil {
		log.Panic("Parse config file failed: ", err)
	}
	log.Println("config file parsed successfully: ", configFile)
	return
}
