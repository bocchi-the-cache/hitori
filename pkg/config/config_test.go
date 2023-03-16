package config

import (
	"testing"
)

func TestConfigInitWithPath(t *testing.T) {
	Init("../../config/hitori.yaml")
	t.Logf("Config: %+v", Cfg)
	if Cfg.Server.Port != 8080 {
		t.Errorf("Port not match")
	}
}
