package rawfilesystem

import (
	"github.com/hanwen/go-fuse/v2/fuse"

	"github.com/chiyutianyi/grpcfuse/pb"
)

type server struct {
	pb.UnimplementedRawFileSystemServer

	fs fuse.RawFileSystem
}

// NewServer returns a new loopback server.
func NewServer(fs fuse.RawFileSystem) *server {
	srv := &server{fs: fs}
	return srv
}
