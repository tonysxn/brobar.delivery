package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	defaultQueryTimeout = 5 * time.Second
)

type dbExecutor interface {
	sqlx.ExtContext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
