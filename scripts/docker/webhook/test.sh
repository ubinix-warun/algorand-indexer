#!/bin/bash

# ./docker/dev/create-toolkit-image.sh

#   /tmp/algorand-indexer daemon \
#     --dev-mode \
#     --server ":$PORT" \
#     -P "$CONNECTION_STRING" \
#     --algod-net "${ALGOD_ADDR}" \
#     --algod-token "${ALGOD_TOKEN}" \
#     --genesis "genesis.json" \
#     --data-dir /tmp \
#     --loglevel "debug" \
#     --logfile "/tmp/indexer-log.txt" >> /tmp/command.txt

docker run --rm -e GOPATH=/opt/webhook/.gopath/ \
    -e PORT='8980' \
    -e ALGOD_ADDR='192.168.1.38:4001' \
    -e ALGOD_TOKEN='aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' \
    -v `pwd`:/opt/webhook -ti build-algorand \
    #  /opt/webhook/scripts/docker/webhook/start.sh $@
     /bin/bash


# docker run --rm -e GOPATH=/opt/webhook/.gopath/ \
#     -e PORT='8980' \
#     -e ALGOD_ADDR='192.168.1.38:4001' \
#     -e ALGOD_TOKEN='aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' \
#     -v `pwd`:/opt/webhook -ti build-algorand \
#      /opt/webhook/scripts/docker/webhook/start.sh $@

# ./docker/dev/create-webhook-image.sh