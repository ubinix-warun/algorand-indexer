ARG GO_VERSION=1.18.1
FROM golang:$GO_VERSION-alpine

ENV HOME /opt/indexer
WORKDIR /opt/indexer

RUN apk add --no-cache bash libtool boost-dev

COPY start.sh /tmp/start.sh
COPY algorand-indexer /tmp/algorand-indexer

CMD ["/tmp/start.sh"]