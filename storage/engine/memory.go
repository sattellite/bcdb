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
	ErrNotFound = errors.New("not found")
)

func NewMemory(l *slog.Logger, done chan struct{}) (*Memory, error) {
	if l == nil {
		return nil, errors.New("logger is required")
	}

	if done == nil {
		return nil, errors.New("done channel is required")
	}

	return &Memory{
		done:   done,
		store:  make(map[string]any),
		logger: l.With("engine", "memory"),
	}, nil
}

type Memory struct {
	done   chan struct{}
	store  map[string]any
	logger *slog.Logger
}

func (m *Memory) Set(ctx context.Context, key string, value any) (err error) {
	defer func(start time.Time) {
		err = m.deferredLog("set", key, start, err)
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
		err = m.deferredLog("get", key, start, err)
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
		return nil, ErrNotFound
	}

	return value, nil
}

func (m *Memory) Del(ctx context.Context, key string) (err error) {
	defer func(start time.Time) {
		err = m.deferredLog("del", key, start, err)
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

	if _, ok := m.store[key]; !ok {
		return ErrNotFound
	}

	delete(m.store, key)
	return nil
}

func (m *Memory) deferredLog(method, key string, start time.Time, err error) error {
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

func (m *Memory) Done() <-chan struct{} {
	return m.done
}

func (m *Memory) Close(_ context.Context) {
	m.logger.Info("closing")
	select {
	case _, ok := <-m.done:
		if ok {
			m.done <- struct{}{}
			close(m.done)
		}
	default:
		m.logger.Warn("already closed")
	}
}
