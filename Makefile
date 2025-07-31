export GO111MODULE=on

# Build variables
BINARY_NAME=argo-apps-viz
BINARY_PATH=$(HOME)/go/bin/$(BINARY_NAME)
GO_MODULE=argo-apps-viz
GO_FILES=$(shell find . -type f -name '*.go')

# Version information
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all
all: clean test build

.PHONY: build
build: fmt vet
	go build $(LDFLAGS) -o $(BINARY_PATH) $(GO_MODULE)/cmd/plugin

.PHONY: test
test:
	go test -v -race -cover ./pkg/... ./cmd/... -coverprofile cover.out
	go tool cover -html=cover.out -o coverage.html

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	@echo "Running go vet..."
	@go vet ./pkg/... ./cmd/... || (echo "Go vet failed. Please fix the issues above."; exit 1)

.PHONY: lint
lint:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		brew install golangci-lint; \
	fi
	golangci-lint run

.PHONY: clean
clean:
	rm -rf bin/ cover.out coverage.html


.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all       - Clean, run tests, and build"
	@echo "  build     - Build the binary"
	@echo "  test      - Run tests with coverage"
	@echo "  fmt       - Format code"
	@echo "  vet       - Run go vet"
	@echo "  lint      - Run golangci-lint"
	@echo "  clean     - Clean build artifacts"
	@echo "  dev       - Run directly without building (use ARGS)"
	@echo "  help      - Show this help message"

.PHONY: dev
dev:
	@go run cmd/plugin/main.go $(ARGS)