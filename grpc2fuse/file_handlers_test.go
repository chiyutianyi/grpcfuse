package grpc2fuse_test

import (
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/mock"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestReadDirRPCError(t *testing.T) {
	tests := []struct {
		name string
		call func(*mock.MockRawFileSystemClient, *fuse.ReadIn, *fuse.DirEntryList)
	}{
		{name: "ReadDir", call: func(client *mock.MockRawFileSystemClient, in *fuse.ReadIn, out *fuse.DirEntryList) {
			fs := grpc2fuse.NewFileSystem(client)
			require.Equal(t, fuse.EIO, fs.ReadDir(nil, in, out))
		}},
		{name: "ReadDirPlus", call: func(client *mock.MockRawFileSystemClient, in *fuse.ReadIn, out *fuse.DirEntryList) {
			fs := grpc2fuse.NewFileSystem(client)
			require.Equal(t, fuse.EIO, fs.ReadDirPlus(nil, in, out))
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := mock.NewMockRawFileSystemClient(ctrl)
			in := &fuse.ReadIn{InHeader: TestInHeader}
			out := fuse.NewDirEntryList(make([]byte, 4096), 0)

			if tt.name == "ReadDir" {
				client.EXPECT().ReadDir(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc failed"))
			} else {
				client.EXPECT().ReadDirPlus(gomock.Any(), gomock.Any()).Return(nil, errors.New("rpc failed"))
			}
			tt.call(client, in, out)
		})
	}
}

func TestRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	in := fuse.ReadIn{
		InHeader: TestInHeader,
		Size:     1,
	}

	buf := make([]byte, 100)
	idx := -1
	msg := []struct {
		buf *pb.ReadResponse
		err error
	}{
		{&pb.ReadResponse{Buffer: []byte("hello ")}, nil},
		{&pb.ReadResponse{Buffer: []byte("world")}, nil},
		{nil, io.EOF},
	}

	readclient := mock.NewMockRawFileSystem_ReadClient(ctrl)

	client.EXPECT().Read(gomock.Any(), gomock.Any()).Return(readclient, nil)
	readclient.EXPECT().Recv().Times(3).DoAndReturn(func() (*pb.ReadResponse, error) {
		idx++
		return msg[idx].buf, msg[idx].err
	})
	rs, status := fs.Read(nil, &in, buf)
	require.Equal(t, fuse.OK, status)
	out, status := rs.Bytes(buf)
	require.Equal(t, fuse.OK, status)
	require.Equal(t, "hello world", string(out))
}
