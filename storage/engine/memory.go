package engine

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

var (
	ErrEmptyKey = errors.New("empty key")
	ErrInternal = errors.New("internal error")
)

func NewMemory(l *slog.Logger) *Memory {
	return &Memory{
		store:  make(map[string]any),
		logger: l.With("engine", "memory"),
	}
}

type Memory struct {
	store  map[string]any
	logger *slog.Logger
}

func (m *Memory) Set(ctx context.Context, key string, value any) (err error) {
	defer func(start time.Time) {
		err = m.defferedLog("set", key, start, err)
	}(time.Now())
	m.logger.Debug("set", slog.String("key", key), slog.Any("value", value))

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.set(key, value)
	}
}

func (m *Memory) set(key string, value any) error {
	if key == "" {
		return ErrEmptyKey
	}

	m.store[key] = value
	return nil
}

func (m *Memory) Get(ctx context.Context, key string) (result any, err error) {
	defer func(start time.Time) {
		err = m.defferedLog("get", key, start, err)
	}(time.Now())

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return m.get(key)
	}
}

func (m *Memory) get(key string) (any, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}

	value, ok := m.store[key]
	if !ok {
		return nil, nil
	}

	return value, nil
}

func (m *Memory) Del(ctx context.Context, key string) (err error) {
	defer func(start time.Time) {
		err = m.defferedLog("del", key, start, err)
	}(time.Now())

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.del(key)
	}
}

func (m *Memory) del(key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	delete(m.store, key)
	return nil
}

func (m *Memory) defferedLog(method, key string, start time.Time, err error) error {
	if rErr := recover(); rErr != nil {
		m.logger.Error(method, slog.String("key", key), slog.Any("error", rErr), slog.Duration("elapsed", time.Since(start)))
		return ErrInternal
	}
	if err != nil {
		m.logger.Error(method, slog.String("key", key), slog.Any("error", err), slog.Duration("elapsed", time.Since(start)))
		return err
	}
	m.logger.Debug(method, slog.String("key", key), slog.Duration("elapsed", time.Since(start)))
	return nil
}
