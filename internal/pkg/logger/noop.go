// default noop logger
package logger

import (
	"context"
	"log/slog"
)

type Noop struct {
	disableWarning bool
}

type OptFn func(*Noop)

func NewNoopLogger(opts ...OptFn) *Noop {
	logger := &Noop{disableWarning: false}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func WithDisableWarning() OptFn {
	return func(n *Noop) {
		n.disableWarning = true
	}
}

func (n *Noop) With(ctx context.Context, presetKeysValues ...any) (context.Context, Logger) {
	return ctx, n
}

func (n *Noop) Debug(msg string, keysValues ...any) {
	if n.disableWarning {
		return
	}

	slog.Warn("logger not init")
}

func (n *Noop) Info(msg string, keysValues ...any) {
	if n.disableWarning {
		return
	}

	slog.Warn("logger not init")
}

func (n *Noop) Warn(msg string, keysValues ...any) {
	if n.disableWarning {
		return
	}

	slog.Warn("logger not init")
}

func (n *Noop) Error(msg string, err error, keysValues ...any) {
	if n.disableWarning {
		return
	}

	slog.Warn("logger not init")
}

func (n *Noop) Panic(msg string, keysValues ...any) {
	if n.disableWarning {
		return
	}

	slog.Warn("logger not init")
}

func (n *Noop) Log(lvl Level, msg string, keysValues ...any) {
	if n.disableWarning {
		return
	}

	slog.Warn("logger not init")
}

func (n *Noop) Level() Level {
	return DebugLevel
}
