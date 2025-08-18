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

package grpc2fuse_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/mock"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestStatFs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	var (
		in  fuse.InHeader = TestInHeader
		out fuse.StatfsOut
		ch  <-chan struct{} = nil
	)

	// Case 1: client.StatFs returns a response with status code 0
	expectedResponse := &pb.StatfsResponse{
		Blocks:  1000,
		Bfree:   500,
		Bavail:  400,
		Files:   2000,
		Ffree:   1000,
		Bsize:   4096,
		NameLen: 255,
		Frsize:  4096,
		Status:  &pb.Status{Code: 0},
	}

	client.EXPECT().StatFs(gomock.Any(), gomock.Any()).Return(expectedResponse, nil)
	st := fs.StatFs(ch, &in, &out)
	require.Equal(t, fuse.OK, st)
	require.Equal(t, uint64(1000), out.Blocks)
	require.Equal(t, uint64(500), out.Bfree)
	require.Equal(t, uint64(400), out.Bavail)
	require.Equal(t, uint64(2000), out.Files)
	require.Equal(t, uint64(1000), out.Ffree)
	require.Equal(t, uint32(4096), out.Bsize)
	require.Equal(t, uint32(255), out.NameLen)
	require.Equal(t, uint32(4096), out.Frsize)

	// Case 2: client.StatFs returns a response with non-zero status code
	client.EXPECT().StatFs(gomock.Any(), gomock.Any()).Return(&pb.StatfsResponse{
		Status: &pb.Status{Code: 5},
	}, nil)
	st = fs.StatFs(ch, &in, &out)
	require.Equal(t, fuse.Status(5), st)

	// Case 3: client.StatFs returns Unimplemented error
	client.EXPECT().StatFs(gomock.Any(), gomock.Any()).Return(nil, grpcstatus.Error(codes.Unimplemented, "Unimplemented"))
	st = fs.StatFs(ch, &in, &out)
	require.Equal(t, fuse.ENOSYS, st)

	// Case 4: client.StatFs returns other error
	client.EXPECT().StatFs(gomock.Any(), gomock.Any()).Return(nil, grpcstatus.Error(codes.Internal, "Internal error"))
	st = fs.StatFs(ch, &in, &out)
	require.Equal(t, fuse.EIO, st)
}
