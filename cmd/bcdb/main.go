package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sattellite/bcdb/compute"

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
		os.Exit(1)
	}

	// set logger config
	log = logger.SetDefault(logger.WithConfig(cfg))

	log.Info("starting bcdb")
	log.Debug("loaded config ", slog.Any("cfg", cfg))

	ctx, cancel := context.WithCancel(context.Background())

	// wait for signals
	wait := make(chan os.Signal, 1)
	signal.Notify(
		wait,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP)

	server(ctx, cancel, log, cfg, wait)
}

func server(ctx context.Context, cancel context.CancelFunc, log *slog.Logger, cfg *config.Config, wait chan os.Signal) {
	// create storage engine
	eng, engineErr := storage.NewEngine(ctx, storage.EngineTypeMemory)
	if engineErr != nil {
		log.Error("failed to create storage engine", slog.Any("error", engineErr))
		cancel()
		return
	}
	// create computer for user requests
	comp, cErr := compute.New(eng, cfg)
	if cErr != nil {
		log.Error("failed to create compute", slog.Any("error", cErr))
		cancel()
		return
	}
	go comp.Run(ctx)

	<-wait
	// send cancel signal
	cancel()
	log.Info("stopping bcdb")
	<-eng.Done()
}
