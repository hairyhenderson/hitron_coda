.DEFAULT_GOAL = build
extension = $(patsubst windows,.exe,$(filter windows,$(1)))
GO := go
PREFIX := .

COMMIT ?= `git rev-parse --short HEAD 2>/dev/null`
VERSION ?= `git describe --abbrev=0 --tags $(git rev-list --tags --max-count=1) 2>/dev/null | sed 's/v\(.*\)/\1/'`

COMMIT_FLAG := -X `go list ./version`.GitCommit=$(COMMIT)
VERSION_FLAG := -X `go list ./version`.Version=$(VERSION)

GOOS ?= $(shell go version | sed 's/^.*\ \([a-z0-9]*\)\/\([a-z0-9]*\)/\1/')
GOARCH ?= $(shell go version | sed 's/^.*\ \([a-z0-9]*\)\/\([a-z0-9]*\)/\2/')

ifeq ($(OS),Windows_NT)
test:
	$(GO) test -coverprofile=c.out ./...
else
test:
	$(GO) test -race -coverprofile=c.out ./...
endif

lint:
	@golangci-lint run --verbose --max-same-issues=0 --max-issues-per-linter=0 --sort-results

ci-lint:
	@golangci-lint run --verbose --max-same-issues=0 --max-issues-per-linter=0 --sort-results --out-format=github-actions

.PHONY: clean test lint
.DELETE_ON_ERROR:
.SECONDARY:
