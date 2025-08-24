SHELL = /bin/sh

PKG := github.com/chiyutianyi/grpcfuse

# Host information
OS := $(shell uname)
ifeq (${OS},Darwin)
    GO_GOOS ?= darwin
else ifeq (${OS},Linux)
    GO_GOOS ?= linux
else
    $(error Unsupported OS: ${OS})
endif

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

.PHONY: all test testfuse2grpc testgrpc2fuse mock example client loopback clean lint vet fmt coverage benchmark proto

all: proto example test

test: testfuse2grpc testgrpc2fuse

testfuse2grpc:
	@echo "Testing fuse2grpc package..."
	go test -v -cover -coverprofile=coverage_fuse2grpc.out ${PKG}/fuse2grpc...
	go tool cover -func=coverage_fuse2grpc.out | grep statements

testgrpc2fuse:
	@echo "Testing grpc2fuse package..."
	go test -v -cover -coverprofile=coverage_grpc2fuse.out ${PKG}/grpc2fuse...
	go tool cover -func=coverage_grpc2fuse.out | grep statements

test-all: test
	@echo "Running all tests with race detection..."
	go test -race -v ${PKG}/...

coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage_fuse2grpc.out -o coverage_fuse2grpc.html
	go tool cover -html=coverage_grpc2fuse.out -o coverage_grpc2fuse.html
	@echo "Coverage reports generated: coverage_fuse2grpc.html, coverage_grpc2fuse.html"

benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ${PKG}/fuse2grpc...
	go test -bench=. -benchmem ${PKG}/grpc2fuse...

lint:
	@echo "Running linter..."
	golangci-lint run

vet:
	@echo "Running go vet..."
	go vet ${PKG}/...

fmt:
	@echo "Formatting code..."
	go fmt ${PKG}/...
	go fmt example/...

fmt-check:
	@echo "Checking code formatting..."
	@if [ "$$(gofmt -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		gofmt -l .; \
		exit 1; \
	fi

proto:
	@echo "Generating protobuf code..."
	@if command -v protoc >/dev/null 2>&1; then \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/*.proto; \
	else \
		echo "protoc not found. Skipping protobuf generation."; \
	fi

mock:
	@echo "Generating mocks..."
	_support/mock.sh

example: client loopback

client:
	@echo "Building client example..."
	GOOS=${GO_GOOS} GOARCH=amd64 go build ${LDFLAGS} -o bin/client example/client/client.go

loopback:
	@echo "Building loopback example..."
	GOOS=${GO_GOOS} GOARCH=amd64 go build ${LDFLAGS} -o bin/loopback example/loopback/server.go

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/*
	rm -f coverage*.out
	rm -f coverage*.html

# Development helpers
dev-setup: install proto mock
	@echo "Development environment setup complete"

pre-commit: fmt-check lint vet test
	@echo "Pre-commit checks passed"

# Docker support
docker-build:
	@echo "Building Docker image..."
	docker build -t grpcfuse:latest .

docker-test:
	@echo "Running tests in Docker..."
	docker run --rm grpcfuse:latest make test

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build everything and run tests"
	@echo "  test         - Run all tests"
	@echo "  coverage     - Generate coverage reports"
	@echo "  benchmark    - Run benchmarks"
	@echo "  lint         - Run linter"
	@echo "  vet          - Run go vet"
	@echo "  fmt          - Format code"
	@echo "  fmt-check    - Check code formatting"
	@echo "  proto        - Generate protobuf code"
	@echo "  mock         - Generate mocks"
	@echo "  example      - Build examples"
	@echo "  install      - Install dependencies"
	@echo "  clean        - Clean build artifacts"
	@echo "  dev-setup    - Setup development environment"
	@echo "  pre-commit   - Run pre-commit checks"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-test  - Run tests in Docker"
	@echo "  help         - Show this help"