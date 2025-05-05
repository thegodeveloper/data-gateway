// Package mongodb
// internal/datasource/mongodb/mongo.go
package mongodb

import (
	"context"
	"errors"
	"github.com/thegodeveloper/data-gateway/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSource struct {
	client *mongo.Client
}

func NewMongoSource(client *mongo.Client) *MongoSource {
	return &MongoSource{client: client}
}

func (m *MongoSource) Query(ctx context.Context, req domain.QueryRequest) (any, error) {
	dbName, ok := req.Params["database"].(string)
	if !ok {
		return nil, errors.New("missing 'database' parameter")
	}
	collectionName, ok := req.Params["collection"].(string)
	if !ok {
		return nil, errors.New("missing 'collection' parameter")
	}

	filterRaw, ok := req.Params["filter"].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing or invalid 'filter' parameter")
	}
	filter := bson.M(filterRaw)

	coll := m.client.Database(dbName).Collection(collectionName)
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
