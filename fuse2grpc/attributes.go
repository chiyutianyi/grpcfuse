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

func (s *server) GetAttr(ctx context.Context, req *pb.GetAttrRequest) (*pb.GetAttrResponse, error) {
	var (
		out    fuse.AttrOut
		header fuse.InHeader
	)
	grpc_logrus.Extract(ctx).WithFields(log.Fields{
		"nodeId": req.Header.NodeId,
	}).Debug("GetAttr")
	toFuseInHeader(req.Header, &header)

	st := s.fs.GetAttr(ctx.Done(), &fuse.GetAttrIn{InHeader: header}, &out)
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

	st := s.fs.SetAttr(ctx.Done(),
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
