package zap

import (
	"context"

	"github.com/LexBokun/ControlAgent/internal/pkg/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Zap struct {
	cfg    logger.Config
	logger *zap.SugaredLogger
}

func NewZapLogger(libCfg logger.Config, presetKeyAndValue ...any) (*Zap, error) {
	// Default is development logger.
	zapCfg := zap.NewDevelopmentConfig()
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapCfg.EncoderConfig = encoderCfg
	// zapCfg.Level = libCfg.Level
	stacktraceLevel := zapcore.ErrorLevel

	zapCfg.Level.Enabled(zapcore.DebugLevel)

	if libCfg.Structured {
		zapCfg = zap.NewProductionConfig()
		stacktraceLevel = zapcore.DPanicLevel
	}

	level, err := zapcore.ParseLevel(libCfg.Level)
	if err != nil {
		return nil, err
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)

	callerSkip := 1 // 1 - because wrapper add one more function.
	logger, err := zapCfg.Build(zap.WithCaller(true), zap.AddCallerSkip(callerSkip), zap.AddStacktrace(stacktraceLevel))
	if err != nil {
		return nil, err
	}

	return &Zap{
		cfg:    libCfg,
		logger: logger.Sugar().With(presetKeyAndValue...),
	}, nil
}

func (z *Zap) Logger() *zap.Logger {
	return z.logger.Desugar()
}

func (z *Zap) With(ctx context.Context, presetKeysVals ...any) (context.Context, logger.Logger) {
	keys := presetKeysVals

	if z.cfg.StoreKeysAtContext {
		ctx, keys = logger.UpsertPresetAtCtx(ctx, presetKeysVals...)
	}

	zapLogger := z.logger.With(keys...)

	return ctx, &Zap{logger: zapLogger}
}

func (z *Zap) Debug(msg string, keyAndValues ...any) {
	z.logger.Debugw(msg, keyAndValues...)
}

func (z *Zap) Info(msg string, keyAndValues ...any) {
	z.logger.Infow(msg, keyAndValues...)
}

func (z *Zap) Warn(msg string, keyAndValues ...any) {
	z.logger.Warnw(msg, keyAndValues...)
}

func (z *Zap) Error(msg string, err error, keyAndValues ...any) {
	keyAndValues = append(keyAndValues, "error", err)
	z.logger.Errorw(msg, keyAndValues...)
}

func (z *Zap) Panic(msg string, keyAndValues ...any) {
	z.logger.Panicw(msg, keyAndValues...)
}

func (z *Zap) Log(lvl logger.Level, msg string, keysValues ...any) {
	zapLvl := zapcore.Level(lvl)

	z.logger.Logw(zapLvl, msg, keysValues...)
}

func (z *Zap) Level() logger.Level {
	lvl := z.logger.Level()

	return logger.Level(lvl)
}

func (z *Zap) Sync() error {
	return z.logger.Sync()
}
