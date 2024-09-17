package storage

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/sattellite/bcdb/logger"
	"github.com/sattellite/bcdb/storage/engine"
)

type EngineType int

const (
	EngineTypeMemory EngineType = iota
)

func (t *EngineType) String() string {
	if *t == EngineTypeMemory {
		return "memory"
	}
	return "unknown"
}

type Engine interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
	Del(ctx context.Context, key string) error

	Done() <-chan struct{}
	Close(ctx context.Context)
}

func NewEngine(ctx context.Context, t EngineType) (Engine, error) {
	l := logger.WithScope("storage")
	l.Info("creating storage engine", slog.String("type", t.String()))
	var eng Engine
	var err error
	done := make(chan struct{})
	if t == EngineTypeMemory {
		eng, err = engine.NewMemory(l, done)
	}
	if err != nil {
		l.Error("failed to create storage engine", slog.Any("error", err))
		return nil, err
	}

	go stopEngine(ctx, eng)
	return eng, nil
}

func stopEngine(ctx context.Context, eng Engine) {
	<-ctx.Done()
	l := logger.WithScope("storage")
	l.Info("stopping storage engine")
	if eng != nil {
		tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		eng.Close(tctx)
		select {
		case <-eng.Done():
			l.Info("storage engine stopped")
		case <-tctx.Done():
			err := tctx.Err()
			if !errors.Is(err, context.Canceled) {
				l.Error("failed to stop storage engine 1", slog.Any("error", tctx.Err()))
				return
			}
			l.Info("storage engine stopped")
		}
	}
}
