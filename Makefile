
check-gow:
	@command -v gow >/dev/null 2>&1 || { echo >&2 "gow is required but it's not installed. Aborting."; exit 1; }

dev-unit: check-gow
	gow test -p 1 -failfast -run Unit_ ./...

dev-integration: check-gow
	gow test -p 1 -failfast -run Integration_ ./...

dev: check-gow
	gow test -p 1 -failfast ./...

test: deps
	go test -p 1 -failfast ./...

test-unit: deps
	go test -p 1 -failfast -run Unit_ ./...

test-integration: deps
	go test -p 1 -failfast -run Integration_ ./...

deps:
	go mod tidy && go mod vendor

VERSION := $(shell git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null || echo "v0.0.0")

install: deps
	cd cmd/gho && go install -ldflags="-X 'main.Version=${VERSION}'"

lint: deps
	golangci-lint run