package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())

	// create storage engine
	eng := storage.NewEngine(ctx, storage.EngineTypeMemory)

	runtime.KeepAlive(eng)

	// wait for signals
	wait := make(chan os.Signal, 1)
	signal.Notify(
		wait,
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP)

	<-wait
	// send cancel signal
	cancel()
	log.Info("stopping bcdb")
	<-eng.Done()
}
