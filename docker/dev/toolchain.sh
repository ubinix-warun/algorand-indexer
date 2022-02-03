#!/bin/bash

cd docker/dev
docker build -t build-indexer -f Dockerfile.Toolchain .