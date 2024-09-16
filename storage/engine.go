package storage

import (
	"context"
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
	switch *t {
	case EngineTypeMemory:
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

func NewEngine(ctx context.Context, t EngineType) Engine {
	l := logger.WithScope("storage")
	l.Info("creating storage engine", slog.String("type", t.String()))
	var eng Engine
	done := make(chan struct{})
	switch t {
	case EngineTypeMemory:
		eng = engine.NewMemory(l, done)
	}

	go stopEngine(ctx, eng)
	return eng
}

func stopEngine(ctx context.Context, eng Engine) {
	<-ctx.Done()
	l := logger.WithScope("storage")
	l.Info("stopping storage engine")
	if eng != nil {
		tctx, _ := context.WithTimeout(ctx, 5*time.Second)
		eng.Close(tctx)
		select {
		case <-eng.Done():
			l.Info("storage engine stopped")
		case <-tctx.Done():
			l.Error("failed to stop storage engine", slog.Any("error", tctx.Err()))
		}
	}
}
