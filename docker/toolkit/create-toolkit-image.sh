#!/bin/bash

cd docker/toolkit
docker build -t build-indexer -f Dockerfile.Toolchain .