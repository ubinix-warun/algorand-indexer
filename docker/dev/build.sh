#!/bin/bash

# ./docker/dev/toolchain.sh

docker run --rm -e GOPATH=/opt/indexer/.gopath/ -v `pwd`:/opt/indexer -ti build-indexer make
cp cmd/algorand-indexer/algorand-indexer docker/dev/algorand-indexer

# ./docker/dev/create-image.sh