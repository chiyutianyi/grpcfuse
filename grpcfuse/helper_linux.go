package grpcfuse

import (
	"github.com/hanwen/go-fuse/v2/fuse"
)

func getUmask(in *fuse.MknodIn) uint16 {
	return uint16(in.Umask)
}

func setFlags(out *fuse.Attr, flags uint32) {
}

func setBlksize(out *fuse.Attr, size uint32) {
	out.Blksize = size
}

func setPadding(out *fuse.Attr, padding uint32) {
	out.Padding = padding
}
