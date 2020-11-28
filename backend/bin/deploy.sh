#!/bin/bash

set -ex

handleLambda () {
    lambda=$1
    buildDir="build/${lambda}"

    mkdir $buildDir
    GOOS=linux GOARCH=amd64 go build -o $buildDir ${lambda}/main.go
    (cd $buildDir && zip function.zip main)
    aws --profile=mLock lambda update-function-code \
        --function-name $lambda \
        --zip-file fileb://${buildDir}/function.zip
}

scriptDir=$(dirname "$0")
cd $scriptDir
scriptDir=$(pwd)

# handle the DB
cd $scriptDir
cd ../db

# handle the apis
cd $scriptDir
cd ../apis

rm -rf build
mkdir build

#lambdas=('helloworld' 'helloworld2' 'users')
lambdas=('users')
for lambda in "${lambdas[@]}" ; do
    handleLambda $lambda
done

