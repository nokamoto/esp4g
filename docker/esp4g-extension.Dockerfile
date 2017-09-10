FROM nokamotohub/esp4g/go

RUN go get github.com/nokamoto/esp4g/esp4g-extension

RUN apk del .build-deps

ENTRYPOINT esp4g-extension
