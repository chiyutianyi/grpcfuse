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

package fuse2grpc_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestStatFs(t *testing.T) {
	server, fs := startTestServices(t, 0)
	defer server.Stop()

	client, conn := newRawFileSystemClient(t, serverSocketPath)
	defer conn.Close()

	ctx, cancel := Context()
	defer cancel()

	// Case 1: StatFs returns OK with populated fields
	fs.EXPECT().StatFs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(cancel <-chan struct{}, in *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
			out.Blocks = 1000
			out.Bfree = 500
			out.Bavail = 400
			out.Files = 2000
			out.Ffree = 1000
			out.Bsize = 4096
			out.NameLen = 255
			out.Frsize = 4096
			return fuse.OK
		})

	resp, err := client.StatFs(ctx, &pb.StatfsRequest{
		Input: TestInHeader,
	})
	require.NoError(t, err)
	require.Equal(t, int32(0), resp.Status.Code)
	require.Equal(t, uint64(1000), resp.Blocks)
	require.Equal(t, uint64(500), resp.Bfree)
	require.Equal(t, uint64(400), resp.Bavail)
	require.Equal(t, uint64(2000), resp.Files)
	require.Equal(t, uint64(1000), resp.Ffree)
	require.Equal(t, uint32(4096), resp.Bsize)
	require.Equal(t, uint32(255), resp.NameLen)
	require.Equal(t, uint32(4096), resp.Frsize)

	// Case 2: StatFs returns ENOSYS, should translate to gRPC Unimplemented
	fs.EXPECT().StatFs(gomock.Any(), gomock.Any(), gomock.Any()).Return(fuse.ENOSYS)

	resp, err = client.StatFs(ctx, &pb.StatfsRequest{
		Input: TestInHeader,
	})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unimplemented, st.Code())

	// Case 3: StatFs returns a specific error code
	fs.EXPECT().StatFs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(cancel <-chan struct{}, in *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
			return fuse.Status(13)
		})

	resp, err = client.StatFs(ctx, &pb.StatfsRequest{
		Input: TestInHeader,
	})
	require.NoError(t, err)
	require.Equal(t, int32(13), resp.Status.Code)
}
