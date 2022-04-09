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
