package device

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo"
	"regexp"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const (
	tableName = "Device_v1"
)

var (
	code = regexp.MustCompile(`3(\d) 3(\d) 3(\d) 3(\d) 0A 0D`)
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

func Get(ctx context.Context, id uuid.UUID) (shared.Device, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Device{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	result, err := dy.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {B: id[:]},
		},
	})
	if err != nil {
		return shared.Device{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.Device{}, false, nil
	}

	item := &shared.Device{}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return shared.Device{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	cleanCodes(item)

	return *item, true, nil
}

func List(ctx context.Context) ([]shared.Device, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.Device{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	items := []shared.Device{}
	for {
		result, err := dy.ScanWithContext(ctx, input)
		if err != nil {
			return []shared.Device{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := shared.Device{}
			if err = dynamodbattribute.UnmarshalMap(i, &item); err != nil {
				return []shared.Device{}, fmt.Errorf("error unmarshaling: %s", err.Error())
			}
			items = append(items, item)
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].RawDevice.Name < items[j].RawDevice.Name
	})

	return items, nil
}

func ListForUnit(ctx context.Context, unit shared.Unit) ([]shared.Device, error) {
	all, err := List(ctx)
	if err != nil {
		return []shared.Device{}, fmt.Errorf("error getting devices: %s", err.Error())
	}

	forU := []shared.Device{}
	for _, d := range all {
		if d.UnitID != nil && *d.UnitID == unit.ID {
			forU = append(forU, d)
		}
	}
	return forU, nil
}

func Put(ctx context.Context, item shared.Device) (shared.Device, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.Device{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	if item.ID == uuid.Nil {
		// Since an ID can easily be forgotten, let's never assume we need to create one.
		return shared.Device{}, fmt.Errorf("an ID is required")
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return shared.Device{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItemWithContext(ctx, input)
	if err != nil {
		return shared.Device{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := Get(ctx, item.ID)
	if err != nil {
		return shared.Device{}, err
	}
	if !ok {
		return shared.Device{}, fmt.Errorf("couldn't find entity after insert")
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

func cleanCodes(in *shared.Device) {
	in.HABThing.Configuration.UsercodeCode1 = cleanCode(in.HABThing.Configuration.UsercodeCode1)
	in.HABThing.Configuration.UsercodeCode1 = cleanCode(in.HABThing.Configuration.UsercodeCode1)
	in.HABThing.Configuration.UsercodeCode2 = cleanCode(in.HABThing.Configuration.UsercodeCode2)
	in.HABThing.Configuration.UsercodeCode3 = cleanCode(in.HABThing.Configuration.UsercodeCode3)
	in.HABThing.Configuration.UsercodeCode4 = cleanCode(in.HABThing.Configuration.UsercodeCode4)
	in.HABThing.Configuration.UsercodeCode5 = cleanCode(in.HABThing.Configuration.UsercodeCode5)
	in.HABThing.Configuration.UsercodeCode6 = cleanCode(in.HABThing.Configuration.UsercodeCode6)
	in.HABThing.Configuration.UsercodeCode7 = cleanCode(in.HABThing.Configuration.UsercodeCode7)
	in.HABThing.Configuration.UsercodeCode8 = cleanCode(in.HABThing.Configuration.UsercodeCode8)
	in.HABThing.Configuration.UsercodeCode9 = cleanCode(in.HABThing.Configuration.UsercodeCode9)
	in.HABThing.Configuration.UsercodeCode10 = cleanCode(in.HABThing.Configuration.UsercodeCode10)
	in.HABThing.Configuration.UsercodeCode11 = cleanCode(in.HABThing.Configuration.UsercodeCode11)
	in.HABThing.Configuration.UsercodeCode12 = cleanCode(in.HABThing.Configuration.UsercodeCode12)
	in.HABThing.Configuration.UsercodeCode13 = cleanCode(in.HABThing.Configuration.UsercodeCode13)
	in.HABThing.Configuration.UsercodeCode14 = cleanCode(in.HABThing.Configuration.UsercodeCode14)
	in.HABThing.Configuration.UsercodeCode15 = cleanCode(in.HABThing.Configuration.UsercodeCode15)
	in.HABThing.Configuration.UsercodeCode16 = cleanCode(in.HABThing.Configuration.UsercodeCode16)
	in.HABThing.Configuration.UsercodeCode17 = cleanCode(in.HABThing.Configuration.UsercodeCode17)
	in.HABThing.Configuration.UsercodeCode18 = cleanCode(in.HABThing.Configuration.UsercodeCode18)
	in.HABThing.Configuration.UsercodeCode19 = cleanCode(in.HABThing.Configuration.UsercodeCode19)
	in.HABThing.Configuration.UsercodeCode20 = cleanCode(in.HABThing.Configuration.UsercodeCode20)
	in.HABThing.Configuration.UsercodeCode21 = cleanCode(in.HABThing.Configuration.UsercodeCode21)
	in.HABThing.Configuration.UsercodeCode22 = cleanCode(in.HABThing.Configuration.UsercodeCode22)
	in.HABThing.Configuration.UsercodeCode23 = cleanCode(in.HABThing.Configuration.UsercodeCode23)
	in.HABThing.Configuration.UsercodeCode24 = cleanCode(in.HABThing.Configuration.UsercodeCode24)
	in.HABThing.Configuration.UsercodeCode25 = cleanCode(in.HABThing.Configuration.UsercodeCode25)
	in.HABThing.Configuration.UsercodeCode26 = cleanCode(in.HABThing.Configuration.UsercodeCode26)
	in.HABThing.Configuration.UsercodeCode27 = cleanCode(in.HABThing.Configuration.UsercodeCode27)
	in.HABThing.Configuration.UsercodeCode28 = cleanCode(in.HABThing.Configuration.UsercodeCode28)
	in.HABThing.Configuration.UsercodeCode29 = cleanCode(in.HABThing.Configuration.UsercodeCode29)
	in.HABThing.Configuration.UsercodeCode30 = cleanCode(in.HABThing.Configuration.UsercodeCode30)
}

func cleanCode(in string) string {
	m := code.FindStringSubmatch(in)
	if len(m) == 0 {
		return in
	}

	return m[1] + m[2] + m[3] + m[4]
}
