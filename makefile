src := $(shell find . -type f -name '*.go')

OUT_DIR := ./dist

build: $(src)
	go build -o $(OUT_DIR)/gbac-opa cmd/gbac-opa/main.go

PHONY: run-opa
run-opa: build
	go run cmd/gbac-opa/main.go run --server --config-file test/redis/config.yaml test/rego/