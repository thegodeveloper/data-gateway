// internal/adapters/repositories/dynamodb/dynamodb_repository.go
package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.opentelemetry.io/otel"

	"github.com/thegodeveloper/data-gateway/internal/core/ports"
)

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(sess *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{client: sess, tableName: tableName}
}

func NewDynamoDBSession(region string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}

func (r *DynamoDBRepository) FetchData(ctx context.Context, path string, query map[string]interface{}) (interface{}, error) {
	ctx, span := otel.Tracer("data-layer").Start(ctx, "dynamoDBRepository.FetchData")
	defer span.End()

	// Example: Assuming the 'path' might influence how we query DynamoDB
	// and the 'query' map contains the key-value pairs for filtering.

	// In a real-world scenario, you would likely need to build a more
	// specific DynamoDB query based on the 'path' and the structure of
	// your DynamoDB table.

	// This is a very basic example assuming direct attribute matching.
	filter := map[string]types.AttributeValue{}
	for key, value := range query {
		av, err := attributevalue.Marshal(value)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to marshal attribute '%s': %w", key, err)
		}
		filter[key] = av
	}

	input := &dynamodb.ScanInput{
		TableName:              aws.String(r.tableName),
		ScanFilter:             filter,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	output, err := r.client.Scan(ctx, input)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to scan DynamoDB table '%s': %w", r.tableName, err)
	}

	var results []map[string]interface{}
	if err := attributevalue.UnmarshalListOfMaps(output.Items, &results); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to unmarshal DynamoDB items: %w", err)
	}

	return results, nil
}

// You might need other methods like PutItem, UpdateItem, DeleteItem based on your requirements.
