#!/bin/bash

git config --global --add safe.directory /opt/indexer/third_party/go-algorand && make
go get github.com/gorilla/websocket