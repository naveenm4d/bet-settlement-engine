BIN?=bet-settlement-engine

default: run
.PHONY : build run

lint:
	golangci-lint run -c .golangci.yml

build:
	go build -o build/${BIN}

run: build
	./build/${BIN}

build-mocks:
	cd mocks/ && rm -rf -- */ && mockery

test:
	go test ./... -tags musl -coverprofile=coverage.txt -covermode count = 1