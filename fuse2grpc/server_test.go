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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chiyutianyi/grpcfuse/mock"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockFS := mock.NewMockRawFileSystem(ctrl)
	server := NewServer(mockFS)
	
	assert.NotNil(t, server)
	assert.Equal(t, msgSizeThreshold, server.msgSizeThreshold)
}

func TestServerSetMsgSizeThreshold(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockFS := mock.NewMockRawFileSystem(ctrl)
	server := NewServer(mockFS)
	
	newThreshold := 512 * 1024 // 512KB
	server.SetMsgSizeThreshold(newThreshold)
	
	assert.Equal(t, newThreshold, server.msgSizeThreshold)
}

func TestServerString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockFS := mock.NewMockRawFileSystem(ctrl)
	server := NewServer(mockFS)
	
	expectedName := "test-filesystem"
	mockFS.EXPECT().String().Return(expectedName)
	
	response, err := server.String(context.Background(), &pb.StringRequest{})
	
	require.NoError(t, err)
	assert.Equal(t, expectedName, response.Value)
}

func TestServerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockFS := mock.NewMockRawFileSystem(ctrl)
	server := NewServer(mockFS)
	
	req := &pb.CreateRequest{
		Header: &pb.InHeader{
			Length:  100,
			Opcode:  101,
			Unique:  102,
			NodeId:  1,
			Caller:  &pb.Caller{Owner: &pb.Owner{Uid: 1000, Gid: 1000}, Pid: 123},
			Padding: 0,
		},
		Name:  "testfile.txt",
		Flags: 0644,
		Mode:  0644,
	}
	
	tests := []struct {
		name           string
		fuseStatus     fuse.Status
		expectedStatus int32
	}{
		{
			name:           "successful create",
			fuseStatus:     fuse.OK,
			expectedStatus: 0,
		},
		{
			name:           "permission denied",
			fuseStatus:     fuse.EACCES,
			expectedStatus: int32(fuse.EACCES),
		},
		{
			name:           "file exists",
			fuseStatus:     fuse.Status(17), // EEXIST
			expectedStatus: 17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedOut := fuse.CreateOut{
				EntryOut: fuse.EntryOut{
					NodeId: 2,
					Attr: fuse.Attr{
						Ino:  2,
						Mode: 0644,
						Size: 0,
					},
				},
				OpenOut: fuse.OpenOut{
					Fh: 1,
				},
			}
			
			mockFS.EXPECT().Create(gomock.Any(), gomock.Any(), "testfile.txt", gomock.Any()).
				DoAndReturn(func(cancel <-chan struct{}, input *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {
					if tt.fuseStatus == fuse.OK {
						*out = expectedOut
					}
					return tt.fuseStatus
				})
			
			response, err := server.Create(context.Background(), req)
			
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.Status.Code)
			
			if tt.fuseStatus == fuse.OK {
				assert.NotNil(t, response.EntryOut)
				assert.NotNil(t, response.OpenOut)
				assert.Equal(t, uint64(2), response.EntryOut.NodeId)
				assert.Equal(t, uint64(1), response.OpenOut.Fh)
			}
		})
	}
}