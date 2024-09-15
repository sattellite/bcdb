package main

import (
	"log/slog"
	"runtime"

	"github.com/sattellite/bcdb/config"
	"github.com/sattellite/bcdb/logger"
	"github.com/sattellite/bcdb/storage"
)

func main() {
	log := logger.Default()

	// load app config
	cfg, cfgErr := config.Load()
	if cfgErr != nil {
		log.Error("failed to load config", slog.Any("error", cfgErr))
	}

	// set logger config
	log = logger.SetDefault(logger.WithConfig(cfg))

	log.Info("starting bcdb")
	log.Debug("loaded config ", slog.Any("cfg", cfg))

	// create storage engine
	eng := storage.NewEngine(storage.EngineTypeMemory)

	runtime.KeepAlive(eng)
}
