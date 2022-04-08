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

package fuse2grpc

import (
	"context"
	"fmt"
	"unsafe"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func (s *server) String(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	grpc_logrus.Extract(ctx).Debug("String")
	return &pb.StringResponse{Value: s.fs.String()}, nil
}

func (s *server) Lookup(ctx context.Context, req *pb.LookupRequest) (*pb.LookupResponse, error) {
	var (
		out    fuse.EntryOut
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"parent": req.Header.NodeId,
		"name":   req.Name,
	}).Debug("Lookup")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Lookup(ch, &header, req.Name, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
	}
	if st != fuse.OK {
		return &pb.LookupResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.LookupResponse{
		EntryOut: &pb.EntryOut{
			NodeId:         out.NodeId,
			Generation:     out.Generation,
			Attr:           toPbAttr(&out.Attr),
			AttrValid:      out.AttrValid,
			AttrValidNsec:  out.AttrValidNsec,
			EntryValid:     out.EntryValid,
			EntryValidNsec: out.EntryValidNsec,
		},
		Status: &pb.Status{Code: 0},
	}, nil
}
func (s *server) Forget(ctx context.Context, req *pb.ForgetRequest) (*emptypb.Empty, error) {
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeid":  req.Nodeid,
		"nlookup": req.Nlookup,
	}).Debug("Forget")
	s.fs.Forget(req.Nodeid, req.Nlookup)
	return &emptypb.Empty{}, nil
}

func (s *server) GetAttr(ctx context.Context, req *pb.GetAttrRequest) (*pb.GetAttrResponse, error) {
	var (
		out    fuse.AttrOut
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
	}).Debug("GetAttr")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.GetAttr(ch, &fuse.GetAttrIn{InHeader: header}, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method GetAttr not implemented")
	}
	if st != fuse.OK {
		return &pb.GetAttrResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.GetAttrResponse{
		AttrOut: &pb.AttrOut{
			Attr:          toPbAttr(&out.Attr),
			AttrValid:     out.AttrValid,
			AttrValidNsec: out.AttrValidNsec,
		},
		Status: &pb.Status{Code: 0},
	}, nil
}

func (s *server) SetAttr(ctx context.Context, req *pb.SetAttrRequest) (*pb.SetAttrResponse, error) {
	var (
		out    fuse.AttrOut
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
	}).Debug("SetAttr")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.SetAttr(ch,
		&fuse.SetAttrIn{
			SetAttrInCommon: fuse.SetAttrInCommon{
				InHeader:  header,
				Valid:     req.Valid,
				Padding:   req.Padding,
				Fh:        req.Fh,
				Size:      req.Size,
				LockOwner: req.LockOwner,
				Atime:     req.Atime,
				Mtime:     req.Mtime,
				Ctime:     req.Ctime,
				Atimensec: req.Atimensec,
				Mtimensec: req.Mtimensec,
				Ctimensec: req.Ctimensec,
				Mode:      req.Mode,
				Unused4:   req.Unused4,
				Owner: fuse.Owner{
					Uid: req.Owner.Uid,
					Gid: req.Owner.Gid,
				},
				Unused5: req.Unused5,
			},
		},
		&out,
	)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method SetAttr not implemented")
	}
	if st != fuse.OK {
		return &pb.SetAttrResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.SetAttrResponse{
		AttrOut: &pb.AttrOut{
			Attr: toPbAttr(&out.Attr),
		},
		Status: &pb.Status{Code: 0},
	}, nil
}

