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
	"sync"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"

	"github.com/chiyutianyi/grpcfuse/pb"
)

// Default message size threshold (1MB < default gRPC message size limit 4MB)
const defaultMsgSizeThreshold = 1 << 20

// Server represents a FUSE to gRPC server implementation
type Server struct {
	pb.UnimplementedRawFileSystemServer

	fs fuse.RawFileSystem

	buffers bufferPool

	msgSizeThreshold int
	logger           *log.Entry
	mu               sync.RWMutex
}

// ServerOption is a function that configures a Server
type ServerOption func(*Server)

// WithMsgSizeThreshold sets the message size threshold for the server
func WithMsgSizeThreshold(threshold int) ServerOption {
	return func(s *Server) {
		s.msgSizeThreshold = threshold
	}
}

// WithLogger sets a custom logger for the server
func WithLogger(logger *log.Entry) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithBufferPool sets a custom buffer pool for the server
func WithBufferPool(bp bufferPool) ServerOption {
	return func(s *Server) {
		s.buffers = bp
	}
}

// NewServer returns a new FUSE to gRPC server with the given options
func NewServer(fs fuse.RawFileSystem, opts ...ServerOption) *Server {
	s := &Server{
		fs:               fs,
		buffers:          newBufferPool(),
		msgSizeThreshold: defaultMsgSizeThreshold,
		logger:           log.WithField("component", "fuse2grpc"),
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SetMsgSizeThreshold sets the message size threshold for streaming operations
func (s *Server) SetMsgSizeThreshold(threshold int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msgSizeThreshold = threshold
}

// GetMsgSizeThreshold returns the current message size threshold
func (s *Server) GetMsgSizeThreshold() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.msgSizeThreshold
}

// String returns the filesystem name
func (s *Server) String(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	s.logger.WithContext(ctx).Debug("String operation requested")
	
	value := s.fs.String()
	s.logger.WithContext(ctx).WithField("value", value).Debug("String operation completed")
	
	return &pb.StringResponse{Value: value}, nil
}

// GetFileSystem returns the underlying FUSE filesystem
func (s *Server) GetFileSystem() fuse.RawFileSystem {
	return s.fs
}

// GetBufferPool returns the buffer pool used by the server
func (s *Server) GetBufferPool() bufferPool {
	return s.buffers
}

// Close releases resources held by the server
func (s *Server) Close() error {
	s.logger.Info("Closing server and releasing resources")
	
	// Close buffer pool if it implements io.Closer
	if closer, ok := s.buffers.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}
