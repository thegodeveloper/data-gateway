package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.28.0"

	"github.com/spf13/viper"
	"github.com/thegodeveloper/data-gateway/internal/adapters/repositories/dynamodb"
	"github.com/thegodeveloper/data-gateway/internal/adapters/repositories/postgres"
	"github.com/thegodeveloper/data-gateway/internal/core/ports"
	"github.com/thegodeveloper/data-gateway/internal/core/services"
)

type Config struct {
	DataSources map[string]interface{} `mapstructure:"data-sources"`
	Paths       map[string]interface{} `mapstructure:"paths"`
}

// NewDataService initializes the appropriate DataService based on the request path
// and the configuration file.
func NewDataService(requestPath string, cfg *Config) (ports.DataService, error) {
	var dataSourceName string

	for pathPattern, sourceConfig := range cfg.Paths {
		if strings.HasPrefix(requestPath, pathPattern) {
			switch v := sourceConfig.(type) {
			case string:
				dataSourceName = v
			case map[interface{}]interface{}:
				if source, ok := v["source"].(string); ok {
					dataSourceName = source
				} else {
					return nil, fmt.Errorf("invalid source configuration for path '%s'", pathPattern)
				}
			default:
				return nil, fmt.Errorf("invalid source configuration for path '%s'", pathPattern)
			}
			break // Found a matching path
		}
	}

	if dataSourceName == "" {
		return nil, fmt.Errorf("no data source configured for path '%s'", requestPath)
	}

	switch dataSourceName {
	case "postgres":
		connStr, ok := cfg.DataSources["postgres"].(string)
		if !ok {
			return nil, fmt.Errorf("postgres connection string not found in configuration: '%s'", dataSourceName)
		}
		db, err := postgres.NewPostgresDB(connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL: '%s': %w", connStr, err)
		}
		return services.NewDataService(postgres.NewPostgresRepository(db)), nil
	case "dynamodb":
		region, ok := cfg.DataSources["dynamodb"].(string)
		if !ok {
			return nil, fmt.Errorf("dynamodb region not found in configuration")
		}
		sess, err := dynamodb.NewDynamoDBSession(region)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize DynamoDB session: %w", err)
		}
		return services.NewDataService(dynamodb.NewDynamoDBRepository(sess, "your-dynamodb-table")), nil // Replace with your table name (can also be in config)
	case "mongodb":
		// Implement MongoDB initialization here based on the config
		uriConfig, ok := cfg.Paths["/invoices"].(map[interface{}]interface{})
		if !ok || uriConfig["source"] != "mongodb" {
			return nil, fmt.Errorf("mongodb configuration not found for path '/invoices'")
		}
		uri, ok := uriConfig["uri"].(string)
		if !ok {
			return nil, fmt.Errorf("mongodb URI not found in configuration")
		}
		// Initialize MongoDB repository (you'll need to create this)
		// return services.NewDataService(mongodb.NewMongoDBRepository(uri)), nil
		return nil, fmt.Errorf("MongoDB support not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported data source: %s", dataSourceName)
	}
}

// LoadConfig loads the configuration from the specified YAML file.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

// InitTracer initializes the OpenTelemetry tracer provider.
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	// ... (same as before) ...
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion("v1.0.0"),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)
	return tp, nil
}
