FROM nokamotohub/esp4g/go

RUN go get github.com/nokamoto/esp4g/examples/calc/esp4g-calc-server

RUN apk del .build-deps

ENTRYPOINT esp4g-calc-server
