package databasex

import (
	"context"
	"database/sql"
	"time"

	"github.com/choyri/gugu/internal/logx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func NewMySQL(ctx context.Context, conf Config) (*bun.DB, error) {
	logger := logx.FromContext(ctx)
	logger.Debug("prepare open mysql db")

	stdDB, err := sql.Open("mysql", conf.DSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = stdDB.PingContext(ctx); err != nil {
		return nil, err
	}

	logger.Debug("db ping succeeded")

	stdDB.SetMaxOpenConns(int(conf.MaxOpenConns))
	stdDB.SetMaxIdleConns(int(conf.MaxIdleConns))

	db := bun.NewDB(stdDB, mysqldialect.New())

	if conf.Debug {
		db.AddQueryHook(newBunLogHook(logger))
		logger.Debug("bun query log hook added")
	}

	return db, nil
}
