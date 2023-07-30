package databasex

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"time"

	"github.com/choyri/gugu/internal/logx"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slog"
)

type logHook struct {
	logger *slog.Logger
}

var _ bun.QueryHook = (*logHook)(nil)

func newBunLogHook(logger *slog.Logger) *logHook {
	return &logHook{
		logger: logger,
	}
}

func (h *logHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (h *logHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	var (
		fields []slog.Attr
		lvl    = slog.LevelDebug
	)

	if event.Err != nil &&
		!errors.Is(event.Err, sql.ErrNoRows) &&
		!errors.Is(event.Err, sql.ErrTxDone) {
		fields = append(
			fields,
			logx.Error(event.Err),
			slog.String("error_type", reflect.TypeOf(event.Err).String()),
		)
		lvl = slog.LevelError
	}

	slog.LogAttrs(ctx, lvl, "bun query log hook", append(fields,
		slog.String("db.query", event.Query),
		slog.String("db.operation", event.Operation()),
		slog.Duration("db.duration", time.Since(event.StartTime).Round(time.Microsecond)),
	)...)
}
