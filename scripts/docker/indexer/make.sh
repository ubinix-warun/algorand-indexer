#!/bin/bash

cd /opt/indexer
git config --global --add safe.directory /opt/indexer && make
#go get github.com/gorilla/websocket