FROM alpine:3.6

RUN apk add --no-cache --virtual .build-deps go git gcc openssl musl-dev && \
    go get github.com/nokamoto/esp4g/examples/ping/esp4g-ping-server && \
    apk del .build-deps

ENV PATH $PATH:/root/go/bin

ENTRYPOINT ["esp4g-ping-server"]
