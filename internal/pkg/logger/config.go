package logger

import (
	"context"
)

type keyValuesKey string

const kvKey keyValuesKey = "logger-kv-map"

type Config struct {
	Level              string `validate:"required"`
	Structured         bool
	StoreKeysAtContext bool
}

func UpsertPresetAtCtx(ctx context.Context, presetKeysVals ...any) (context.Context, []any) {
	newKV := make([]any, len(presetKeysVals))
	copy(newKV, presetKeysVals)

	ctxKV, ok := ctx.Value(kvKey).(map[any]any)
	if !ok {
		ctxKV = make(map[any]any)
	}

	for k, v := range ctxKV {
		//nolint:makezero // non-zero initialized length slice
		newKV = append(newKV, k, v)
	}

	ctxKV = savePreset(ctxKV, presetKeysVals...)

	ctx = context.WithValue(ctx, kvKey, ctxKV)

	return ctx, newKV
}

func savePreset(kvMap map[any]any, presetKV ...any) map[any]any {
	for i := range len(presetKV) / 2 {
		keyIdx := 2 * i
		valueIdx := 2*i + 1

		key := presetKV[keyIdx]
		value := presetKV[valueIdx]

		kvMap[key] = value
	}

	return kvMap
}
