#!/bin/bash

# ./docker/dev/create-toolkit-image.sh

docker run --rm -e GOPATH=/opt/webhook/.gopath/ -v `pwd`:/opt/webhook -ti build-algorand /opt/webhook/scripts/docker/webhook/make.sh
cp cmd/webhook/webhook scripts/docker/webhook/webhook

# ./docker/dev/create-webhook-image.sh