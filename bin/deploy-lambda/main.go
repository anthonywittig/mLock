package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type lambdaDeployer struct {
	cfg        aws.Config
	lambdaDir  string
	name       string
	projectDir string
	svc        *lambda.Client
}

func main() {
	ctx := context.Background()
	relativePath := os.Args[1]

	projectDir, err := filepath.Abs("./../../")
	if err != nil {
		log.Fatalf("unable to create project directory path, %v", err)
	}

	lambdaDir := fmt.Sprintf("%s/%s", projectDir, relativePath)
	if err != nil {
		log.Fatalf("unable to create file path, %v", err)
	}
	if _, err := os.Stat(lambdaDir); err != nil {
		log.Fatalf("unable to stat lambda directory, %v", err)
	}

	name := path.Base(lambdaDir)
	fmt.Printf("Deploying lambda: %s\n", name)

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedConfigProfile("mLock-dev"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := lambda.NewFromConfig(cfg)

	l := lambdaDeployer{
		cfg:        cfg,
		lambdaDir:  lambdaDir,
		name:       name,
		projectDir: projectDir,
		svc:        svc,
	}

	if err := l.createIfNotExists(ctx); err != nil {
		log.Fatalf("unable to create if not exists, %v", err)
	}

	if err := l.updateCode(ctx); err != nil {
		log.Fatalf("unable to update code, %v", err)
	}
}

func (l *lambdaDeployer) buildAndGetZip() ([]byte, error) {
	if err := mkDir(fmt.Sprintf("%s/build", l.projectDir)); err != nil {
		return []byte{}, fmt.Errorf("unable to create directory, %v", err)
	}
	buildDir := fmt.Sprintf("%s/build/%s", l.projectDir, l.name)
	if err := rmDir(buildDir); err != nil {
		return []byte{}, fmt.Errorf("unable to remove directory, %v", err)
	}
	if err := mkDir(buildDir); err != nil {
		return []byte{}, fmt.Errorf("unable to create directory, %v", err)
	}

	if err := cpIfExists(fmt.Sprintf("%s/.env", l.lambdaDir), buildDir); err != nil {
		return []byte{}, fmt.Errorf("unable to copy env file, %v", err)
	}

	if err := buildLambdaBinary(l.lambdaDir, buildDir); err != nil {
		return []byte{}, fmt.Errorf("unable to build go binary, %v", err)
	}

	if err := createLambdaZip(buildDir); err != nil {
		return []byte{}, fmt.Errorf("unable to create lambda zip, %v", err)
	}

	zipContents, err := ioutil.ReadFile(fmt.Sprintf("%s/function.zip", buildDir))
	if err != nil {
		return []byte{}, fmt.Errorf("error reading zip file, %v", err)
	}

	return zipContents, nil
}

func (l *lambdaDeployer) createIfNotExists(ctx context.Context) error {
	listResp, err := l.svc.ListFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return fmt.Errorf("unable to list functions, %v", err)
	}

	for _, fn := range listResp.Functions {
		if *fn.FunctionName == l.name {
			return nil
		}
	}

	zipContents, err := l.buildAndGetZip()
	if err != nil {
		return fmt.Errorf("unable to build and get zip, %v", err)
	}

	role, err := l.createRoleIfNotExists(ctx)
	if err != nil {
		return fmt.Errorf("unable to create role if not exists, %v", err)
	}

	if _, err := l.svc.CreateFunction(
		ctx,
		&lambda.CreateFunctionInput{
			Code: &types.FunctionCode{
				ZipFile: zipContents,
			},
			FunctionName: &l.name,
			Role:         role,
			Handler:      aws.String("main"),
			PackageType:  "Zip",
			Publish:      true,
			Runtime:      "go1.x",
			Timeout:      aws.Int32(30), // Seems like an ok default, but some will need more.
		},
	); err != nil {
		return fmt.Errorf("unable to create function -- YOU MAY NEED TO WAIT 30 SECONDS AND TRY AGAIN --, %v", err)
	}

	return nil
}

func (l *lambdaDeployer) createRoleIfNotExists(ctx context.Context) (*string, error) {
	iamSvc := iam.NewFromConfig(l.cfg)

	roleName := fmt.Sprintf("%s-lambda-role", l.name)

	list, err := iamSvc.ListRoles(ctx, &iam.ListRolesInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to list roles, %v", err)
	}

	for _, r := range list.Roles {
		if *r.RoleName == roleName {
			return r.Arn, nil
		}
	}

	// In reality we'll need more permissions than this.
	create, err := iamSvc.CreateRole(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal":{
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`),
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create role, %v", err)
	}

	return create.Role.Arn, nil
}

func (l *lambdaDeployer) updateCode(ctx context.Context) error {
	zipContents, err := l.buildAndGetZip()
	if err != nil {
		return fmt.Errorf("unable to build and get zip, %v", err)
	}

	if _, err := l.svc.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		Publish:      true,
		FunctionName: &l.name,
		ZipFile:      zipContents,
	}); err != nil {
		return fmt.Errorf("unable to update code, %v", err)
	}

	return nil
}
