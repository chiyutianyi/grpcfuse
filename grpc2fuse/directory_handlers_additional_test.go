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
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/mock"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestOpenDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	var (
		in  fuse.OpenIn  = fuse.OpenIn{InHeader: TestInHeader, Flags: 0, Mode: 0}
		out fuse.OpenOut
		ch  <-chan struct{} = nil
	)

	// Case 1: client.OpenDir returns a response with status code 0
	expectedResponse := &pb.OpenDirResponse{
		Status: &pb.Status{Code: 0},
		OpenOut: &pb.OpenOut{
			Fh:        42,
			OpenFlags: 123,
			Padding:   456,
		},
	}

	client.EXPECT().OpenDir(gomock.Any(), gomock.Any()).Return(expectedResponse, nil)
	st := fs.OpenDir(ch, &in, &out)
	require.Equal(t, fuse.OK, st)
	require.Equal(t, uint64(42), out.Fh)
	require.Equal(t, uint32(123), out.OpenFlags)
	require.Equal(t, uint32(456), out.Padding)

	// Case 2: client.OpenDir returns a response with non-zero status code
	client.EXPECT().OpenDir(gomock.Any(), gomock.Any()).Return(&pb.OpenDirResponse{
		Status: &pb.Status{Code: 13},
	}, nil)
	st = fs.OpenDir(ch, &in, &out)
	require.Equal(t, fuse.Status(13), st)

	// Case 3: client.OpenDir returns Unimplemented error
	client.EXPECT().OpenDir(gomock.Any(), gomock.Any()).Return(nil, grpcstatus.Error(codes.Unimplemented, "Unimplemented"))
	st = fs.OpenDir(ch, &in, &out)
	require.Equal(t, fuse.ENOSYS, st)

	// Case 4: client.OpenDir returns other error
	client.EXPECT().OpenDir(gomock.Any(), gomock.Any()).Return(nil, grpcstatus.Error(codes.Internal, "Internal error"))
	st = fs.OpenDir(ch, &in, &out)
	require.Equal(t, fuse.EIO, st)
}

func TestReadDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	var (
		in  fuse.ReadIn = fuse.ReadIn{InHeader: TestInHeader, Fh: 42, Offset: 0, Size: 4096}
		ch  <-chan struct{} = nil
	)

	// Create a buffer for DirEntryList
	buf := make([]byte, 4096)
	list := fuse.NewDirEntryList(buf, 0)

	// Case 1: client.ReadDir returns a stream with entries and then EOF
	readDirClient := mock.NewMockRawFileSystem_ReadDirClient(ctrl)
	entries := []*pb.DirEntry{
		{Ino: 1, Name: []byte("file1"), Mode: 0100644},
		{Ino: 2, Name: []byte("file2"), Mode: 0100644},
		{Ino: 3, Name: []byte("dir1"), Mode: 040755},
	}

	client.EXPECT().ReadDir(gomock.Any(), gomock.Any()).Return(readDirClient, nil)
	readDirClient.EXPECT().Recv().Return(&pb.ReadDirResponse{
		Status:  &pb.Status{Code: 0},
		Entries: entries,
	}, nil)
	readDirClient.EXPECT().Recv().Return(nil, io.EOF)

	st := fs.ReadDir(ch, &in, list)
	require.Equal(t, fuse.OK, st)

	// Case 2: client.ReadDir returns a stream that returns Unimplemented error
	readDirClient = mock.NewMockRawFileSystem_ReadDirClient(ctrl)
	client.EXPECT().ReadDir(gomock.Any(), gomock.Any()).Return(readDirClient, nil)
	readDirClient.EXPECT().Recv().Return(nil, grpcstatus.Error(codes.Unimplemented, "Unimplemented"))

	list = fuse.NewDirEntryList(buf, 0)
	st = fs.ReadDir(ch, &in, list)
	require.Equal(t, fuse.ENOSYS, st)

	// Case 3: client.ReadDir returns a stream that returns a response with non-zero status
	readDirClient = mock.NewMockRawFileSystem_ReadDirClient(ctrl)
	client.EXPECT().ReadDir(gomock.Any(), gomock.Any()).Return(readDirClient, nil)
	readDirClient.EXPECT().Recv().Return(&pb.ReadDirResponse{
		Status: &pb.Status{Code: 13},
	}, nil)

	list = fuse.NewDirEntryList(buf, 0)
	st = fs.ReadDir(ch, &in, list)
	require.Equal(t, fuse.Status(13), st)

}

