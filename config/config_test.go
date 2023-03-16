package config

import (
	"github.com/anqur/yasch"
	"testing"
)

func TestConfigInitWithPath(t *testing.T) {
	Init("example/config.yaml")
	t.Logf("Config: %+v", Cfg)
	if Cfg.Server.Port != 8080 {
		t.Errorf("Port not match")
	}
}

func TestGenerateJsonSchema(t *testing.T) {
	yasch.WriteFile(Cfg, "config.schema.json")
}
