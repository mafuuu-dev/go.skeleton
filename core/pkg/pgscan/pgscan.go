package pgscan

import (
	"backend/core/pkg/errorsx"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Scannable interface {
	Scan(dest ...any) error
}

func Scan[T any](source any, mapRow func(Scannable) (T, error)) ([]T, error) {
	switch s := source.(type) {

	case pgx.Rows:
		defer s.Close()
		var result []T
		for s.Next() {
			item, err := mapRow(s)
			if err != nil {
				return nil, errorsx.Error(err)
			}

			result = append(result, item)
		}

		return result, nil

	case pgx.Row:
		item, err := mapRow(s)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}

			return nil, errorsx.Error(err)
		}

		return []T{item}, nil

	default:
		return nil, errorsx.Errorf("Unsupported input type for Scan")
	}
}

func ScanOne[T any](row pgx.Row, mapRow func(Scannable) (T, error)) (*T, error) {
	results, err := Scan(row, mapRow)
	if err != nil || len(results) == 0 {
		return nil, errorsx.Error(err)
	}

	return &results[0], nil
}
