// Package dynamodb
// internal/datasource/dynamodb/dynamodb.go
package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	sdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/thegodeveloper/data-gateway/internal/domain"
)

type Source struct {
	client *sdynamodb.Client
}

func NewSource(client *sdynamodb.Client) *Source {
	return &Source{client: client}
}

func (s *Source) Query(ctx context.Context, req domain.QueryRequest) (any, error) {

	tableName, ok := req.Params["table"].(string)
	if !ok || tableName == "" {
		return nil, errors.New("missing or invalid 'table' parameter")
	}

	keyMap, ok := req.Params["key"].(map[string]interface{})
	if !ok || len(keyMap) == 0 {
		return nil, errors.New("missing or invalid 'key' parameter")
	}

	keyCondition := ""
	exprAttrValues := make(map[string]types.AttributeValue)
	index := 0

	for key, value := range keyMap {
		placeholder := fmt.Sprintf(":v%d", index)
		keyCondition += fmt.Sprintf("%s = %s", key, placeholder)
		if index < len(keyMap)-1 {
			keyCondition += " AND "
		}
		av, err := attributevalue.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal key value: %w", err)
		}
		exprAttrValues[placeholder] = av
		index++
	}

	out, err := s.client.Query(ctx, &sdynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: exprAttrValues,
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb query failed: %w", err)
	}

	var results []map[string]interface{}
	for _, item := range out.Items {
		var record map[string]interface{}
		if err := attributevalue.UnmarshalMap(item, &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal result: %w", err)
		}
		results = append(results, record)
	}

	return results, nil
}