func TestReadDirPlus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	var (
		in  fuse.ReadIn = fuse.ReadIn{InHeader: TestInHeader, Fh: 42, Offset: 0, Size: 4096}
		ch  <-chan struct{} = nil
	)

	// Create a buffer for DirEntryList
	buf := make([]byte, 4096)
	list := fuse.NewDirEntryList(buf, 0)

	// Case 1: client.ReadDirPlus returns a stream with entries and then EOF
	readDirPlusClient := mock.NewMockRawFileSystem_ReadDirPlusClient(ctrl)
	entries := []*pb.DirEntry{
		{Ino: 1, Name: []byte("file1"), Mode: 0100644},
		{Ino: 2, Name: []byte("file2"), Mode: 0100644},
		{Ino: 3, Name: []byte("dir1"), Mode: 040755},
	}

	client.EXPECT().ReadDirPlus(gomock.Any(), gomock.Any()).Return(readDirPlusClient, nil)
	readDirPlusClient.EXPECT().Recv().Return(&pb.ReadDirResponse{
		Status:  &pb.Status{Code: 0},
		Entries: entries,
	}, nil)
	readDirPlusClient.EXPECT().Recv().Return(nil, io.EOF)

	st := fs.ReadDirPlus(ch, &in, list)
	require.Equal(t, fuse.OK, st)

	// Case 2: client.ReadDirPlus returns a stream that returns Unimplemented error
	readDirPlusClient = mock.NewMockRawFileSystem_ReadDirPlusClient(ctrl)
	client.EXPECT().ReadDirPlus(gomock.Any(), gomock.Any()).Return(readDirPlusClient, nil)
	readDirPlusClient.EXPECT().Recv().Return(nil, grpcstatus.Error(codes.Unimplemented, "Unimplemented"))

	list = fuse.NewDirEntryList(buf, 0)
	st = fs.ReadDirPlus(ch, &in, list)
	require.Equal(t, fuse.ENOSYS, st)

	// Case 3: client.ReadDirPlus returns a stream that returns a response with non-zero status
	readDirPlusClient = mock.NewMockRawFileSystem_ReadDirPlusClient(ctrl)
	client.EXPECT().ReadDirPlus(gomock.Any(), gomock.Any()).Return(readDirPlusClient, nil)
	readDirPlusClient.EXPECT().Recv().Return(&pb.ReadDirResponse{
		Status: &pb.Status{Code: 13},
	}, nil)

	list = fuse.NewDirEntryList(buf, 0)
	st = fs.ReadDirPlus(ch, &in, list)
	require.Equal(t, fuse.Status(13), st)

}

func TestFsyncDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	var (
		in  fuse.FsyncIn = fuse.FsyncIn{InHeader: TestInHeader, Fh: 42, FsyncFlags: 1, Padding: 0}
		ch  <-chan struct{} = nil
	)

	// Case 1: client.FsyncDir returns a response with status code 0
	client.EXPECT().FsyncDir(gomock.Any(), gomock.Any()).Return(&pb.FsyncResponse{
		Status: &pb.Status{Code: 0},
	}, nil)
	st := fs.FsyncDir(ch, &in)
	require.Equal(t, fuse.OK, st)

	// Case 2: client.FsyncDir returns a response with non-zero status code
	client.EXPECT().FsyncDir(gomock.Any(), gomock.Any()).Return(&pb.FsyncResponse{
		Status: &pb.Status{Code: 13},
	}, nil)
	st = fs.FsyncDir(ch, &in)
	require.Equal(t, fuse.Status(13), st)

	// Case 3: client.FsyncDir returns Unimplemented error
	client.EXPECT().FsyncDir(gomock.Any(), gomock.Any()).Return(nil, grpcstatus.Error(codes.Unimplemented, "Unimplemented"))
	st = fs.FsyncDir(ch, &in)
	require.Equal(t, fuse.ENOSYS, st)

	// Case 4: client.FsyncDir returns other error
	client.EXPECT().FsyncDir(gomock.Any(), gomock.Any()).Return(nil, grpcstatus.Error(codes.Internal, "Internal error"))
	st = fs.FsyncDir(ch, &in)
	require.Equal(t, fuse.EIO, st)
}

func TestReleaseDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	var (
		in fuse.ReleaseIn = fuse.ReleaseIn{InHeader: TestInHeader, Fh: 42, Flags: 0, ReleaseFlags: 0, LockOwner: 0}
	)

	// Test that ReleaseDir doesn't panic when client returns nil error
	client.EXPECT().ReleaseDir(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)
	require.NotPanics(t, func() {
		fs.ReleaseDir(&in)
	})

	// Test that ReleaseDir doesn't panic when client returns error
	client.EXPECT().ReleaseDir(gomock.Any(), gomock.Any()).Return(nil, grpcstatus.Error(codes.Internal, "Internal error"))
	require.NotPanics(t, func() {
		fs.ReleaseDir(&in)
	})
}
