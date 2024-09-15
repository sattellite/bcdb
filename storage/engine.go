package storage

import (
	"context"
	"log/slog"

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
}

func NewEngine(t EngineType) Engine {
	l := logger.WithScope("storage")
	l.Info("creating storage engine", slog.String("type", t.String()))
	switch t {
	case EngineTypeMemory:
		return engine.NewMemory(l)
	}
	return nil
}
