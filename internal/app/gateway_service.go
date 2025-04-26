// Package app
// internal/app/gateway_service.go
package app

import (
	"fmt"
	"github.com/thegodeveloper/data-gateway/internal/domain"
)

type GatewayService struct {
	sources map[string]domain.DataSource
}

func NewGatewayService(sources map[string]domain.DataSource) *GatewayService {
	return &GatewayService{sources: sources}
}

func (g *GatewayService) Query(req domain.QueryRequest) (any, error) {
	source, ok := g.sources[req.Source]
	if !ok {
		return nil, fmt.Errorf("data source %s not found", req.Source)
	}
	return source.Query(req)
}
