

```
docker run --rm -e GOPATH=/opt/indexer/.gopath/ -v `pwd`:/opt/indexer -ti build-indexer make cmd/algorand-indexer/offchain-worker


docker run --rm -v `pwd`:/opt/indexer -ti build-indexer ./cmd/offchain-worker/offchain-worker --addr=192.168.1.107:1323

docker run --rm -v `pwd`:/opt/indexer -ti build-indexer ./cmd/offchain-worker/offchain-worker --addr=10.42.0.87:1323

```