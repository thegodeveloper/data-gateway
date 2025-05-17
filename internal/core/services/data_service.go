// internal/core/services/data_service.go
package services

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"

	"github.com/thegodeveloper/data-gateway/internal/core/ports"
)

type dataService struct {
	repo ports.DataPort
}

func NewDataService(repo ports.DataPort) ports.DataService {
	return &dataService{repo: repo}
}

func (s *dataService) ProcessRequest(ctx context.Context, req ports.DataRequest) (ports.DataResponse, error) {
	ctx, span := otel.Tracer("data-layer").Start(ctx, "dataService.ProcessRequest")
	defer span.End()

	// Here you would implement the logic to transform the generic request
	// based on the 'path' and 'payload' into a format suitable for the
	// underlying data source. This is the key abstraction point.

	// For this example, we'll just pass the payload as the query.
	// In a real-world scenario, you'd have more sophisticated logic here.
	query := req.Payload

	data, err := s.repo.FetchData(ctx, req.Path, query)
	if err != nil {
		span.RecordError(err)
		return ports.DataResponse{Error: fmt.Sprintf("failed to fetch data for path '%s': %v", req.Path, err)}, nil
	}

	return ports.DataResponse{Data: data}, nil
}

// Helper function to convert map[string]interface{} to JSON string (if needed by the repository)
func mapToJSONString(m map[string]interface{}) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
