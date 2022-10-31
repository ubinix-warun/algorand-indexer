#!/bin/bash

cd scripts/docker
docker build -t build-algorand -f Dockerfile.Toolchain .