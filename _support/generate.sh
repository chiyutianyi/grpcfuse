#/bin/sh

protoc -I proto -I $GOPATH/src shared.proto raw_file_system.proto  --go_out=plugins=grpc,paths=source_relative:pb
