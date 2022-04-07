package fuse2grpc

import (
	"github.com/hanwen/go-fuse/v2/fuse"

	"github.com/chiyutianyi/grpcfuse/pb"
)

// msgSizeThreshold 1mb < default grpc message size limit 4mb
const msgSizeThreshold = 1 << 20

type server struct {
	pb.UnimplementedRawFileSystemServer

	fs fuse.RawFileSystem

	buffers bufferPool

	msgSizeThreshold int
}

// NewServer returns a new loopback server.
func NewServer(fs fuse.RawFileSystem) *server {
	return &server{fs: fs, buffers: bufferPool{}, msgSizeThreshold: msgSizeThreshold}
}
