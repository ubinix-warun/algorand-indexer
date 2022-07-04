#!/bin/bash

# ./docker/dev/create-toolkit-image.sh

docker run --rm -e GOPATH=/opt/indexer/.gopath/ -v `pwd`:/opt/indexer -ti build-indexer /opt/indexer/docker/toolkit/make.sh
cp cmd/algorand-indexer/algorand-indexer docker/toolkit/algorand-indexer

# ./docker/dev/create-indexer-image.sh