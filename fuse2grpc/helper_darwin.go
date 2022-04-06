package fuse2grpc

import (
	"github.com/chiyutianyi/grpcfuse/pb"

	"github.com/hanwen/go-fuse/v2/fuse"
)

func toPbAttr(in *fuse.Attr) *pb.Attr {
	return &pb.Attr{
		Ino:       in.Ino,
		Size:      in.Size,
		Blocks:    in.Blocks,
		Atime:     in.Atime,
		Mtime:     in.Mtime,
		Ctime:     in.Ctime,
		Atimensec: in.Atimensec,
		Mtimensec: in.Mtimensec,
		Ctimensec: in.Ctimensec,
		Mode:      in.Mode,
		Nlink:     in.Nlink,
		Owner: &pb.Owner{
			Uid: in.Uid,
			Gid: in.Gid,
		},
		Rdev:  in.Rdev,
		Flags: in.Flags_,
	}
}
