FROM alpine:3.6

RUN apk add --no-cache --virtual .build-deps go git gcc openssl musl-dev

ENV PATH $PATH:/root/go/bin
