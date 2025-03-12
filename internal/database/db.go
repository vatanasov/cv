package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNoRowsAffected = errors.New("no rows affected")
)

type Db interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	return db, nil
}

func Insert(db Db, query string, args []any) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected != 1 {
		return 0, ErrNoRowsAffected
	}

	return result.LastInsertId()
}

func GetAll[T any](db Db, query string, args []any, resultsF func(*T) []any) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var resultList []T

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return resultList, err
	}
	defer rows.Close()

	for rows.Next() {
		var result T

		err := rows.Scan(resultsF(&result)...)
		if err != nil {
			return resultList, err
		}
		resultList = append(resultList, result)
	}

	return resultList, nil
}

func Get(db Db, query string, args []any, scanInto []any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, query, args...)

	return row.Scan(scanInto...)
}
