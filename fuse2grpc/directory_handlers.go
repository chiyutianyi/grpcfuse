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
		// fuse.DirEntryList.Add()
		// padding := (8 - len(name)&7) & 7
		padding := (8 - e.NameLen&7) & 7
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

func (s *server) FsyncDir(ctx context.Context, req *pb.FsyncRequest) (*pb.FsyncResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":     req.Header.NodeId,
		"fh":         req.Fh,
		"fsyncFlags": req.FsyncFlags,
		"padding":    req.Padding,
	}).Debug("FsyncDir")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.FsyncDir(ch, &fuse.FsyncIn{InHeader: header, Fh: req.Fh, FsyncFlags: req.FsyncFlags, Padding: req.Padding})
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method FsyncDir not implemented")
	}
	return &pb.FsyncResponse{Status: &pb.Status{Code: int32(st)}}, nil
}
