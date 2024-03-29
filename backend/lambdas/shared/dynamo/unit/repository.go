package unit

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type Repository struct{}

const (
	tableName = "Unit_v2"
)

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	// TODO: under what circumstances would we want to stop this?

	// No audit trail for deletes. :(

	if _, err = dy.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberB{Value: id[:]},
		},
		TableName: aws.String(tableName),
	}); err != nil {
		return fmt.Errorf("error deleting item: %s", err.Error())
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (shared.Unit, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	result, err := dy.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberB{Value: id[:]},
		},
	})
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.Unit{}, false, nil
	}

	item := shared.Unit{}
	err = dynamo.UnmarshalMapWithOptions(result.Item, &item)
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return item, true, nil
}

func (r *Repository) List(ctx context.Context) ([]shared.Unit, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.Unit{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	items := []shared.Unit{}
	for {
		result, err := dy.Scan(ctx, input)
		if err != nil {
			return []shared.Unit{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := shared.Unit{}
			if err = dynamo.UnmarshalMapWithOptions(i, &item); err != nil {
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

func (r *Repository) ListByName(ctx context.Context) (map[string]shared.Unit, error) {
	units, err := r.List(ctx)
	if err != nil {
		return map[string]shared.Unit{}, fmt.Errorf("error getting units: %s", err.Error())
	}

	byName := map[string]shared.Unit{}
	for _, u := range units {
		byName[u.Name] = u
	}

	return byName, nil
}

func (r *Repository) ListByID(ctx context.Context) (map[uuid.UUID]shared.Unit, error) {
	units, err := r.List(ctx)
	if err != nil {
		return map[uuid.UUID]shared.Unit{}, fmt.Errorf("error getting units: %s", err.Error())
	}

	byID := map[uuid.UUID]shared.Unit{}
	for _, u := range units {
		byID[u.ID] = u
	}

	return byID, nil
}

func (r *Repository) Put(ctx context.Context, item shared.Unit) (shared.Unit, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	if item.ID == uuid.Nil {
		// Since an ID can easily be forgotten, let's never assume we need to create one.
		return shared.Unit{}, fmt.Errorf("an ID is required")
	}

	item.Name = strings.TrimSpace(item.Name)
	item.RemotePropertyURL = strings.TrimSpace(item.RemotePropertyURL)

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.Unit{}, fmt.Errorf("no current user")
	}
	item.UpdatedBy = currentUser.Email

	av, err := dynamo.MarshalMapWithOptions(item)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItem(ctx, input)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := r.Get(ctx, item.ID)
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
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: "B",
			},
		},
		BillingMode: "PAY_PER_REQUEST",
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       "HASH",
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dy.CreateTable(ctx, input)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	log.Printf("created table: %s - %+v", tableName, result)

	return nil
}

func migrateData(ctx context.Context) error {
	return nil
}
