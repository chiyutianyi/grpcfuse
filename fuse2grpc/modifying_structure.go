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

func (s *server) Mkdir(ctx context.Context, req *pb.MkdirRequest) (*pb.MkdirResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.EntryOut
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
		"name":   req.Name,
		"mode":   req.Mode,
		"umask":  req.Umask,
	}).Debug("Mknod")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Mkdir(ch, &fuse.MkdirIn{InHeader: header, Mode: req.Mode, Umask: req.Umask}, req.Name, &out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Mkdir not implemented")
	}
	return &pb.MkdirResponse{EntryOut: toPbEntryOut(&out), Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) Unlink(ctx context.Context, req *pb.UnlinkRequest) (*pb.UnlinkResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
		"name":   req.Name,
	}).Debug("Unlink")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Unlink(ch, &header, req.Name)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Unlink not implemented")
	}
	return &pb.UnlinkResponse{Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) Rmdir(ctx context.Context, req *pb.RmdirRequest) (*pb.RmdirResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
		"name":   req.Name,
	}).Debug("Rmdir")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Rmdir(ch, &header, req.Name)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Rmdir not implemented")
	}
	return &pb.RmdirResponse{Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) Rename(ctx context.Context, req *pb.RenameRequest) (*pb.RenameResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":  req.Header.NodeId,
		"oldName": req.OldName,
		"newName": req.NewName,
		"flags":   req.Flags,
		"padding": req.Padding,
	}).Debug("Rename")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.Rename(ch, &fuse.RenameIn{InHeader: header, Newdir: req.Newdir, Flags: req.Flags, Padding: req.Padding}, req.OldName, req.NewName)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method Rename not implemented")
	}
	return &pb.RenameResponse{Status: &pb.Status{Code: int32(st)}}, nil
}