func (s *server) Access(ctx context.Context, req *pb.AccessRequest) (*pb.AccessResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
	}).Debug("Access")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Access(ch, &fuse.AccessIn{InHeader: header, Mask: req.Mask, Padding: req.Padding})
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Access not implemented")
	}
	return &pb.AccessResponse{Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) GetXAttr(context.Context, *pb.GetXAttrRequest) (*pb.GetXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetXAttr not implemented")
}
func (s *server) ListXAttr(context.Context, *pb.ListXAttrRequest) (*pb.ListXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListXAttr not implemented")
}
func (s *server) SetXAttr(context.Context, *pb.SetXAttrRequest) (*pb.SetXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetXAttr not implemented")
}
func (s *server) RemoveXAttr(context.Context, *pb.RemoveXAttrRequest) (*pb.RemoveXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveXAttr not implemented")
}
func (s *server) Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func (s *server) Open(ctx context.Context, req *pb.OpenRequest) (*pb.OpenResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.OpenOut
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.OpenIn.Header.NodeId,
		"flags":  req.OpenIn.Flags,
		"mode":   req.OpenIn.Mode,
	}).Debug("Open")
	toFuseInHeader(req.OpenIn.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Open(ch, &fuse.OpenIn{InHeader: header, Flags: req.OpenIn.Flags, Mode: req.OpenIn.Mode}, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Open not implemented")
	}
	if st != fuse.OK {
		return &pb.OpenResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.OpenResponse{
		OpenOut: &pb.OpenOut{
			Fh:        out.Fh,
			OpenFlags: out.OpenFlags,
			Padding:   out.Padding,
		},
		Status: &pb.Status{Code: 0},
	}, nil
}

func (s *server) Read(req *pb.ReadRequest, stream pb.RawFileSystem_ReadServer) error {
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (s *server) LSeek(context.Context, *pb.LSeekRequest) (*pb.LSeekResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LSeek not implemented")
}
func (s *server) GetLk(context.Context, *pb.LkRequest) (*pb.GetLkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLk not implemented")
}
func (s *server) SetLk(context.Context, *pb.LkRequest) (*pb.SetLkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLk not implemented")
}
func (s *server) SetLkw(context.Context, *pb.LkRequest) (*pb.SetLkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLkw not implemented")
}

func (s *server) Release(ctx context.Context, req *pb.ReleaseRequest) (*emptypb.Empty, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":       req.Header.NodeId,
		"fh":           req.Fh,
		"flags":        req.Flags,
		"releaseFlags": req.ReleaseFlags,
		"lockOwner":    req.LockOwner,
	}).Debug("Release")
	ch := newCancel(ctx)
	defer releaseCancel(ch)

	toFuseInHeader(req.Header, &header)
	s.fs.Release(ch, &fuse.ReleaseIn{InHeader: header, Fh: req.Fh, Flags: req.Flags, ReleaseFlags: req.ReleaseFlags, LockOwner: req.LockOwner})
	return &emptypb.Empty{}, nil
}

func (s *server) Write(context.Context, *pb.WriteRequest) (*pb.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (s *server) CopyFileRange(context.Context, *pb.CopyFileRangeRequest) (*pb.CopyFileRangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CopyFileRange not implemented")
}

func (s *server) Flush(ctx context.Context, req *pb.FlushRequest) (*pb.FlushResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":    req.Header.NodeId,
		"fh":        req.Fh,
		"unused":    req.Unused,
		"padding":   req.Padding,
		"lockOwner": req.LockOwner,
	}).Debug("OpenDir")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Flush(ch, &fuse.FlushIn{InHeader: header, Fh: req.Fh, Unused: req.Unused, Padding: req.Padding, LockOwner: req.LockOwner})
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
	}
	return &pb.FlushResponse{Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) Fsync(context.Context, *pb.FsyncRequest) (*pb.FsyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fsync not implemented")
}
func (s *server) Fallocate(context.Context, *pb.FallocateRequest) (*pb.FallocateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fallocate not implemented")
}

