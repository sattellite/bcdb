// Copyright (c) 2024 Alexander <sattellite> Groshev
// SPDX-License-Identifier: MIT

package logger

import (
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/sattellite/bcdb/config"
)

var l = slog.New(slog.NewTextHandler(os.Stdout, defaultOptions()))
var defaultChanged atomic.Bool

func defaultOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339Nano)),
				}
			}
			return a
		},
	}
}

func Default() *slog.Logger {
	return l
}

func SetDefault(dl *slog.Logger) *slog.Logger {
	if defaultChanged.CompareAndSwap(false, true) {
		l = dl
	}
	return l
}

func WithConfig(c *config.Config) *slog.Logger {
	opts := defaultOptions()

	if c != nil && c.Debug {
		opts.Level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func WithScope(scope string) *slog.Logger {
	return l.With("scope", scope)
}
