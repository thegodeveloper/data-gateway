// Package dynamodb
// internal/datasource/dynamodb/dynamodb.go
package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thegodeveloper/data-gateway/internal/domain"
)

type DynamoDBSource struct {
	client *dynamodb.Client
}

func NewDynamoDBSource() (*DynamoDBSource, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBSource{client: client}, nil
}

func (d *DynamoDBSource) Query(req domain.QueryRequest) (any, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(req.Params["table"].(string)),
	}
	result, err := d.client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
