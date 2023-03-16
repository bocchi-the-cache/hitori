package main

import (
	"github.com/bocchi-the-cache/hitori/pkg/config"
	"github.com/bocchi-the-cache/hitori/pkg/logger"
	"github.com/bocchi-the-cache/hitori/pkg/origin"
	"github.com/bocchi-the-cache/hitori/pkg/proxy"
)

func main() {
	logger.Init("log")
	config.Init("conf/hitori.yaml")
	origin.Init(&config.Cfg.Mapping)
	proxy.Init(&config.Cfg)

	if err := proxy.Serve(); err != nil {
		panic(err)
	}
}
