// Package postgres
// internal/datasource/postgres/pg.go
package postgres

import (
	"database/sql"
	"github.com/thegodeveloper/data-gateway/internal/domain"
)

type PostgresSource struct {
	db *sql.DB
}

func NewPostgresSource(connStr string) (*PostgresSource, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresSource{db: db}, nil
}

func (p *PostgresSource) Query(req domain.QueryRequest) (any, error) {
	rows, err := p.db.Query(req.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	result := []map[string]interface{}{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			rowMap[colName] = *val
		}
		result = append(result, rowMap)
	}

	return result, nil
}
