#!/bin/bash

CODE_DIR="./code"

if [ ! -d "$CODE_DIR" ]; then
	echo "'$CODE_DIR' is not exist"
	exit 1
fi

cd "$CODE_DIR"

for dir in */; do
	dir=${dir%/}
	if [ -d "$dir" ]; then
		cd "$dir"
		./docker.sh
		cd ..
	fi
done
