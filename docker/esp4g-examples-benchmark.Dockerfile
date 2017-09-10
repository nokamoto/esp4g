FROM nokamotohub/esp4g/go

RUN go get github.com/nokamoto/esp4g/examples/benchmark/esp4g-benchmark-server

RUN apk del .build-deps

ENTRYPOINT esp4g-benchmark-server
