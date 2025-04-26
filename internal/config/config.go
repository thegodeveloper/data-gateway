// Package config
// internal/config/config.go
package config

import "os"

type Config struct {
	PostgresConnStr string
	MongoURI        string
	HTTPPort        string
}

func Load() *Config {
	return &Config{
		PostgresConnStr: os.Getenv("POSTGRES_CONN_STR"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
