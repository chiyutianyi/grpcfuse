package rawfilesystem

import (
	"github.com/hanwen/go-fuse/v2/fuse"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func toFuseInHeader(in *pb.InHeader) *fuse.InHeader {
	return &fuse.InHeader{
		Length: in.Length,
		Opcode: in.Opcode,
		Unique: in.Unique,
		NodeId: in.NodeId,
		Caller: fuse.Caller{
			Owner: fuse.Owner{
				Uid: in.Caller.Owner.Uid,
				Gid: in.Caller.Owner.Gid,
			},
			Pid: in.Caller.Pid,
		},
		Padding: in.Padding,
	}
}
