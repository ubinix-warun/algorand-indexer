#!/bin/bash

cd scripts/docker/indexer
docker build -t algorand-indexer -f Dockerfile .