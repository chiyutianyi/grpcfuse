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

package grpc2fuse

import (
	"context"
	"fmt"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

const (
	defaultName = "grpcfuse"
)

// FileSystem represents a gRPC-based FUSE filesystem client
type FileSystem struct {
	fuse.RawFileSystem

	client pb.RawFileSystemClient
	opts   []grpc.CallOption
	logger *log.Entry
}

// NewFileSystem creates a new FileSystem instance with the given gRPC client
func NewFileSystem(client pb.RawFileSystemClient, opts ...grpc.CallOption) *FileSystem {
	return &FileSystem{
		RawFileSystem: fuse.NewDefaultRawFileSystem(),
		client:        client,
		opts:          opts,
		logger:        log.WithField("component", "grpc2fuse"),
	}
}

// String returns the filesystem name
func (fs *FileSystem) String() string {
	ctx := context.Background()
	res, err := fs.client.String(ctx, &pb.StringRequest{}, fs.opts...)
	if err != nil {
		fs.logger.WithError(err).Warn("Failed to get filesystem name, using default")
		return defaultName
	}
	return res.Value
}

// SetLogger sets a custom logger for the filesystem
func (fs *FileSystem) SetLogger(logger *log.Entry) {
	fs.logger = logger
}

// GetClient returns the underlying gRPC client
func (fs *FileSystem) GetClient() pb.RawFileSystemClient {
	return fs.client
}

// GetOptions returns the gRPC call options
func (fs *FileSystem) GetOptions() []grpc.CallOption {
	return fs.opts
}

// handleGRPCError converts gRPC errors to appropriate FUSE errors
func (fs *FileSystem) handleGRPCError(err error, operation string) syscall.Errno {
	if err == nil {
		return 0
	}

	fs.logger.WithError(err).WithField("operation", operation).Debug("gRPC operation failed")

	st, ok := status.FromError(err)
	if !ok {
		fs.logger.WithError(err).WithField("operation", operation).Warn("Unknown error type")
		return syscall.EIO
	}

	switch st.Code() {
	case codes.NotFound:
		return syscall.ENOENT
	case codes.PermissionDenied:
		return syscall.EACCES
	case codes.InvalidArgument:
		return syscall.EINVAL
	case codes.ResourceExhausted:
		return syscall.ENOMEM
	case codes.Unavailable:
		return syscall.EAGAIN
	case codes.DeadlineExceeded:
		return syscall.ETIMEDOUT
	default:
		fs.logger.WithError(err).WithField("operation", operation).Warn("Unhandled gRPC error")
		return syscall.EIO
	}
}
