.DEFAULT_GOAL = build
extension = $(patsubst windows,.exe,$(filter windows,$(1)))
GO := go
PKG_NAME := hitron
PREFIX := .

COMMIT ?= `git rev-parse --short HEAD 2>/dev/null`
VERSION ?= `git describe --abbrev=0 --tags $(git rev-list --tags --max-count=1) 2>/dev/null | sed 's/v\(.*\)/\1/'`

COMMIT_FLAG := -X `go list ./internal/version`.GitCommit=$(COMMIT)
VERSION_FLAG := -X `go list ./internal/version`.Version=$(VERSION)

GOOS ?= $(shell go version | sed 's/^.*\ \([a-z0-9]*\)\/\([a-z0-9]*\)/\1/')
GOARCH ?= $(shell go version | sed 's/^.*\ \([a-z0-9]*\)\/\([a-z0-9]*\)/\2/')

ifeq ("$(TARGETVARIANT)","")
ifneq ("$(GOARM)","")
TARGETVARIANT := v$(GOARM)
endif
else
ifeq ("$(GOARM)","")
GOARM ?= $(subst v,,$(TARGETVARIANT))
endif
endif

clean:
	rm -Rf $(PREFIX)/bin/*

$(PREFIX)/bin/$(PKG_NAME)_%v5$(call extension,$(GOOS)): $(shell find $(PREFIX) -type f -name "*.go")
	GOOS=$(shell echo $* | cut -f1 -d-) GOARCH=$(shell echo $* | cut -f2 -d- ) GOARM=5 CGO_ENABLED=0 \
		$(GO) build \
			-ldflags "-w -s $(COMMIT_FLAG) $(VERSION_FLAG)" \
			-o $@ \
			./cmd/$(PKG_NAME)

$(PREFIX)/bin/$(PKG_NAME)_%v6$(call extension,$(GOOS)): $(shell find $(PREFIX) -type f -name "*.go")
	GOOS=$(shell echo $* | cut -f1 -d-) GOARCH=$(shell echo $* | cut -f2 -d- ) GOARM=6 CGO_ENABLED=0 \
		$(GO) build \
			-ldflags "-w -s $(COMMIT_FLAG) $(VERSION_FLAG)" \
			-o $@ \
			./cmd/$(PKG_NAME)

$(PREFIX)/bin/$(PKG_NAME)_%v7$(call extension,$(GOOS)): $(shell find $(PREFIX) -type f -name "*.go")
	GOOS=$(shell echo $* | cut -f1 -d-) GOARCH=$(shell echo $* | cut -f2 -d- ) GOARM=7 CGO_ENABLED=0 \
		$(GO) build \
			-ldflags "-w -s $(COMMIT_FLAG) $(VERSION_FLAG)" \
			-o $@ \
			./cmd/$(PKG_NAME)

$(PREFIX)/bin/$(PKG_NAME)_windows-%.exe: $(shell find $(PREFIX) -type f -name "*.go")
	GOOS=windows GOARCH=$* GOARM= CGO_ENABLED=0 \
		$(GO) build \
			-ldflags "-w -s $(COMMIT_FLAG) $(VERSION_FLAG)" \
			-o $@ \
			./cmd/$(PKG_NAME)

$(PREFIX)/bin/$(PKG_NAME)_%$(TARGETVARIANT)$(call extension,$(GOOS)): $(shell find $(PREFIX) -type f -name "*.go")
	GOOS=$(shell echo $* | cut -f1 -d-) GOARCH=$(shell echo $* | cut -f2 -d- ) GOARM=$(GOARM) CGO_ENABLED=0 \
		$(GO) build \
			-ldflags "-w -s $(COMMIT_FLAG) $(VERSION_FLAG)" \
			-o $@ \
			./cmd/$(PKG_NAME)

$(PREFIX)/bin/$(PKG_NAME)$(call extension,$(GOOS)): $(PREFIX)/bin/$(PKG_NAME)_$(GOOS)-$(GOARCH)$(TARGETVARIANT)$(call extension,$(GOOS))
	cp $< $@

build: $(PREFIX)/bin/$(PKG_NAME)_$(GOOS)-$(GOARCH)$(TARGETVARIANT)$(call extension,$(GOOS)) $(PREFIX)/bin/$(PKG_NAME)$(call extension,$(GOOS))

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

.PHONY: clean test lint ci-lint
.DELETE_ON_ERROR:
.SECONDARY:
