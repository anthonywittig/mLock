package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

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

	if err := createIfNotExists(ctx, svc, name); err != nil {
		log.Fatalf("unable to create if not exists, %v", err)
	}

	if err := updateCode(ctx, svc, name, lambdaDir, projectDir); err != nil {
		log.Fatalf("unable to update code, %v", err)
	}
}

func createIfNotExists(ctx context.Context, svc *lambda.Client, name string) error {
	listResp, err := svc.ListFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return fmt.Errorf("unable to list functions, %v", err)
	}

	for _, fn := range listResp.Functions {
		if *fn.FunctionName == name {
			// Already exists, nothing to do.
			return nil
		}
	}

	return fmt.Errorf("function %s does not exist, need to implement creation", name)
}

func updateCode(ctx context.Context, svc *lambda.Client, name string, lambdaDir string, projectDir string) error {
	// Create the build directory.
	if err := mkDir(fmt.Sprintf("%s/build", projectDir)); err != nil {
		return fmt.Errorf("unable to create directory, %v", err)
	}
	buildDir := fmt.Sprintf("%s/build/%s", projectDir, name)
	if err := rmDir(buildDir); err != nil {
		return fmt.Errorf("unable to remove directory, %v", err)
	}
	if err := mkDir(buildDir); err != nil {
		return fmt.Errorf("unable to create directory, %v", err)
	}

	if err := cpIfExists(fmt.Sprintf("%s/.env", lambdaDir), buildDir); err != nil {
		return fmt.Errorf("unable to copy env file, %v", err)
	}

	if err := buildLambdaBinary(lambdaDir, buildDir); err != nil {
		return fmt.Errorf("unable to build go binary, %v", err)
	}

	if err := createLambdaZip(buildDir); err != nil {
		return fmt.Errorf("unable to create lambda zip, %v", err)
	}

	zipContents, err := ioutil.ReadFile(fmt.Sprintf("%s/function.zip", buildDir))
	if err != nil {
		return fmt.Errorf("error reading zip file, %v", err)
	}

	// Upload
	if _, err := svc.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		Publish:      true,
		FunctionName: &name,
		ZipFile:      zipContents,
	}); err != nil {
		return fmt.Errorf("unable to update code, %v", err)
	}

	return nil
}
