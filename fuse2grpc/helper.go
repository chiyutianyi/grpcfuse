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

// modeToType converts a file *type* (as used in _Dirent.Typ)
// to a file *mode* (as used in syscall.Stat_t.Mode).
func typeToMode(typ uint32) uint32 {
	return (typ << 12) & 0170000
}
