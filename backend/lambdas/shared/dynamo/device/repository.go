package device

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo"
	"mlock/lambdas/shared/dynamo/auditlog"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Repository struct{}

const (
	tableName = "Device_v1"
)

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) AppendToAuditLog(ctx context.Context, device shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) error {
	if len(managedLockCodes) > 0 {
		al, exists, err := auditlog.Get(ctx, device.ID)
		if err != nil {
			return fmt.Errorf("error getting audit log: %s", err.Error())
		}

		if !exists {
			al = shared.AuditLog{ID: device.ID}
		}

		for _, mlc := range managedLockCodes {
			al.Entries = append(
				al.Entries,
				shared.AuditLogEntry{
					CreatedAt: time.Now(),
					Log:       fmt.Sprintf("Code: %s; Start: %s; End: %s; Note: %s.", mlc.Code, mlc.StartAt.Format(time.RFC3339), mlc.EndAt.Format(time.RFC3339), mlc.Note),
				},
			)
		}

		// It'd be nice to tie this to the `device.Put`.
		if _, err := auditlog.Put(ctx, al); err != nil {
			return fmt.Errorf("error putting audit log: %s", err.Error())
		}
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
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

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (shared.Device, bool, error) {
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

	return *item, true, nil
}

func (r *Repository) List(ctx context.Context) ([]shared.Device, error) {
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

func (r *Repository) ListForUnit(ctx context.Context, unit shared.Unit) ([]shared.Device, error) {
	all, err := r.List(ctx)
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

func (r *Repository) Put(ctx context.Context, item shared.Device) (shared.Device, error) {
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

	entity, ok, err := r.Get(ctx, item.ID)
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
