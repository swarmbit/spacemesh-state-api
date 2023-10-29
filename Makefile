export CGO_ENABLED := 1
export CGO_CFLAGS := $(CGO_CFLAGS) -DSQLITE_ENABLE_DBSTAT_VTAB=1
BIN_DIR ?= $(PROJ_DIR)../build/

build: server
.PHONY: build

server:
	cd server; go build -o $(BIN_DIR)$@ .
.PHONY: server

docker-build-api:
	docker build -t ghcr.io/swarmbit/spacemesh-state-api-v2:v2.2.7 .

docker-push-api:
	docker push ghcr.io/swarmbit/spacemesh-state-api-v2:v2.2.7