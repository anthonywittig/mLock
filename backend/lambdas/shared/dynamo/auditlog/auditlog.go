package auditlog

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const (
	tableName = "AuditLog_v1"
)

func Delete(ctx context.Context, id uuid.UUID) error {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	// TODO: under what circumstances would we want to stop this?

	// No audit trail for deletes. :(

	if _, err = dy.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {B: id[:]},
		},
		TableName: aws.String(tableName),
	}); err != nil {
		return fmt.Errorf("error deleting item: %s", err.Error())
	}

	return nil
}

func Get(ctx context.Context, id uuid.UUID) (shared.AuditLog, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.AuditLog{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	result, err := dy.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {B: id[:]},
		},
	})
	if err != nil {
		return shared.AuditLog{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.AuditLog{}, false, nil
	}

	item := &shared.AuditLog{}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return shared.AuditLog{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return *item, true, nil
}

func List(ctx context.Context) ([]shared.AuditLog, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.AuditLog{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	items := []shared.AuditLog{}
	for {
		result, err := dy.ScanWithContext(ctx, input)
		if err != nil {
			return []shared.AuditLog{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := shared.AuditLog{}
			if err = dynamodbattribute.UnmarshalMap(i, &item); err != nil {
				return []shared.AuditLog{}, fmt.Errorf("error unmarshaling: %s", err.Error())
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

func Put(ctx context.Context, item shared.AuditLog) (shared.AuditLog, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.AuditLog{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	if item.ID == uuid.Nil {
		// Since an ID can easily be forgotten, let's never assume we need to create one.
		return shared.AuditLog{}, fmt.Errorf("an ID is required")
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return shared.AuditLog{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItemWithContext(ctx, input)
	if err != nil {
		return shared.AuditLog{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := Get(ctx, item.ID)
	if err != nil {
		return shared.AuditLog{}, err
	}
	if !ok {
		return shared.AuditLog{}, fmt.Errorf("couldn't find entity after insert")
	}

	return entity, nil
}

func Migrate(ctx context.Context) error {
	if err := migrateCreateTable(ctx); err != nil {
		return fmt.Errorf("error creating table: %s", err.Error())
	}

	if err := migrateData(ctx); err != nil {
		return fmt.Errorf("error migrating data: %s", err.Error())
	}

	return nil
}

func migrateCreateTable(ctx context.Context) error {
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
				AttributeName: aws.String("id"),
				AttributeType: aws.String("B"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
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

func migrateData(ctx context.Context) error {
	return nil
}
