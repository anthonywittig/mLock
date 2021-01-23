package user

import (
	"context"
	"fmt"
	"log"
	"mlock/shared"
	"mlock/shared/dynamo"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UserServiceImpl struct{}

const (
	tableName = "Users"
	userType  = "1" // Not really using ATM.
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func Delete(ctx context.Context, email string) error {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	email = strings.ToLower(email)

	// No audit trail for deletes. :(

	if _, err = dy.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Type":  {S: aws.String(userType)},
			"Email": {S: aws.String(email)},
		},
		TableName: aws.String(tableName),
	}); err != nil {
		return fmt.Errorf("error deleting item: %s", err.Error())
	}

	return nil
}

func Get(ctx context.Context, email string) (shared.User, bool, error) {
	return (&UserServiceImpl{}).Get(ctx, email)
}

func (u *UserServiceImpl) Get(ctx context.Context, email string) (shared.User, bool, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error getting client: %s", err.Error())
	}

	email = strings.ToLower(email)

	result, err := dy.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Type":  {S: aws.String(userType)},
			"Email": {S: aws.String(email)},
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

func List(ctx context.Context) ([]shared.User, error) {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return []shared.User{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	/*
		filt := expression.Name("Type").Equal(expression.Value(userType))
		expr, err := expression.NewBuilder().WithFilter(filt).Build()
		if err != nil {
			return []shared.User{}, fmt.Errorf("error building expression: %s", err.Error())
		}
	*/

	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"Type": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(userType)},
				},
			},
		},
	}

	items := []shared.User{}
	for {
		// Get the list of tables
		result, err := dy.QueryWithContext(ctx, input)
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

	return items, nil
}

func Put(ctx context.Context, email string) error {
	dy, err := dynamo.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	email = strings.ToLower(email)

	if !isEmailValid(email) {
		// Should indicate it's a 4xx; we should probably do some validation on the frontend too.
		return fmt.Errorf("email isn't formatted correctly")
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return fmt.Errorf("no current user")
	}

	item := shared.User{
		Type:      userType,
		Email:     email,
		CreatedBy: currentUser.Email,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dy.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("error putting item: %s", err.Error())
	}

	return nil
}

func MigrateUsers(ctx context.Context) error {
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
				AttributeName: aws.String("Type"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Email"),
				AttributeType: aws.String("S"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Type"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Email"),
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

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
