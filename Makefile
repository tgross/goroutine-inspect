MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := build

GO_SRC := $(wildcard ./*.go)

.PHONY: build
build: dist/goroutine-inspect

dist/goroutine-inspect: $(GO_SRC) protos/profile.pb.go
	@mkdir -p ./dist
	go build -trimpath -o dist/goroutine-inspect .

protos/profile.pb.go: protos/profile.proto
	protoc --proto_path=./protos \
		--go_out=./protos \
		--go_opt=Mprofile.proto=github.com/tgross/goroutine-inspect/protos \
		--go_opt=paths=source_relative \
		./protos/profile.proto

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
