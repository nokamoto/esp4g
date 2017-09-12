#!/bin/bash

set -ex

cd `dirname $0`

wget https://raw.githubusercontent.com/nokamoto/esp4g/master/examples/ping/protobuf/service.proto -O service.proto

protoc -I. --include_imports --include_source_info ./service.proto --descriptor_set_out ./descriptor.pb

cat << EOF > config.yaml
logs:
  zap:
    level: info
    encoding: json
    outputPaths:
      - stdout
    errorOutputPaths:
      - stderr

usage:
  rules:
    - selector: /eps4g.ping.PingService/Send
      allow_unregistered_calls: true
EOF

esp4g -c ./config.yaml -d ./descriptor.pb
