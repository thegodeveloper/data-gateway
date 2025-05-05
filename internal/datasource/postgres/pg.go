// Package postgres
// internal/datasource/postgres/pg.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/thegodeveloper/data-gateway/internal/domain"
)

type PostgresSource struct {
	db *sql.DB
}

func NewPostgresSource(db *sql.DB) *PostgresSource {
	return &PostgresSource{db: db}
}

func (p *PostgresSource) Query(ctx context.Context, req domain.QueryRequest) (any, error) {
	queryStr, ok := req.Params["query"].(string)
	if !ok {
		return nil, errors.New("missing or invalid 'query' parameter")
	}

	rows, err := p.db.QueryContext(ctx, queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}
		results = append(results, rowMap)
	}

	return results, nil
}
