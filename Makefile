BIN := eth-indexer
RECONCILER_BIN := eth-indexer-reconciler
DB_FILE := tracker_db 
BUILD_CONF := CGO_ENABLED=1 GOOS=linux GOARCH=amd64
BUILD_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null)
DEBUG := DEV=true

.PHONY: build build-reconciler run run-reconciler run-bootstrap clean clean-debug

clean:
	rm ${BIN}

clean-db:
	rm ${DB_FILE}

build:
	${BUILD_CONF} go build -ldflags="-X main.build=${BUILD_COMMIT} -s -w" -o ${BIN} cmd/service/*.go

build-reconciler:
	${BUILD_CONF} go build -ldflags="-X main.build=${BUILD_COMMIT} -s -w" -o ${RECONCILER_BIN} cmd/reconciler/*.go

run:
	${BUILD_CONF} ${DEBUG} go run cmd/service/*.go

run-reconciler:
	${BUILD_CONF} ${DEBUG} go run cmd/reconciler/*.go