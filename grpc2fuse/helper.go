/*
 * Copyright 2022 Han Xin, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpc2fuse

import (
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func toFuseOpenOut(out *fuse.OpenOut, in *pb.OpenOut) {
	out.Fh = in.Fh
	out.OpenFlags = in.OpenFlags
	out.Padding = in.Padding
}

func dealGrpcError(method string, err error) fuse.Status {
	if err == nil {
		return fuse.OK
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() == codes.Unimplemented {
			log.Warnf("%s unimplemented", method)
			return fuse.ENOSYS
		}
	}
	log.Errorf("%s: %v", method, err)
	return fuse.EIO
}
