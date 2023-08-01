package databasex

import (
	"context"
	"database/sql"
	"time"

	"github.com/choyri/gugu/internal/logx"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"golang.org/x/exp/slog"
)

func NewSQLite(ctx context.Context, conf Config) (*bun.DB, error) {
	logger := logx.FromContext(ctx)
	logger.Debug("prepare open sqlite db", slog.String("dsn", conf.DSN))

	stdDB, err := sql.Open(sqliteshim.ShimName, conf.DSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = stdDB.PingContext(ctx); err != nil {
		return nil, err
	}

	logger.Debug("db ping succeeded")

	stdDB.SetMaxOpenConns(1)

	db := bun.NewDB(stdDB, sqlitedialect.New())

	if conf.Debug {
		db.AddQueryHook(newBunLogHook(logger))
		logger.Debug("bun query log hook added")
	}

	return db, nil
}
