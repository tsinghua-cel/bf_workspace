#!/usr/bin/env sh
docker build --no-cache -t tscel/bf.prysm:v5.0.1 -f Dockerfile.prysm .
docker build --no-cache -t tscel/prysmctl:v5.0.1 -f Dockerfile.prysmctl .

#docker push tscel/bf.prysm:v5.0.1
#docker push tscel/prysmctl:v5.0.1