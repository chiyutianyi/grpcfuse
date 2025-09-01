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

all: example test

test: testfuse2grpc testgrpc2fuse testutils testintegration

testfuse2grpc:
	go test -cover -coverprofile=fuse2grpc_coverage.out ${PKG}/fuse2grpc...
	go tool cover -func=fuse2grpc_coverage.out | grep statements

testgrpc2fuse:
	go test -cover -coverprofile=grpc2fuse_coverage.out ${PKG}/grpc2fuse...
	go tool cover -func=grpc2fuse_coverage.out | grep statements

testutils:
	go test -cover -coverprofile=utils_coverage.out ${PKG}/pkg/utils...
	go tool cover -func=utils_coverage.out | grep statements

testintegration:
	go test -cover -coverprofile=integration_coverage.out ${PKG}...
	go tool cover -func=integration_coverage.out | grep statements

testall:
	go test -cover -coverprofile=all_coverage.out ./...
	go tool cover -func=all_coverage.out

bench:
	go test -bench=. -benchmem ./...

perftest:
	go test -run=TestPerformance -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

mock:
	_support/mock.sh

example: client loopback

client:
	GOOS=${GO_GOOS} GOARCH=amd64 go build -o bin/client example/client/client.go

loopback:
	GOOS=${GO_GOOS} GOARCH=amd64 go build -o bin/loopback example/loopback/server.go

clean:
	rm -f bin/*