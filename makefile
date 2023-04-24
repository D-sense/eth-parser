SHELL := /bin/bash
.PHONY: api-run clean docker-build test tidy fmt

# ==============================================================================
# Building without container
build-run:
	go build -o api/bin/eth_parser api/main.go
	./api/bin/eth_parser

# ==============================================================================
# Building container

VERSION := 1.0

docker-build:
	docker build \
		-f api/Dockerfile \
		-t ethereum-parser-api \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

docker-run:
	docker run -d -p 8080:8080 ethereum-parser-api

# ==============================================================================
# Testing coverage support
test:
	go test ./... -v -count=1

# ==============================================================================
# Modules support

tidy:
	go mod tidy


# ==============================================================================
# Code Format support
fmt:
	go fmt ./...

