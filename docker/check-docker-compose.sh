#!/bin/bash

set -ex

cd `dirname $0`

cd ../examples/ping

docker-compose up -d
docker-compose stop

cd ../calc
docker-compose up -d
docker-compose stop

cd ../benchmark
docker-compose up -d
docker-compose stop
