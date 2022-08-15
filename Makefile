.PHONY: build

build:
	go build -v ./cmd/api/main.go

.DEFAULT_GOAL := build