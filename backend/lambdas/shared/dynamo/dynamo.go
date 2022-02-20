package dynamo

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetClient(ctx context.Context) (*dynamodb.Client, error) {
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.DY != nil {
		return cd.DY, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-1"))
	if err != nil {
		return nil, fmt.Errorf("error getting aws config: %s", err.Error())
	}

	cd.DY = dynamodb.NewFromConfig(cfg)

	return cd.DY, nil
}

func MarshalMapWithOptions(in interface{}) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMapWithOptions(
		in,
		func(x *attributevalue.EncoderOptions) {
			x.TagKey = "json"
		},
	)
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

func UnmarshalMapWithOptions(m map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.UnmarshalMapWithOptions(
		m,
		out,
		func(x *attributevalue.DecoderOptions) {
			x.TagKey = "json"
		},
	)
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
		result, err := dy.ListTables(ctx, input)
		if err != nil {
			return []string{}, fmt.Errorf("error calling dynamo: %s", err.Error())
		}

		tables = append(tables, result.TableNames...)

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
