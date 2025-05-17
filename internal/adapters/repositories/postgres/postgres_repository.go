// internal/adapters/repositories/postgres/postgres_repository.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/thegodeveloper/data-gateway/internal/core/ports"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return db, nil
}

func (r *PostgresRepository) FetchData(ctx context.Context, path string, query map[string]interface{}) (interface{}, error) {
	ctx, span := otel.Tracer("data-layer").Start(ctx, "postgresRepository.FetchData")
	defer span.End()

	// Example: Assuming the 'path' corresponds to a table name
	tableName := path

	// Build the query based on the 'query' map. This is a simplified example.
	// In a real application, you'd need more robust query building logic
	// to prevent SQL injection and handle various query parameters.
	whereClause := ""
	values := []interface{}{}
	if len(query) > 0 {
		whereClause = "WHERE "
		i := 0
		for key, value := range query {
			if i > 0 {
				whereClause += " AND "
			}
			whereClause += fmt.Sprintf("%s = $%d", key, i+1)
			values = append(values, value)
			i++
		}
	}

	sqlStatement := fmt.Sprintf("SELECT * FROM %s %s", tableName, whereClause)

	rows, err := r.db.QueryContext(ctx, sqlStatement, values...)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to query PostgreSQL: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get column names: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Handle different data types if needed
			row[col] = val
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return results, nil
}

// You might need other methods like Create, Update, Delete based on your requirements.
