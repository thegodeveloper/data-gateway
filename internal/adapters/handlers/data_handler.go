// internal/adapters/handlers/data_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/gorilla/mux"
	"github.com/thegodeveloper/data-gateway/internal/core/ports"
)

type DataHandler struct {
	service ports.DataService
}

func NewDataHandler(service ports.DataService) *DataHandler {
	return &DataHandler{service: service}
}

func (h *DataHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("data-layer").Start(ctx, "dataHandler.HandleRequest")
	defer span.End()

	vars := mux.Vars(r)
	path := vars["path"]

	var req ports.DataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		span.RecordError(err)
		return
	}
	req.Path = path // Ensure the path from the URL is included in the request

	response, err := h.service.ProcessRequest(ctx, req)
	if err != nil {
		// Log the error internally, the response error is already set
		span.RecordError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		span.RecordError(err)
		return
	}
}
