package logx

import (
	"context"

	"golang.org/x/exp/slog"
)

type ctxKey struct{}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok && l == logger {
		return ctx
	}

	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}
