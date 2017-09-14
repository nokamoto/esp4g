#!/bin/bash

check () {
    docker-compose up -d

    docker-compose exec $1 echo ok

    curl localhost:8080/metrics

    for service in esp4g prometheus grafana
    do
        docker-compose exec $service echo ok
    done

    if [ "$2" != "" ]
    then
        $2 -test
    fi

    docker-compose stop
}

set -ex

cd `dirname $0`

cd ../examples/ping
check ping esp4g-ping

cd ../calc
check calc esp4g-calc

cd ../benchmark
check benchmark
