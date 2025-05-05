// Package main
// cmd/gateway/main.go
package main

import (
	"context"
	"github.com/thegodeveloper/data-gateway/internal/app"
	"github.com/thegodeveloper/data-gateway/internal/config"
	"github.com/thegodeveloper/data-gateway/internal/datasource/dynamodb"
	"github.com/thegodeveloper/data-gateway/internal/datasource/mongodb"
	"github.com/thegodeveloper/data-gateway/internal/datasource/postgres"
	"github.com/thegodeveloper/data-gateway/internal/domain"
	"github.com/thegodeveloper/data-gateway/internal/transport/http"
	"github.com/thegodeveloper/data-gateway/pkg/common"
	"github.com/thegodeveloper/data-gateway/pkg/otel"
)

func main() {
	cfg := config.Load()

	shutdown, err := otel.InitTracer("data-gateway")
	if err != nil {
		common.Error("failed to init OpenTelemetry: %v", err)
		return
	}
	defer shutdown(context.Background())

	pg, err := postgres.NewPostgresSource(cfg.PostgresConnStr)
	if err != nil {
		common.Error("Postgres init failed: %v", err)
		return
	}

	dynamo, err := dynamodb.NewDynamoDBSource()
	if err != nil {
		common.Error("DynamoDB init failed: %v", err)
		return
	}

	mongo, err := mongodb.NewMongoSource(cfg.MongoURI)
	if err != nil {
		common.Error("MongoDB init failed: %v", err)
		return
	}

	svc := app.NewGatewayService(map[string]domain.DataSource{
		"postgres": pg,
		"dynamodb": dynamo,
		"mongodb":  mongo,
	})

	common.Info("Starting HTTP server on port %s", cfg.HTTPPort)
	http.StartServer(svc, cfg.HTTPPort)
}
