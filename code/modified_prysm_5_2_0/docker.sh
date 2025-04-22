#!/usr/bin/env sh
#docker build --no-cache -t tscel/bf.beacon:v5.2.0 -f Dockerfile.beacon .
#docker build --no-cache -t tscel/bf.validator:v5.2.0 -f Dockerfile.validator .
#docker push tscel/bf.beacon:v5.2.0
#docker push tscel/bf.validator:v5.2.0

docker build --no-cache -t tscel/bf.prysm:v5.2.0 -f Dockerfile.prysm .
docker build --no-cache -t tscel/prysmctl:v5.2.0 -f Dockerfile.prysmctl .
docker push tscel/bf.prysm:v5.2.0
docker push tscel/prysmctl:v5.2.0