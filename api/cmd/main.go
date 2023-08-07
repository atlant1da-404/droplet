// Package main creates and runs application instance.
package main

import (
	"github.com/atlant1da-404/droplet/app"
	"github.com/atlant1da-404/droplet/config"
	"github.com/atlant1da-404/droplet/pkg/logger"
)

func main() {
	log := logger.New("main")

	cfg := config.Get()
	log.Info("read config", "config", cfg)

	app.Run(cfg)
}
