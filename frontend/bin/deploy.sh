#!/bin/bash

set -ex

# start in the script directory
cd "$(dirname "$0")"
# move up one level
cd ../

./bin/build.sh

aws --profile=mLock-dev s3 cp build s3://mlock-site --recursive

# clear the cache
aws --profile=mLock-dev cloudfront create-invalidation --distribution-id E3RMAC8N7J8VMP --paths "/*"
