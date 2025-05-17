package main

import (
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/gorilla/mux"
	"github.com/thegodeveloper/data-gateway/internal/adapters/handlers"
	"github.com/thegodeveloper/data-gateway/internal/app"
)

const serviceName = "data-layer"

func main() {
	// Initialize OpenTelemetry
	tp, err := app.InitTracer(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(nil); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Load the configuration file
	cfg, err := app.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the handler
	dataHandler := handlers.NewDataHandler(func(path string) (ports.DataService, error) {
		return app.NewDataService(path, cfg)
	})

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/{path}", dataHandler.HandleRequest).Methods(http.MethodPost)

	// Wrap the router with OpenTelemetry instrumentation
	handler := otelhttp.NewHandler(router, serviceName)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))

}
