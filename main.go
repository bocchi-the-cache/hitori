package main

import (
	"github.com/bocchi-the-cache/hitori/internal/proxies"
	"github.com/bocchi-the-cache/hitori/pkg/cache"
	"github.com/bocchi-the-cache/hitori/pkg/config"
	"github.com/bocchi-the-cache/hitori/pkg/logger"
	"github.com/bocchi-the-cache/hitori/pkg/origin"
)

func main() {
	logger.Init("log")
	config.Init("conf/hitori.yaml")
	err := cache.Init(&config.Cfg)
	if err != nil {
		panic(err)
	}
	origin.Init(&config.Cfg.Mapping)
	proxies.Init(&config.Cfg)

	if err := proxies.Serve(); err != nil {
		panic(err)
	}
}
