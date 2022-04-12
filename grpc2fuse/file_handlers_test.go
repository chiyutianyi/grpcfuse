package grpc2fuse_test

import (
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
