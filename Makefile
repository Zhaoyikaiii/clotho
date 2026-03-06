.PHONY: all install clean help

# Build variables
BINARY_NAME=clotho
BUILD_DIR=build
CMD_DIR=cmd/$(BINARY_NAME)

# Version
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short=8 HEAD 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%FT%T%z)
GO_VERSION=$(shell go version | awk '{print $$3}')
INTERNAL=github.com/Zhaoyikaiii/clotho/cmd/clotho/internal
LDFLAGS=-ldflags '-X $(INTERNAL).version=$(VERSION) -X $(INTERNAL).gitCommit=$(GIT_COMMIT) -X $(INTERNAL).buildTime=$(BUILD_TIME) -X $(INTERNAL).goVersion=$(GO_VERSION) -s -w'

# Installation
INSTALL_PREFIX?=$(HOME)/.local
INSTALL_BIN_DIR=$(INSTALL_PREFIX)/bin

# OS detection
UNAME_S:=$(shell uname -s)
UNAME_M:=$(shell uname -m)

ifeq ($(UNAME_S),Linux)
	PLATFORM=linux
	ifeq ($(UNAME_M),x86_64)
		ARCH=amd64
	else ifeq ($(UNAME_M),aarch64)
		ARCH=arm64
	else
		ARCH=$(UNAME_M)
	endif
else ifeq ($(UNAME_S),Darwin)
	PLATFORM=darwin
	ifeq ($(UNAME_M),x86_64)
		ARCH=amd64
	else
		ARCH=arm64
	endif
endif

all: install

## install: Install dev dependencies (go mod, lefthook) and build binary
install:
	@echo "Installing dev dependencies..."
	@go mod download
	@if ! command -v lefthook &> /dev/null; then \
		echo "Installing lefthook..."; \
		brew install lefthook; \
	fi
	@lefthook install -f lefthook.yml
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@mkdir -p $(INSTALL_BIN_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_BIN_DIR)/$(BINARY_NAME)
	@echo "Installed to $(INSTALL_BIN_DIR)/$(BINARY_NAME)"

## build: Build binary without installing
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

## clean: Remove build artifacts
clean:
	@rm -rf $(BUILD_DIR)

## help: Show this help message
help:
	@echo "Clotho Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sort | awk -F': ' '{printf "  %-16s %s\n", substr($$1, 4), $$2}'