func (s *server) OpenDir(ctx context.Context, req *pb.OpenDirRequest) (*pb.OpenDirResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.OpenOut
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.OpenIn.Header.NodeId,
		"flags":  req.OpenIn.Flags,
		"mode":   req.OpenIn.Mode,
	}).Debug("OpenDir")
	toFuseInHeader(req.OpenIn.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.OpenDir(ch, &fuse.OpenIn{InHeader: header, Flags: req.OpenIn.Flags, Mode: req.OpenIn.Mode}, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method OpenDir not implemented")
	}
	if st != fuse.OK {
		return &pb.OpenDirResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.OpenDirResponse{
		OpenOut: &pb.OpenOut{
			Fh:        out.Fh,
			OpenFlags: out.OpenFlags,
			Padding:   out.Padding,
		},
		Status: &pb.Status{Code: 0},
	}, nil
}

func (s *server) doReadDir(
	req *pb.ReadDirRequest,
	stream pb.RawFileSystem_ReadDirServer,
	reader func(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status,
	readerName string,
	prefix uint32,
) error {
	var (
		header           fuse.InHeader
		batchSize, delta int
		pos              uint32
		batch            []*pb.DirEntry
	)
	ctx := stream.Context()
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":    req.ReadIn.Header.NodeId,
		"fh":        req.ReadIn.Fh,
		"offset":    req.ReadIn.Offset,
		"size":      req.ReadIn.Size,
		"readFlags": req.ReadIn.ReadFlags,
	}).Debug(readerName)
	toFuseInHeader(req.ReadIn.Header, &header)

	buf := s.buffers.AllocBuffer(req.ReadIn.Size)
	defer s.buffers.FreeBuffer(buf)

	out := fuse.NewDirEntryList(buf, req.ReadIn.Offset)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := reader(ch,
		&fuse.ReadIn{
			InHeader:  header,
			Fh:        req.ReadIn.Fh,
			Offset:    req.ReadIn.Offset,
			Size:      req.ReadIn.Size,
			ReadFlags: req.ReadIn.ReadFlags,
		}, out)

	if st == fuse.ENOSYS {
		return status.Errorf(codes.Unimplemented, fmt.Sprintf("method %s not implemented", readerName))
	}

	if st != fuse.OK {
		stream.Send(&pb.ReadDirResponse{Status: &pb.Status{Code: int32(st)}})
		return nil
	}

	flushFunc := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := stream.Send(&pb.ReadDirResponse{
			Entries: batch,
			Status:  &pb.Status{Code: 0},
		}); err != nil {
			return err
		}
		return nil
	}

	buf = (*DirEntryList)(unsafe.Pointer(out)).buf
	bufsize := len(buf)

	for {
		if int(pos) >= bufsize || len(buf[pos:]) < int(direntSize) {
			break
		}

		pos += prefix

		e := (*_Dirent)(unsafe.Pointer(&buf[pos]))
		if e.Off == 0 {
			break
		}
		// uint64 Ino uint64 Offset uint32 NameLen uint32 Typ
		delta = deltaSize(e)
		if batchSize+delta > s.msgSizeThreshold {
			if err := flushFunc(); err != nil {
				return err
			}
			batch = nil
			batchSize = 0
		}
		dirEntry := &pb.DirEntry{Mode: typeToMode(e.Typ), Ino: e.Ino, Name: buf[pos+direntSize : pos+direntSize+e.NameLen]}
		batch = append(batch, dirEntry)
		batchSize += delta
		padding := (8 - e.NameLen&7)
		pos += padding + direntSize + e.NameLen
	}
	if len(batch) == 0 {
		return nil
	}
	return flushFunc()
}

func (s *server) ReadDir(req *pb.ReadDirRequest, stream pb.RawFileSystem_ReadDirServer) error {
	return s.doReadDir(req, stream, s.fs.ReadDir, "ReadDir", 0)
}

func deltaSize(e *_Dirent) int {
	// uint32 Mode, uint64 Ino, string name
	return int(4 + 8 + e.NameLen)
}

func (s *server) ReadDirPlus(req *pb.ReadDirRequest, stream pb.RawFileSystem_ReadDirPlusServer) error {
	return s.doReadDir(req, stream, s.fs.ReadDirPlus, "ReadDirPlus", entryOutSize)
}

func (s *server) ReleaseDir(ctx context.Context, req *pb.ReleaseRequest) (*emptypb.Empty, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":       req.Header.NodeId,
		"fh":           req.Fh,
		"flags":        req.Flags,
		"releaseFlags": req.ReleaseFlags,
		"lockOwner":    req.LockOwner,
	}).Debug("ReleaseDir")
	toFuseInHeader(req.Header, &header)
	s.fs.ReleaseDir(&fuse.ReleaseIn{InHeader: header, Fh: req.Fh, Flags: req.Flags, ReleaseFlags: req.ReleaseFlags, LockOwner: req.LockOwner})
	return &emptypb.Empty{}, nil
}

func (s *server) FsyncDir(context.Context, *pb.FsyncRequest) (*pb.FsyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FsyncDir not implemented")
}

func (s *server) StatFs(ctx context.Context, req *pb.StatfsRequest) (*pb.StatfsResponse, error) {
	var (
		out    fuse.StatfsOut
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Input.NodeId,
	}).Debug("StatFs")
	toFuseInHeader(req.Input, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)
	st := s.fs.StatFs(ch, &header, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method StatFS not implemented")
	}
	if st != fuse.OK {
		return &pb.StatfsResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.StatfsResponse{
		Blocks:  out.Blocks,
		Bfree:   out.Bfree,
		Bavail:  out.Bavail,
		Files:   out.Files,
		Ffree:   out.Ffree,
		Bsize:   out.Bsize,
		NameLen: out.NameLen,
		Frsize:  out.Frsize,
		Status:  &pb.Status{Code: 0},
	}, nil
}
