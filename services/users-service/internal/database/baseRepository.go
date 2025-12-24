package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type BaseSqlRepository[T any] struct {
	DB *sql.DB
}

func (repo *BaseSqlRepository[T]) SelectMultiple(mapRow func(*sql.Rows, *T) error, query string, args ...any) ([]*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var list []*T

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var t T
		if err := mapRow(rows, &t); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (repo *BaseSqlRepository[T]) SelectSingle(mapRow func(*sql.Row, *T) error, query string, args ...any) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := repo.DB.QueryRowContext(ctx, query, args...)
	var t T
	if err := mapRow(row, &t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &t, nil
}

func (repo *BaseSqlRepository[T]) SelectSingleWithContext(ctx context.Context, mapRow func(*sql.Row, *T) error, query string, args ...any) (*T, error) {
	row := repo.DB.QueryRowContext(ctx, query, args...)
	var t T
	if err := mapRow(row, &t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &t, nil
}

func (repo *BaseSqlRepository[T]) SelectMultipleWithContext(ctx context.Context, mapRow func(*sql.Rows, *T) error, query string, args ...any) ([]*T, error) {
	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var list []*T

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var t T
		if err := mapRow(rows, &t); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (repo *BaseSqlRepository[T]) Insert(ctx context.Context, query string, args ...any) (string, error) {
	var id string
	query = query + " RETURNING id"

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (repo *BaseSqlRepository[T]) ExecuteQuery(query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *BaseSqlRepository[T]) GetBoolValue(query string, args ...any) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var boolValue bool
	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&boolValue) // Use QueryRow and Scan
	if err != nil {
		return false, err
	}
	return boolValue, nil
}
