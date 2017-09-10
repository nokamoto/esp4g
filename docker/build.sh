#!/bin/bash

set -ex

cd `dirname $0`

docker build -f go.Dockerfile -t nokamotohub/esp4g/go .

for package in {esp4g,esp4g-extension,esp4g-examples-ping,esp4g-examples-calc,esp4g-examples-benchmark}
do
    docker build -f ${package}.Dockerfile -t nokamotohub/$package .

    docker push nokamotohub/$package
done
