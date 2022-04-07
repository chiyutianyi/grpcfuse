package grpc2fuse

import (
	"github.com/hanwen/go-fuse/v2/fuse"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func toPbHeader(header *fuse.InHeader) *pb.InHeader {
	return &pb.InHeader{
		Length: header.Length,
		Opcode: header.Opcode,
		Unique: header.Unique,
		NodeId: header.NodeId,
		Caller: &pb.Caller{
			Owner: &pb.Owner{
				Uid: header.Uid,
				Gid: header.Gid,
			},
			Pid: header.Pid,
		},
	}
}

func toFuseAttr(out *fuse.Attr, in *pb.Attr) {
	out.Ino = in.Ino
	out.Size = in.Size
	out.Blocks = in.Blocks

	out.Atime = in.Atime
	out.Mtime = in.Mtime
	out.Ctime = in.Ctime
	out.Atimensec = in.Atimensec
	out.Mtimensec = in.Mtimensec
	out.Ctimensec = in.Ctimensec

	out.Mode = in.Mode
	out.Nlink = in.Nlink

	out.Uid = in.Owner.Uid
	out.Gid = in.Owner.Gid

	out.Rdev = in.Rdev
	setFlags(out, in.Flags)
	setBlksize(out, in.Blksize)
	setPadding(out, in.Padding)
}

func toFuseEntryOut(out *fuse.EntryOut, in *pb.EntryOut) {
	out.NodeId = in.NodeId
	out.Generation = in.Generation
	out.AttrValid = in.AttrValid
	out.AttrValidNsec = in.AttrValidNsec
	out.EntryValid = in.EntryValid
	out.EntryValidNsec = in.EntryValidNsec
	toFuseAttr(&out.Attr, in.Attr)
}

func toFuseAttrOut(out *fuse.AttrOut, in *pb.AttrOut) {
	out.AttrValid = in.AttrValid
	out.AttrValidNsec = uint32(in.AttrValidNsec)
	toFuseAttr(&out.Attr, in.Attr)
}
