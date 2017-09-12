# esp4g

[![CircleCI](https://circleci.com/gh/nokamoto/esp4g/tree/master.svg?style=svg)](https://circleci.com/gh/nokamoto/esp4g/tree/master)

- [Overview](#overview)
- [Installing](#installing)
- [Quickstart](#quickstart)
- [Examples](#examples)
- [Configuratin](#configuration)

## Overview
![Overview](/.md/overview.png)

Supports:

- Unary/Stream gRPC proxy
- Access log
  - with logging
  - with prometheus exporter
- Access control
  - with API keys

## Installing
### Using Docker
```
docker run \
    -v /path/to/your/config.yaml:/config.yaml \
    -v /path/to/your/descriptor.pb:/descriptor.pb \
    -p 9000:9000 \
    nokamotohub/esp4g -c /config.yaml -d /descriptor.pb -proxy [host:port]
```

To make _descriptor.pb_ file, run _protoc_ with `--descriptor_set_out` option.

This starts esp to forward incoming gRPC requests to `-proxy` address.

See also [Configuration](#configuration).

### From Source
```
go get github.com/nokamoto/esp4g/esp4g
esp4g -c /path/to/your/config.yaml -d /path/to/your/descriptor.pb
```

## Quickstart

Install:

```bash
go get github.com/nokamoto/esp4g/esp4g
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

esp4g -c ./config.yaml -d ./descriptor.pb
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
histogram_quantile(0.99, rate(grpc_response_seconds_bucket[1m])) * 1000
```

![gRPC rps](/.md/rps.png)

```
rate(grpc_request_bytes_count[10s])
```

## Configuration
### esp4g
```
$ esp4g -h
Usage of esp4g:
  -c string
    	The application config file (default "./config.yaml")
  -d string
    	FileDescriptorSet protocol buffer file (default "./descriptor.pb")
  -e string
    	The gRPC extension service address (default: in-process)
  -p int
    	The gRPC server port (default 9000)
  -u string
    	The gRPC upstream address (default "localhost:8000")
```

### config.yaml
```yaml
logs:
  # If 'logging' is set `true`, gRPC access logs to stdout.
  # default: false
  logging: true

  prometheus:
    # The port to scrape Prometheus metrics from. If 'port' is not supplied, Prometheus exporter is not available.
    # default: nil
    port: 8080

    histograms:
      # The histogram for gRPC latency.
      # default: nil
      response_seconds:
        name: "grpc_response_seconds"
        help: "gRPC latency distributions."
        buckets: [0.001, 0.002, 0.003, 0.004, 0.005, 0.006, 0.007, 0.008, 0.009, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]

      # The histogram for gRPC request size.
      # default: nil
      request_bytes:
        name: "grpc_request_bytes"
        help: "gRPC request content size distributions."
        buckets: [1, 2, 4, 8, 16, 32, 64, 128]

      # The histogram for gRPC response size.
      # default: nil
      response_bytes:
        name: "grpc_response_bytes"
        help: "gRPC response content size distributions."
        buckets: [1, 2, 4, 8, 16, 32, 64, 128]

usage:
  # The gRPC access control list.
  # Esp forwards an incoming request to the upstream server only if it satisfies one of the following rules,
  # otherwise esp sends back `Unauthenticated` status code.
  rules:
      # The gRPC service method name.
    - selector: "/esp4g.benchmark.UnaryService/Send"

      # Esp always allows any incoming request if 'allow_unregistered_calls' is `true`.
      # default: false
      allow_unregistered_calls: true

    - selector: "/eps4g.ping.PingService/Send"

      # Esp allows an incoming request if 'x-api-key' metadata contains one of the following keys.
      # default: nil
      registered_api_keys:
        - guest
```

#### Histogram configuration
[prometheus/client_golang](https://github.com/prometheus/client_golang/blob/3eb912b336976e0f66a5eb9a434adfbba8dff646/prometheus/histogram.go#L113-L154)
