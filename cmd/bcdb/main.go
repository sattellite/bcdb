package main

import (
	"log/slog"

	"github.com/sattellite/bcdb/config"
	"github.com/sattellite/bcdb/logger"
)

func main() {
	log := logger.Default()
	// load app config
	cfg, cfgErr := config.Load()
	if cfgErr != nil {
		log.Error("failed to load config: ", cfgErr)
	}
	// set logger config
	log = logger.SetDefault(logger.WithConfig(cfg))

	log.Info("starting bcdb")
	log.Debug("loaded config ", slog.Any("cfg", cfg))
}
