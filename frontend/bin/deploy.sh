#!/bin/bash

set -ex

# start in the script directory
cd "$(dirname "$0")"
# move up one level
cd ../

npm run build

aws --profile=mLock s3 cp build s3://mlock-site --recursive
