// Package app
// internal/app/gateway_service.go
package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/thegodeveloper/data-gateway/internal/domain"
)

type GatewayService struct {
	dataSources map[string]domain.DataSource
}

func NewGatewayService(dataSources map[string]domain.DataSource) *GatewayService {
	return &GatewayService{
		dataSources: dataSources,
	}
}

// HandleQuery processes the request and routes it to the correct data source.
func (s *GatewayService) HandleQuery(ctx context.Context, req domain.QueryRequest) (any, error) {
	if req.Source == "" {
		return nil, errors.New("missing 'source' field in request")
	}

	ds, ok := s.dataSources[req.Source]
	if !ok {
		return nil, fmt.Errorf("data source '%s' not supported", req.Source)
	}

	result, err := ds.Query(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("query failed for '%s': %w", req.Source, err)
	}

	return result, nil
}
