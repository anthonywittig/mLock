package property

import (
	"context"
	"fmt"
	"log"
	"mlock/shared"
	"mlock/shared/dynamo"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const (
	tableName = "Property"
	itemType  = "1" // Not really using ATM.
)

func Delete(ctx context.Context, name string) error {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	name = strings.TrimSpace(name)

	// TODO: don't delete if in use?

	// No audit trail for deletes. :(

	if _, err = dy.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"type": {S: aws.String(itemType)},
			"name": {S: aws.String(name)},
		},
		TableName: aws.String(tableName),
	}); err != nil {
		return fmt.Errorf("error deleting item: %s", err.Error())
	}

	return nil
}

func Get(ctx context.Context, name string) (shared.Property, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Property{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	name = strings.TrimSpace(name)

	result, err := dy.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"type": {S: aws.String(itemType)},
			"name": {S: aws.String(name)},
		},
	})
	if err != nil {
		return shared.Property{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.Property{}, false, nil
	}

	item := shared.Property{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return shared.Property{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return item, true, nil
}

func List(ctx context.Context) ([]shared.Property, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.Property{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"type": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(itemType)},
				},
			},
		},
	}

	items := []shared.Property{}
	for {
		// Get the list of tables
		result, err := dy.QueryWithContext(ctx, input)
		if err != nil {
			return []shared.Property{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := shared.Property{}
			if err = dynamodbattribute.UnmarshalMap(i, &item); err != nil {
				return []shared.Property{}, fmt.Errorf("error unmarshaling: %s", err.Error())
			}
			items = append(items, item)
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	return items, nil
}

func Put(ctx context.Context, oldKey string, name string) (shared.Property, error) {
	return PutID(ctx, oldKey, name, uuid.New().String())
}

func PutID(ctx context.Context, oldKey string, name string, id string) (shared.Property, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Property{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	name = strings.TrimSpace(name)

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.Property{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.Property{}, fmt.Errorf("no current user")
	}

	item := shared.Property{
		Type:      itemType,
		Name:      name,
		ID:        id,
		CreatedBy: currentUser.Email,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return shared.Property{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItemWithContext(ctx, input)
	if err != nil {
		return shared.Property{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := Get(ctx, name)
	if err != nil {
		return shared.Property{}, err
	}
	if !ok {
		return shared.Property{}, fmt.Errorf("couldn't find entity after insert")
	}

	if oldKey != "" && oldKey != entity.Name {
		if err := Delete(ctx, oldKey); err != nil {
			return shared.Property{}, fmt.Errorf("error deleting old item: %s", err.Error())
		}
	}

	return entity, nil
}

func Migrate(ctx context.Context) error {
	exists, err := dynamo.TableExists(ctx, tableName)
	if err != nil {
		return fmt.Errorf("error checking for table: %s", err.Error())
	}
	if exists {
		return nil
	}

	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("type"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("type"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("RANGE"),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dy.CreateTableWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	log.Printf("created table: %s - %+v", tableName, result)

	return nil
}
