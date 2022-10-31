.PHONY: build

default: build

lint:
	golangci-lint run --timeout=5m

tests:
	./scripts/tests.sh

build: lint tests

tidy:
	go mod tidy -v
