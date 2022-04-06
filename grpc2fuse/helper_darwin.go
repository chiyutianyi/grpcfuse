package grpc2fuse

import (
	"github.com/hanwen/go-fuse/v2/fuse"
)

func getUmask(in *fuse.MknodIn) uint16 {
	return 0
}

func setFlags(out *fuse.Attr, flags uint32) {
	out.Flags_ = flags
}

func setBlksize(out *fuse.Attr, size uint32) {
}

func setPadding(out *fuse.Attr, padding uint32) {
}
