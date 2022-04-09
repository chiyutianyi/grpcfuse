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
