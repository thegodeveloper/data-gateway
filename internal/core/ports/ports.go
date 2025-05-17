// internal/core/ports/ports.go
package ports

import (
	"context"
)

// DataRequest represents the generic request format from microservices.
type DataRequest struct {
	Path    string                 `json:"path"`
	Payload map[string]interface{} `json:"payload"`
}

// DataResponse represents the generic response format to microservices.
type DataResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
}

// DataPort defines the interface for the data repository.
type DataPort interface {
	FetchData(ctx context.Context, path string, query map[string]interface{}) (interface{}, error)
}

// DataService defines the interface for the data service.
type DataService interface {
	ProcessRequest(ctx context.Context, req DataRequest) (DataResponse, error)
}
