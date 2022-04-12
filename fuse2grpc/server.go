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

	"github.com/chiyutianyi/grpcfuse/pb"
)

// msgSizeThreshold 1mb < default grpc message size limit 4mb
const msgSizeThreshold = 1 << 20

type server struct {
	pb.UnimplementedRawFileSystemServer

	fs fuse.RawFileSystem

	buffers bufferPool

	msgSizeThreshold int
}

// NewServer returns a new loopback server.
func NewServer(fs fuse.RawFileSystem) *server {
	return &server{fs: fs, buffers: bufferPool{}, msgSizeThreshold: msgSizeThreshold}
}

func (s *server) SetMsgSizeThreshold(threshold int) {
	s.msgSizeThreshold = threshold
}

func (s *server) String(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	grpc_logrus.Extract(ctx).Debug("String")
	return &pb.StringResponse{Value: s.fs.String()}, nil
}
