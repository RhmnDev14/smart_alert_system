package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(databaseURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{DB: db}, nil
}

type txKey struct{}

func (p *PostgresDB) Ext(ctx context.Context) interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
} {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return p.DB
}

func (p *PostgresDB) RunInTx(ctx context.Context, fn func(txCtx context.Context) error) (err error) {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // err is not nil; rollback
		} else {
			err = tx.Commit() // err is nil; if Commit fails, return that error
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Execute the transaction block
	err = fn(txCtx)

	return err
}

func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
