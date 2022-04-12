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

test: testfuse2grpc testgrpc2fuse

testfuse2grpc:
	go test -cover -coverprofile=coverage.out ${PKG}/fuse2grpc...
	go tool cover -func=coverage.out | grep statements

testgrpc2fuse:
	go test -cover -coverprofile=coverage.out ${PKG}/grpc2fuse...
	go tool cover -func=coverage.out | grep statements

mock:
	_support/mock.sh

example: client loopback

client:
	GOOS=${GO_GOOS} GOARCH=amd64 go build -o bin/client example/client/client.go

loopback:
	GOOS=${GO_GOOS} GOARCH=amd64 go build -o bin/loopback example/loopback/server.go

clean:
	rm -f bin/*