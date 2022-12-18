#!/bin/bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR

cd ../backend
go generate ./...
go vet ./...
go test ./...
cd $SCRIPT_DIR

./deploy-lambda/run.sh backend/lambdas/apis/devices
./deploy-lambda/run.sh backend/lambdas/jobs/pollschedules
./deploy-lambda/run.sh backend/lambdas/apis/units
./deploy-lambda/run.sh backend/lambdas/apis/users
./deploy-lambda/run.sh backend/lambdas/apis/signin
./deploy-lambda/run.sh backend/lambdas/apis/properties
./deploy-lambda/run.sh backend/lambdas/db/migrations
