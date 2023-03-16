//go:build tools

package main

import (
	"path/filepath"

	"github.com/anqur/yasch"

	"github.com/bocchi-the-cache/hitori/pkg/config"
)

func main() {
	yasch.WriteFile(&config.Cfg, filepath.Join("pkg", "config", "config.schema.json"))
}
