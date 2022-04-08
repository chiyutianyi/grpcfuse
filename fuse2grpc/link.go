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

func (s *server) Link(ctx context.Context, req *pb.LinkRequest) (*pb.LinkResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.EntryOut
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":    req.Header.NodeId,
		"oldnodeid": req.Oldnodeid,
		"filename":  req.Filename,
	}).Debug("Link")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Link(ch, &fuse.LinkIn{InHeader: header, Oldnodeid: req.Oldnodeid}, req.Filename, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Link not implemented")
	}
	return &pb.LinkResponse{EntryOut: toPbEntryOut(&out), Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) Symlink(ctx context.Context, req *pb.SymlinkRequest) (*pb.SymlinkResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.EntryOut
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":    req.Header.NodeId,
		"pointedTo": req.PointedTo,
		"linkName":  req.LinkName,
	}).Debug("Symlink")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Symlink(ch, &header, req.PointedTo, req.LinkName, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Symlink not implemented")
	}
	return &pb.SymlinkResponse{EntryOut: toPbEntryOut(&out), Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) Readlink(ctx context.Context, req *pb.ReadlinkRequest) (*pb.ReadlinkResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
	}).Debug("Readlink")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	out, st := s.fs.Readlink(ch, &header)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Readlink not implemented")
	}
	return &pb.ReadlinkResponse{Out: out, Status: &pb.Status{Code: int32(st)}}, nil
}
