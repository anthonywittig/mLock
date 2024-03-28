package climatecontrol

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type Repository struct{}

const (
	tableName = "ClimateControl_v1"
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

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (shared.ClimateControl, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.ClimateControl{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	result, err := dy.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberB{Value: id[:]},
		},
	})
	if err != nil {
		return shared.ClimateControl{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.ClimateControl{}, false, nil
	}

	item := &shared.ClimateControl{}
	err = dynamo.UnmarshalMapWithOptions(result.Item, item)
	if err != nil {
		return shared.ClimateControl{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return *item, true, nil
}

func (r *Repository) GroupByFriendlyNamePrefix(all []shared.ClimateControl) map[string][]shared.ClimateControl {
	byP := map[string][]shared.ClimateControl{}
	for _, e := range all {
		prefix := e.GetFriendlyNamePrefix()
		es, ok := byP[prefix]
		if !ok {
			es = []shared.ClimateControl{}
		}
		es = append(es, e)
		byP[prefix] = es
	}

	return byP
}

func (r *Repository) List(ctx context.Context) ([]shared.ClimateControl, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.ClimateControl{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	items := []shared.ClimateControl{}
	for {
		result, err := dy.Scan(ctx, input)
		if err != nil {
			return []shared.ClimateControl{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := shared.ClimateControl{}
			if err = dynamo.UnmarshalMapWithOptions(i, &item); err != nil {
				return []shared.ClimateControl{}, fmt.Errorf("error unmarshaling: %s", err.Error())
			}
			items = append(items, item)
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].RawClimateControl.Attributes.FriendlyName < items[j].RawClimateControl.Attributes.FriendlyName
	})

	return items, nil
}

func (r *Repository) ListByFriendlyNamePrefix(ctx context.Context) (map[string][]shared.ClimateControl, error) {
	all, err := r.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting devices: %s", err.Error())
	}

	return r.GroupByFriendlyNamePrefix(all), nil
}

func (r *Repository) Put(ctx context.Context, item shared.ClimateControl) (shared.ClimateControl, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.ClimateControl{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	if item.ID == uuid.Nil {
		// Since an ID can easily be forgotten, let's never assume we need to create one.
		return shared.ClimateControl{}, fmt.Errorf("an ID is required")
	}

	av, err := dynamo.MarshalMapWithOptions(item)
	if err != nil {
		return shared.ClimateControl{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItem(ctx, input)
	if err != nil {
		return shared.ClimateControl{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := r.Get(ctx, item.ID)
	if err != nil {
		return shared.ClimateControl{}, err
	}
	if !ok {
		return shared.ClimateControl{}, fmt.Errorf("couldn't find entity after insert")
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
