package fuse2grpc

import (
	"github.com/hanwen/go-fuse/v2/fuse"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func toFuseInHeader(in *pb.InHeader, out *fuse.InHeader) {
	out.Length = in.Length
	out.Opcode = in.Opcode
	out.Unique = in.Unique
	out.NodeId = in.NodeId
	out.Uid = in.Caller.Owner.Uid
	out.Gid = in.Caller.Owner.Gid
	out.Pid = in.Caller.Pid
	out.Padding = in.Padding
}
