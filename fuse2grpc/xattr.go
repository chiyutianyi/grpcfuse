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

func (s *server) GetXAttr(ctx context.Context, req *pb.GetXAttrRequest) (*pb.GetXAttrResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
		"attr":   req.Attr,
		"dest":   req.Dest,
	}).Debug("GetXAttr")
	toFuseInHeader(req.Header, &header)

	sz, st := s.fs.GetXAttr(ctx.Done(), &header, req.Attr, req.Dest)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method GetXAttr not implemented")
	}
	return &pb.GetXAttrResponse{Size: sz, Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) ListXAttr(ctx context.Context, req *pb.ListXAttrRequest) (*pb.ListXAttrResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
		"dest":   req.Dest,
	}).Debug("ListXAttr")
	toFuseInHeader(req.Header, &header)

	sz, st := s.fs.ListXAttr(ctx.Done(), &header, req.Dest)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method ListXAttr not implemented")
	}
	return &pb.ListXAttrResponse{Size: sz, Status: &pb.Status{Code: int32(st)}}, nil
}

func (s *server) RemoveXAttr(ctx context.Context, req *pb.RemoveXAttrRequest) (*pb.RemoveXAttrResponse, error) {
	var (
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
		"attr":   req.Attr,
	}).Debug("RemoveXAttr")
	toFuseInHeader(req.Header, &header)

	st := s.fs.RemoveXAttr(ctx.Done(), &header, req.Attr)
	if st == fuse.ENOSYS {
		return nil, status.Errorf(codes.Unimplemented, "method RemoveXAttr not implemented")
	}
	return &pb.RemoveXAttrResponse{Status: &pb.Status{Code: int32(st)}}, nil
}
