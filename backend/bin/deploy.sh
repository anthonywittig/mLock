#!/bin/bash

set -ex

# start in the script directory
cd "$(dirname "$0")"
# move where we need to be
cd ../lambdas

rm -rf build
mkdir build

#lambdas=('helloworld' 'helloworld2' 'users')
lambdas=('users')
for lambda in "${lambdas[@]}" ; do
    buildDir="build/${lambda}"
    mkdir $buildDir
    GOOS=linux GOARCH=amd64 go build -o $buildDir ${lambda}/main.go
    (cd $buildDir && zip function.zip main)

    aws --profile=mLock lambda update-function-code \
        --function-name $lambda \
        --zip-file fileb://${buildDir}/function.zip
        #--handler main
        #--runtime go1.x \
        #--role arn:aws:iam::123456789012:role/execution_role
done

