package app

import (
	"database/sql"
	"github.com/thegodeveloper/data-gateway/internal/datasource/dynamodb"
	"github.com/thegodeveloper/data-gateway/internal/datasource/mongodb"
	"github.com/thegodeveloper/data-gateway/internal/datasource/postgres"
	"github.com/thegodeveloper/data-gateway/internal/domain"
	"github.com/thegodeveloper/data-gateway/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	router       *gin.Engine
	dataSources  map[string]domain.DataSource
	pgDB         *sql.DB
	mongoClient  *mongo.Client
	dynamoClient *dynamodb.Client
}

// NewServer creates a new server instance and sets up all dependencies.
func NewServer(pgDB *sql.DB, mongoClient *mongo.Client, dynamoClient *dynamodb.Client) *Server {
	s := &Server{
		router:       gin.Default(),
		dataSources:  make(map[string]domain.DataSource),
		pgDB:         pgDB,
		mongoClient:  mongoClient,
		dynamoClient: dynamoClient,
	}

	s.registerMiddlewares()
	s.registerRoutes()
	s.registerDataSources()

	return s
}

// Start runs the HTTP server on the given address.
func (s *Server) Start(addr string) {
	log.Printf("Starting server at %s", addr)
	if err := s.router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// registerMiddlewares adds global middlewares to the router.
func (s *Server) registerMiddlewares() {
	// e.g., s.router.Use(cors.Default())
	// Instrument with OpenTelemetry (if added): s.router.Use(otelgin.Middleware("data-gateway"))
}

// registerRoutes sets up API routes and handlers.
func (s *Server) registerRoutes() {
	h := handler.NewQueryHandler(s.dataSources)
	api := s.router.Group("/api/v1")
	{
		api.POST("/query/:source", h.HandleQuery)
	}
}

// registerDataSources wires up each supported data source.
func (s *Server) registerDataSources() {
	// PostgreSQL
	postgresSource := postgres.NewSource(s.pgDB)
	s.dataSources["postgres"] = postgresSource

	// MongoDB
	mongoSource := mongodb.NewSource(s.mongoClient)
	s.dataSources["mongodb"] = mongoSource

	// DynamoDB
	dynamoSource := dynamodb.NewSource(s.dynamoClient)
	s.dataSources["dynamodb"] = dynamoSource
}
