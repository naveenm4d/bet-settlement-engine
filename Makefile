BIN?=bet-settlement-engine

default: run
.PHONY : build run

lint:
	golangci-lint run -c .golangci.yml

build:
	go build -o build/${BIN}

run: build
	./build/${BIN}

populate-cache:
	go run cmd/main.go
