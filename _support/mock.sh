#!/bin/sh
set -ex

mockgen -source=pb/raw_file_system.pb.go -destination=mock/client_mock.go -package=mock RawFileSystemClient

if hash gofumpt 2>/dev/null; then
  gofumpt -w mock/
else
  gofmt -w mock/
fi
