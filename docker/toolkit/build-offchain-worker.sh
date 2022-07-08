#!/bin/bash

# ./docker/dev/create-toolkit-image.sh

docker run --rm -e GOPATH=/opt/indexer/.gopath/ -v `pwd`:/opt/indexer -ti build-indexer make cmd/algorand-indexer/offchain-worker
cp cmd/offchain-worker/offchain-worker docker/toolkit/offchain-worker

# ./docker/dev/create-indexer-image.sh