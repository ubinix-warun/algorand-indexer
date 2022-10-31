#!/bin/bash

cd /opt/webhook
git config --global --add safe.directory /opt/webhook && make cmd/webhook/webhook
#go get github.com/gorilla/websocket