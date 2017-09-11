# esp4g

[![CircleCI](https://circleci.com/gh/nokamoto/esp4g/tree/master.svg?style=svg)](https://circleci.com/gh/nokamoto/esp4g/tree/master)

- [Overview](#overview)
- [Quickstart](#quickstart)
- [Examples](#examples)

## Overview
![Overview](/.md/overview.png)

Supports:

- Unary/Stream gRPC proxy
- Access log
  - with logging
  - with prometheus exporter
- Access control
  - with API keys

## Quickstart

Install:

```bash
go get github.com/nokamoto/esp4g/esp4g
go get github.com/nokamoto/esp4g/esp4g-extension
```

and the example gRPC server and client:

```bash
go get github.com/nokamoto/esp4g/examples/ping/esp4g-ping
go get github.com/nokamoto/esp4g/examples/ping/esp4g-ping-server
```

Configure:

descriptor.pb

```bash
wget https://raw.githubusercontent.com/nokamoto/esp4g/master/examples/ping/protobuf/service.proto -O service.proto

protoc -I. --include_imports --include_source_info ./service.proto --descriptor_set_out ./descriptor.pb
```

config.yaml

```bash
cat << EOF > config.yaml
logs:
  logging: true

usage:
  rules:
    - selector: /eps4g.ping.PingService/Send
      allow_unregistered_calls: true
EOF
```

Run:

```bash
esp4g-ping-server

esp4g-extension -c ./config.yaml -d ./descriptor.pb

esp4g -d ./descriptor.pb
```

Access _PingService_ through the esp:

```bash
esp4g-ping
```

## Examples

To run docker-compose, make _descriptor.pb_ files into _examples_ directories.

```
make
```

Run:

```
cd examples/ping
docker-compose up
```

### Example Services
| Service | Port | |
| --- | --- | --- |
| esp4g | 9000 | esp4g-ping |
| esp4g-extension | 8080 | http://localhost:8080/metrics |
| ping | 8000 | esp4g-ping -p 8000 |
| prometheus | 9090 | http://localhost:9090 |
| grafana | 3000 | http://localhost:3000 |

#### Example Monitoring
Setup Grafana [data sources](http://localhost:3000/datasources).

| Data Source | Type | Url | Access |
| --- | --- | --- | --- |
| esp4g | Prometheus | http://prometheus.esp4g:9090 | proxy |


Add any query to the [dashboard](http://localhost:3000/dashboard).

![gRPC Latency Quantile](/.md/latency.png)

```
histogram_quantile(0.99, rate(grpc_response_time_seconds_bucket[1m])) * 1000
```

![gRPC rps](/.md/rps.png)

```
rate(grpc_request_size_bytes_count[10s])
```
