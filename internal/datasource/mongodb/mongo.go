// Package mongodb
// internal/datasource/mongodb/mongo.go
package mongodb

import (
	"context"
	"github.com/thegodeveloper/data-gateway/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoSource struct {
	client *mongo.Client
}

func NewMongoSource(uri string) (*MongoSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoSource{client: client}, nil
}

func (m *MongoSource) Query(req domain.QueryRequest) (any, error) {
	db := m.client.Database(req.Params["database"].(string))
	col := db.Collection(req.Params["collection"].(string))

	filter := bson.M{}
	if f, ok := req.Params["filter"].(map[string]any); ok {
		filter = f
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
