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

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.CreateOut
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
		"name":   req.Name,
		"flags":  req.Flags,
		"mode":   req.Mode,
	}).Debug("Create")
	toFuseInHeader(req.Header, &header)

	st := s.fs.Create(ctx.Done(), &fuse.CreateIn{InHeader: header, Flags: req.Flags, Mode: req.Mode}, req.Name, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
	}
	if st != fuse.OK {
		return &pb.CreateResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.CreateResponse{
		EntryOut: toPbEntryOut(&out.EntryOut),
		OpenOut: &pb.OpenOut{
			Fh:        out.Fh,
			OpenFlags: out.OpenFlags,
			Padding:   out.Padding,
		},
		Status: &pb.Status{Code: 0},
	}, nil
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

	st := s.fs.Open(ctx.Done(), &fuse.OpenIn{InHeader: header, Flags: req.OpenIn.Flags, Mode: req.OpenIn.Mode}, &out)
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
	var (
		header fuse.InHeader
		pos    int
		batch  []byte
	)
	ctx := stream.Context()

	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":    req.ReadIn.Header.NodeId,
		"fh":        req.ReadIn.Fh,
		"offset":    req.ReadIn.Offset,
		"size":      req.ReadIn.Size,
		"readFlags": req.ReadIn.ReadFlags,
	}).Debug("Read")
	toFuseInHeader(req.ReadIn.Header, &header)

	buf := s.buffers.AllocBuffer(req.ReadIn.Size)
	defer s.buffers.FreeBuffer(buf)

	res, st := s.fs.Read(ctx.Done(),
		&fuse.ReadIn{
			InHeader:  header,
			Fh:        req.ReadIn.Fh,
			Offset:    req.ReadIn.Offset,
			Size:      req.ReadIn.Size,
			ReadFlags: req.ReadIn.ReadFlags,
		}, buf)

	if st == fuse.ENOSYS {
		return status.Errorf(codes.Unimplemented, "method Read not implemented")
	}

	if st != fuse.OK {
		stream.Send(&pb.ReadResponse{Status: &pb.Status{Code: int32(st)}})
		return nil
	}

	data, st := res.Bytes(buf)

	if st != fuse.OK {
		stream.Send(&pb.ReadResponse{Status: &pb.Status{Code: int32(st)}})
		return nil
	}

	flushFunc := func() error {
		if batch == nil {
			return nil
		}
		if err := stream.Send(&pb.ReadResponse{
			Buffer: batch,
			Status: &pb.Status{Code: 0},
		}); err != nil {
			return err
		}
		return nil
	}

	for {
		if pos+s.msgSizeThreshold >= res.Size() {
			batch = data[pos:]
			flushFunc()
			break
		}

		batch = data[pos : pos+s.msgSizeThreshold]
		pos += s.msgSizeThreshold
		flushFunc()
	}
	return nil
}
func (s *server) Lseek(ctx context.Context, req *pb.LseekRequest) (*pb.LseekResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.LseekOut
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":  req.Header.NodeId,
		"fh":      req.Fh,
		"offset":  req.Offset,
		"whence":  req.Whence,
		"padding": req.Padding,
	}).Debug("Lseek")
	toFuseInHeader(req.Header, &header)

	st := s.fs.Lseek(ctx.Done(), &fuse.LseekIn{InHeader: header, Fh: req.Fh, Offset: req.Offset, Whence: req.Whence, Padding: req.Padding}, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Lseek not implemented")
	}
	if st != fuse.OK {
		return &pb.LseekResponse{Status: &pb.Status{Code: int32(st)}}, nil
	}
	return &pb.LseekResponse{
		Offset: out.Offset,
		Status: &pb.Status{Code: 0},
	}, nil
}
