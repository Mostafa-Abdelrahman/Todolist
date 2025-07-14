.PHONY: build run test clean

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test ./...

clean:
	rm -rf bin/
	rm -f todos.db

install:
	go mod download

deps:
	go mod tidy