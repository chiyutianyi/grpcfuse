#!/bin/sh

protoc -I vendor -I proto \
        proto/shared.proto \
        proto/raw_file_system.proto \
        --go_opt=paths=source_relative \
        --go_out=pb \
        --go-grpc_opt=paths=source_relative \
        --go-grpc_out=pb
