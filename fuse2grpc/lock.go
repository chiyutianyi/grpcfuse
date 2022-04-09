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

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func (s *server) GetLk(ctx context.Context, req *pb.LkRequest) (*pb.GetLkResponse, error) {
	var (
		header fuse.InHeader
		out    fuse.LkOut
	)

	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":    req.Header.NodeId,
		"fh":        req.Fh,
		"owner":     req.Owner,
		"lockStart": req.Lk.Start,
		"lockEnd":   req.Lk.End,
		"lockType":  req.Lk.Type,
		"lockPid":   req.Lk.Pid,
		"lkFlags":   req.LkFlags,
		"padding":   req.Padding,
	}).Debug("GetLk")
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := s.fs.GetLk(ch,
		&fuse.LkIn{
			InHeader: header,
			Fh:       req.Fh,
			Owner:    req.Owner,
			Lk: fuse.FileLock{
				Start: req.Lk.Start,
				End:   req.Lk.End,
				Typ:   req.Lk.Type,
				Pid:   req.Lk.Pid,
			},
		},
		&out)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method GetLk not implemented")
	}
	return &pb.GetLkResponse{Lk: &pb.FileLock{Start: out.Lk.Start, End: out.Lk.End, Type: out.Lk.Typ, Pid: out.Lk.Pid}, Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) SetLk(ctx context.Context, req *pb.LkRequest) (*pb.SetLkResponse, error) {
	return s.doSetLk(ctx, req, s.fs.SetLk, "SetLk")
}

func (s *server) SetLkw(ctx context.Context, req *pb.LkRequest) (*pb.SetLkResponse, error) {
	return s.doSetLk(ctx, req, s.fs.SetLk, "SetLkw")
}

func (s *server) doSetLk(
	ctx context.Context,
	req *pb.LkRequest,
	fn func(<-chan struct{}, *fuse.LkIn) fuse.Status,
	funcName string) (*pb.SetLkResponse, error) {
	var (
		header fuse.InHeader
	)

	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId":    req.Header.NodeId,
		"fh":        req.Fh,
		"owner":     req.Owner,
		"lockStart": req.Lk.Start,
		"lockEnd":   req.Lk.End,
		"lockType":  req.Lk.Type,
		"lockPid":   req.Lk.Pid,
		"lkFlags":   req.LkFlags,
		"padding":   req.Padding,
	}).Debug(funcName)
	toFuseInHeader(req.Header, &header)

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	st := fn(ch,
		&fuse.LkIn{
			InHeader: header,
			Fh:       req.Fh,
			Owner:    req.Owner,
			Lk: fuse.FileLock{
				Start: req.Lk.Start,
				End:   req.Lk.End,
				Typ:   req.Lk.Type,
				Pid:   req.Lk.Pid,
			},
		})
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, fmt.Sprintf("method %s not implemented", funcName))
	}
	return &pb.SetLkResponse{Status: &pb.Status{Code: int32(st)}}, nil
}
