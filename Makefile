SERVER_BINARY_NAME = server
SERVER_PACKAGE_PATH = ./cmd/server
SERVER_ADDR = localhost:8080

DB_DRIVER=clickhouse
DB_ADDR=localhost:9000
DB_NAME=events
DATABASE_DSN = $(DB_DRIVER)://$(DB_ADDR)/$(DB_NAME)
# clickhouse://postgres:postgres@localhost:9000/database?dial_timeout=200ms&max_execution_time=60

DOC_ADDR = localhost:6060

# ====================
# HELPERS
# ====================

## help: show this help message
help:
	@echo
	@echo 'usage: make <target>'
	@echo
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# ====================
# QUALITY
# ====================

## tidy: format code and tidy mod file
tidy:
	go fmt ./...
	go mod tidy -v

# ====================
# DEVELOPMENT
# ====================

## test: run all tests
test:
	go test ./...

## test/cover: run all tests and display coverage
test/cover:
	go test ./... -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out
	rm /tmp/coverage.out

## build/local: build the server locally
build/local:
	go build -o /tmp/bin/$(SERVER_BINARY_NAME) $(SERVER_PACKAGE_PATH)

## run/local: run the server locally
run/local: build/local
	/tmp/bin/$(SERVER_BINARY_NAME) -a=$(SERVER_ADDR) -dsn=$(DATABASE_DSN) -d=$(DB_DRIVER)

## build: build the server with docker-compose
build:
	docker-compose build

## run: build the server with docker-compose
run: build-docker
	docker-compose up

# ====================
# DOCUMENTATION
# ====================

## doc: generate documentation on http port
doc:
	@echo 'open http://$(DOC_ADDR)/pkg/github.com/pavlegich/events-store/?m=all'
	godoc -http=$(DOC_ADDR)

.PHONY: help tidy build/local run/local build run doc test test/cover
