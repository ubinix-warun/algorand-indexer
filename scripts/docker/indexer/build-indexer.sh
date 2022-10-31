#!/bin/bash

# ./docker/dev/create-toolkit-image.sh

docker run --rm -e GOPATH=/opt/indexer/.gopath/ -v `pwd`:/opt/indexer -ti build-algorand /opt/indexer/scripts/docker/indexer/make.sh
cp cmd/algorand-indexer/algorand-indexer scripts/docker/indexer/algorand-indexer

# ./docker/dev/create-indexer-image.sh