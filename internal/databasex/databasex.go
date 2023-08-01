package databasex

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
)

func New(ctx context.Context, conf Config) (*bun.DB, error) {
	if strings.HasPrefix(conf.DSN, "file:") {
		return NewSQLite(ctx, conf)
	}

	return NewMySQL(ctx, conf)
}
