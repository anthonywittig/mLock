package unit

import (
	"context"
	"fmt"
	"log"
	"mlock/shared"
	"mlock/shared/dynamo"
	"mlock/shared/dynamo/property"
	"mlock/shared/dynamo/unit/last"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const (
	tableName = "Unit_v2"
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

func Get(ctx context.Context, id uuid.UUID) (shared.Unit, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	result, err := dy.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {B: id[:]},
		},
	})
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.Unit{}, false, nil
	}

	item := shared.Unit{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return item, true, nil
}

func List(ctx context.Context) ([]shared.Unit, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.Unit{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	items := []shared.Unit{}
	for {
		result, err := dy.ScanWithContext(ctx, input)
		if err != nil {
			return []shared.Unit{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := shared.Unit{}
			if err = dynamodbattribute.UnmarshalMap(i, &item); err != nil {
				return []shared.Unit{}, fmt.Errorf("error unmarshaling: %s", err.Error())
			}
			items = append(items, item)
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}

func Put(ctx context.Context, item shared.Unit) (shared.Unit, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	if item.ID == uuid.Nil {
		// Since an ID can easily be forgotten, let's never assume we need to create one.
		return shared.Unit{}, fmt.Errorf("an ID is required")
	}

	item.Name = strings.TrimSpace(item.Name)
	item.CalendarURL = strings.TrimSpace(item.CalendarURL)

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.Unit{}, fmt.Errorf("no current user")
	}
	item.UpdatedBy = currentUser.Email

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItemWithContext(ctx, input)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := Get(ctx, item.ID)
	if err != nil {
		return shared.Unit{}, err
	}
	if !ok {
		return shared.Unit{}, fmt.Errorf("couldn't find entity after insert")
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
	existingItems, err := List(ctx)
	if err != nil {
		return fmt.Errorf("error getting existing items: %s", err.Error())
	}
	if len(existingItems) > 0 {
		log.Println("already migrated unit data")
		return nil
	}

	items, err := last.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting items: %s", err.Error())
	}

	properties, err := property.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting properties: %s", err.Error())
	}

	for _, item := range items {
		cd, err := shared.GetContextData(ctx)
		if err != nil {
			return fmt.Errorf("error getting context data: %s", err.Error())
		}

		propertyID := uuid.Nil
		for _, p := range properties {
			if p.Name == item.PropertyName {
				propertyID = p.ID
			}
		}
		if propertyID == uuid.Nil {
			return fmt.Errorf("couldn't find property: %s", item.PropertyName)
		}

		cd.User = &shared.User{Email: item.UpdatedBy}
		if _, err := Put(ctx, shared.Unit{
			ID:          uuid.New(),
			Name:        item.Name,
			PropertyID:  propertyID,
			CalendarURL: item.CalendarURL,
			UpdatedBy:   item.UpdatedBy,
		}); err != nil {
			return fmt.Errorf("error getting context data: %s", err.Error())
		}
	}

	return nil
}
