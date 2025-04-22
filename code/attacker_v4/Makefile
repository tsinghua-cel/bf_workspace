.PHONY: default attacker reward all clean docker

GOBIN = $(shell pwd)/build/bin

VERSION := "v4"
COMMIT_SHA1 := "aabb"
AppName := attacker

default: attacker

all: attacker reward docker

BUILD_FLAGS = -tags netgo -ldflags "\
	-X github.com/tsinghua-cel/attacker-service/versions.AppName=${AppName} \
	-X github.com/tsinghua-cel/attacker-service/versions.TagVersion="v4" \
	-X 'github.com/tsinghua-cel/attacker-service/versions.BuildTime=`date`' \
	-X github.com/tsinghua-cel/attacker-service/versions.CommitSha1=${COMMIT_SHA1}  \
	-X 'github.com/tsinghua-cel/attacker-service/versions.GoVersion=`go version`' \
	-X 'github.com/tsinghua-cel/attacker-service/versions.GitBranch="v4"' \
	"

attacker:
	go build $(BUILD_FLAGS) -o=${GOBIN}/$@ -gcflags "all=-N -l" ./cmd/attacker
	@echo "Done building."

reward:
	go build $(BUILD_FLAGS) -o=${GOBIN}/$@ -gcflags "all=-N -l" ./cmd/rewards
	@echo "Done building."

docs:
	@swag init -g ./openapi/server.go

clean:
	rm -fr build/*

docker:
	docker build --no-cache -t tscel/bunnyfinder:capella .
