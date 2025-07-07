package logger

import (
	"context"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
)

type Logger interface {
	With(ctx context.Context, presetKeysValues ...any) (context.Context, Logger)
	Debug(msg string, keysValues ...any)
	Info(msg string, keysValues ...any)
	Warn(msg string, keysValues ...any)
	Error(msg string, err error, keysValues ...any)
	Panic(msg string, keysValues ...any)

	Log(lvl Level, msg string, keysValues ...any)

	Level() Level
}

//nolint:gochecknoglobals // Singleton
var defaultLogger Logger = NewNoopLogger()

func GetDefault() Logger {
	return defaultLogger
}

func SetDefault(logger Logger) {
	defaultLogger = logger
}

func With(ctx context.Context, presetKeysValues ...any) (context.Context, Logger) {
	return defaultLogger.With(ctx, presetKeysValues...)
}
