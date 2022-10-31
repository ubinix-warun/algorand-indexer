#!/bin/bash

# Start webhook daemon. There are various configurations controlled by
# environment variables.
#
# Configuration:
#   DISABLED          - If set start a server that returns an error instead of indexer.
#   SNAPSHOT          - snapshot to import, if set don't connect to algod.
#   PORT              - port to start indexer on.
#   ALGOD_ADDR        - host:port to connect to for algod.
#   ALGOD_TOKEN       - token to use when connecting to algod.
set -e

start_with_algod() {
  echo "Starting webhook against algod."

  for i in 1 2 3 4 5; do
    wget "${ALGOD_ADDR}"/genesis -O genesis.json && break
    echo "Algod not responding... waiting."
    sleep 15
  done

  if [ ! -f genesis.json ]; then
    echo "Failed to create genesis file!"
    exit 1
  fi

  # /tmp/webhook daemon \
  /opt/webhook/scripts/docker/webhook/webhook daemon \
    --dev-mode \
    --dummydb true \
    --server ":$PORT" \
    --algod-net "${ALGOD_ADDR}" \
    --algod-token "${ALGOD_TOKEN}" \
    --genesis "genesis.json" \
    --data-dir /tmp \
    --loglevel "debug" \
    --logfile "/tmp/indexer-log.txt" >> /tmp/command.txt
}

disabled() {
  echo "disabled!"
  # go run /tmp/disabled.go -port "$PORT" -code 400 -message "Indexer disabled for this configuration."
}

if [ ! -z "$DISABLED" ]; then
  disabled
elif [ -z "${SNAPSHOT}" ]; then
  start_with_algod
fi

sleep infinity