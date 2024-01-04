package db

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	connStr := "user=bank-transfer dbname=postgres password=123456 port=5432 host=localhost search_path=bank sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Queries struct {
	DB DBTX
}

func (q *Queries) ExecuteTx(ctx context.Context, f func(*Queries) error) error {
	var db *sql.DB
	switch t := q.DB.(type) {
	case *sql.Tx:
		return errors.New("Nested transaction")
	case *sql.DB:
		db = t
	}
	if db == nil {
		return errors.New("*sql.DB is null")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := &Queries{
		DB: tx,
	}
	err = f(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
