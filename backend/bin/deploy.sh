#!/bin/bash

set -ex

handleLambda () {
    lambda=$1
    buildDir="build/${lambda}"

    mkdir -p $buildDir # some might throw a few extra files in here before they call this function
    if test -f "${lambda}/.env"; then
        cp ${lambda}/.env $buildDir
    fi
    GOOS=linux GOARCH=amd64 go build -o $buildDir ${lambda}/main.go
    (cd $buildDir && zip -r function.zip .env *) # Couldn't get just the hidden files without pulling in "." and "..".
    aws --profile=mLock lambda update-function-code \
        --function-name $lambda \
        --zip-file fileb://${buildDir}/function.zip
}

handleDb () {
    cd $scriptDir
    cd ../lambdas/db

    rm -rf build
    mkdir build

    #mkdir -p build/migrations/goosemigrations
    #cp -r migrations/goosemigrations/ build/migrations/goosemigrations

    handleLambda migrations
    # TODO: execute lambda with "forward" option to run the migrations forward.
}

handleApis () {
    cd $scriptDir
    cd ../lambdas/apis

    rm -rf build
    mkdir build

    lambdas=('users' 'signin' 'properties' 'units')
    for lambda in "${lambdas[@]}" ; do
        handleLambda $lambda
    done
}

scriptDir=$(dirname "$0")
cd $scriptDir
scriptDir=$(pwd)

go vet ../...
go test ../...

handleDb
handleApis
