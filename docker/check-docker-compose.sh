#!/bin/bash

set -ex

cd `dirname $0`

for package in {ping,calc,benchmark}
do
    cd ../examples/$package

    docker-compose up -d
    docker-compose stop
done
