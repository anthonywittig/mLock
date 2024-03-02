package miscellaneous

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type Repository struct{}

const (
	tableName = "Miscellaneous_v1"
)

var magicID = uuid.MustParse("8c957e18-d01f-4d86-a20a-32ab637793db")

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Get(ctx context.Context) (shared.Miscellaneous, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Miscellaneous{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	result, err := dy.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberB{Value: magicID[:]},
		},
	})
	if err != nil {
		return shared.Miscellaneous{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.Miscellaneous{}, false, nil
	}

	item := &shared.Miscellaneous{}
	err = dynamo.UnmarshalMapWithOptions(result.Item, item)
	if err != nil {
		return shared.Miscellaneous{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return *item, true, nil
}

func (r *Repository) Put(ctx context.Context, item shared.Miscellaneous) (shared.Miscellaneous, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Miscellaneous{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	item.ID = magicID

	av, err := dynamo.MarshalMapWithOptions(item)
	if err != nil {
		return shared.Miscellaneous{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItem(ctx, input)
	if err != nil {
		return shared.Miscellaneous{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := r.Get(ctx)
	if err != nil {
		return shared.Miscellaneous{}, err
	}
	if !ok {
		return shared.Miscellaneous{}, fmt.Errorf("couldn't find entity after insert")
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
