#!/bin/bash

set -ex

# start in the script directory
cd "$(dirname "$0")"
# move up one level
cd ../

npm install
CI=true npm run build
CI=true npm test
