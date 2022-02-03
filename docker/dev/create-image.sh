#!/bin/bash

cd docker/dev
docker build -t algorand-indexer -f Dockerfile .
