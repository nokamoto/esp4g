FROM nokamotohub/esp4g/go

RUN go get github.com/nokamoto/esp4g/examples/ping/esp4g-ping-server

RUN apk del .build-deps

ENTRYPOINT esp4g-ping-server
