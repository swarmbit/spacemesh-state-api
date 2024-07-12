export CGO_ENABLED := 1
export CGO_CFLAGS := $(CGO_CFLAGS) -DSQLITE_ENABLE_DBSTAT_VTAB=1
BIN_DIR ?= $(PROJ_DIR)../build/

build: server
.PHONY: build

build-mainnet-accounts: mainnet_accounts
.PHONY: build-mainnet-accounts

mainnet_accounts:
	cd mainnet_accounts; go build -o $(BIN_DIR)$@ .
.PHONY: mainnet_accounts

server:
	cd server; go build -o $(BIN_DIR)$@ .
.PHONY: server

run-local: build
	./build/server ./local/config.json

docker-build-api:
	docker build -t ghcr.io/swarmbit/spacemesh-state-api-v2:v2.4.2 .

docker-push-api: docker-build-api
	docker push ghcr.io/swarmbit/spacemesh-state-api-v2:v2.4.2
