.PHONY: build run clean test

build:
	CGO_ENABLED=1 go build -o bin/organizr ./cmd/api

run: build
	./bin/organizr

clean:
	rm -rf bin/

test:
	go test -v ./...

tidy:
	go mod tidy
