package model

import (
	"context"
	"database/sql"
	"regexp"
)

var ReGroupBy = regexp.MustCompile(`(?i)group\s+by`)

// DBConnect interface for db connection
type DBConnect interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// GetDBToDBConnect create DBConnect from getter
func GetDBToDBConnect(getter func() *sql.DB) func() DBConnect {
	return func() DBConnect {
		return getter()
	}
}

// TxToDBConnect create DBConnect from tx
func TxToDBConnect(tx *sql.Tx) func() DBConnect {
	return func() DBConnect {
		return tx
	}
}
