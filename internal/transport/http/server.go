// Package http
// internal/transport/http/server.go
package http

import (
	"encoding/json"
	"github.com/thegodeveloper/data-gateway/internal/app"
	"github.com/thegodeveloper/data-gateway/internal/domain"
	"net/http"
)

func StartServer(service *app.GatewayService, port string) {
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req domain.QueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.Query(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.ListenAndServe(":"+port, nil)
}
