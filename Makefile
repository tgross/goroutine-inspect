MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := build

GO_SRC := $(wildcard ./*.go)

.PHONY: build
build: dist/goroutine-inspect

dist/goroutine-inspect: $(GO_SRC)
	@mkdir -p ./dist
	go build -trimpath -o dist/goroutine-inspect .

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: check
check:
	go vet ./...
	go mod tidy

.PHONY: clean
clean:
	rm -rf ./dist
