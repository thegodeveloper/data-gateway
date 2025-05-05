// Package domain
// domain/datasource.go
package domain

import "context"

type QueryRequest struct {
	Source string                 `json:"source"`
	Params map[string]interface{} `json:"params"`
}

type DataSource interface {
	Query(ctx context.Context, req QueryRequest) (any, error)
}
