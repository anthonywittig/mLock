package user

import (
	"context"
	"fmt"
	"log"
	"mlock/shared"
	"mlock/shared/dynamo"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type UserService struct {
	tableName string
}

const (
	currentTableName = "Users_v2"
	lastTableName    = "Users"
	userType         = "1" // kill soon
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewUserService() *UserService {
	return &UserService{tableName: currentTableName}
}

func (u *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	// No audit trail for deletes. :(

	if _, err = dy.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {B: id[:]},
		},
		TableName: aws.String(u.tableName),
	}); err != nil {
		return fmt.Errorf("error deleting item: %s", err.Error())
	}

	return nil
}

func (u *UserService) Get(ctx context.Context, id uuid.UUID) (shared.User, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	result, err := dy.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(u.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {B: id[:]},
		},
	})
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error getting item: %s", err.Error())
	}
	if result.Item == nil {
		return shared.User{}, false, nil
	}

	item := shared.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return item, true, nil
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (shared.User, bool, error) {
	email = strings.ToLower(email)

	users, err := u.List(ctx)
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error getting users: %s", err.Error())
	}

	for _, u := range users {
		if u.Email == email {
			return u, true, nil
		}
	}

	return shared.User{}, false, nil
}

func (u *UserService) List(ctx context.Context) ([]shared.User, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.User{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	/*
		input := &dynamodb.QueryInput{
			TableName: aws.String(u.tableName),
			KeyConditions: map[string]*dynamodb.Condition{
				"Type": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dynamodb.AttributeValue{
						{S: aws.String(userType)},
					},
				},
			},
		}
	*/

	input := &dynamodb.ScanInput{
		TableName: aws.String(u.tableName),
	}

	items := []shared.User{}
	for {
		// Get the list of tables
		//result, err := dy.QueryWithContext(ctx, input)
		result, err := dy.ScanWithContext(ctx, input)
		if err != nil {
			return []shared.User{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := shared.User{}
			if err = dynamodbattribute.UnmarshalMap(i, &item); err != nil {
				return []shared.User{}, fmt.Errorf("error unmarshaling: %s", err.Error())
			}
			items = append(items, item)
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Email < items[j].Email
	})

	return items, nil
}

func (u *UserService) ListOld(ctx context.Context) ([]shared.User, error) {
	type User2 struct {
		Type      string
		Email     string
		CreatedBy string
		UpdatedBy string
	}

	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.User{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(u.tableName),
	}

	items := []shared.User{}
	for {
		// Get the list of tables
		//result, err := dy.QueryWithContext(ctx, input)
		result, err := dy.ScanWithContext(ctx, input)
		if err != nil {
			return []shared.User{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		for _, i := range result.Items {
			item := User2{}
			if err = dynamodbattribute.UnmarshalMap(i, &item); err != nil {
				return []shared.User{}, fmt.Errorf("error unmarshaling: %s", err.Error())
			}
			items = append(items, shared.User{
				ID:        uuid.New(),
				Type:      item.Type,
				Email:     item.Email,
				CreatedBy: item.CreatedBy,
				UpdatedBy: item.CreatedBy,
			})
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			break
		}
	}

	return items, nil
}

func (u *UserService) Put(ctx context.Context, user shared.User) (shared.User, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.User{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	if user.ID == uuid.Nil {
		// Since an ID can easily be forgotten, let's never assume we need to create one.
		return shared.User{}, fmt.Errorf("an ID is required")
	}

	user.Email = strings.ToLower(user.Email)

	if !isEmailValid(user.Email) {
		// Should indicate it's a 4xx; we should probably do some validation on the frontend too.
		return shared.User{}, fmt.Errorf("email isn't formatted correctly")
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.User{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.User{}, fmt.Errorf("no current user")
	}

	user.UpdatedBy = currentUser.Email
	// CreatedBy is deprecated, but treat it as UpdatedBy.
	user.CreatedBy = user.UpdatedBy

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return shared.User{}, fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(u.tableName),
	}

	_, err = dy.PutItemWithContext(ctx, input)
	if err != nil {
		return shared.User{}, fmt.Errorf("error putting item: %s", err.Error())
	}

	entity, ok, err := u.Get(ctx, user.ID)
	if err != nil {
		return shared.User{}, err
	}
	if !ok {
		return shared.User{}, fmt.Errorf("couldn't find entity after insert")
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
	u := NewUserService()
	exists, err := dynamo.TableExists(ctx, u.tableName)
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
		TableName: aws.String(u.tableName),
	}

	result, err := dy.CreateTableWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("error creating table: %s", err.Error())
	}

	log.Printf("created table: %s - %+v", u.tableName, result)

	return nil
}

func migrateData(ctx context.Context) error {
	oService := NewUserService()
	oService.tableName = lastTableName

	users, err := oService.ListOld(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %s", err.Error())
	}

	nService := NewUserService()
	for _, u := range users {
		cd, err := shared.GetContextData(ctx)
		if err != nil {
			return fmt.Errorf("error getting context data: %s", err.Error())
		}
		cd.User = &shared.User{Email: u.Email, CreatedBy: u.CreatedBy}
		if _, err := nService.Put(ctx, shared.User{
			ID:    uuid.New(),
			Email: u.Email,
		}); err != nil {
			return fmt.Errorf("error getting context data: %s", err.Error())
		}
	}

	return nil
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
