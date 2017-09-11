#!/bin/bash

check () {
    docker-compose up -d

    docker-compose exec $1 echo ok

    curl localhost:8080/metrics

    for service in esp4g prometheus grafana
    do
        docker-compose exec $service echo ok
    done

    docker-compose stop
}

set -ex

cd `dirname $0`

cd ../examples/ping
check ping

cd ../calc
check calc

cd ../benchmark
check benchmark
