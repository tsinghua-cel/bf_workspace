#!/usr/bin/env sh
docker build --no-cache -t tscel/bf.prysm:v5.2.0 -f Dockerfile.prysm .
docker build --no-cache -t tscel/prysmctl:v5.2.0 -f Dockerfile.prysmctl .
