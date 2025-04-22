#!/usr/bin/env sh
docker build --no-cache -t tscel/bf.prysm:v4.2.1 -f Dockerfile.prysm .
docker build --no-cache -t tscel/prysmctl:v4.2.1 -f Dockerfile.prysmctl .

