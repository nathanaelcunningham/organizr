.PHONY: build run clean test test-coverage test-race tidy

build:
	CGO_ENABLED=1 go build -o bin/organizr ./cmd/api

run: build
	./bin/organizr

clean:
	rm -rf bin/

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-race:
	go test -race ./...

tidy:
	go mod tidy
