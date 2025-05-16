package app

import "context"

type DataSource interface {
	Query(ctx context.Context, query map[string]interface{}) (interface{}, error)
}
