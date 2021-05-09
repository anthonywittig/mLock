package dynamo

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetClient(ctx context.Context) (*dynamodb.DynamoDB, error) {
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.DY != nil {
		return cd.DY, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		return nil, fmt.Errorf("error getting aws session: %s", err.Error())
	}

	cd.DY = dynamodb.New(sess)
	return cd.DY, nil
}

func TableExists(ctx context.Context, name string) (bool, error) {
	tables, err := listTables(ctx)
	if err != nil {
		return false, fmt.Errorf("error listing tables: %s", err.Error())
	}

	for _, t := range tables {
		if t == name {
			return true, nil
		}
	}

	return false, nil
}

func listTables(ctx context.Context) ([]string, error) {
	dy, err := GetClient(ctx)
	if err != nil {
		return []string{}, fmt.Errorf("error getting client: %s", err.Error())
	}

	input := &dynamodb.ListTablesInput{}

	tables := []string{}
	for {
		// Get the list of tables
		result, err := dy.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				return []string{}, fmt.Errorf("error calling dynamo: %s -- %s", aerr.Error(), err.Error())
			} else {
				return []string{}, fmt.Errorf("error calling dynamo: %s", err.Error())
			}
		}

		for _, n := range result.TableNames {
			tables = append(tables, *n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return tables, nil
}
