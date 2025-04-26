// Package domain
// domain/datasource.go
package domain

type QueryRequest struct {
	Source string
	Query  string
	Params map[string]any
}

type DataSource interface {
	Query(req QueryRequest) (any, error)
}
