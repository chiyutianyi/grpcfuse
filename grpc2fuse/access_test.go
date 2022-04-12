package grpc2fuse_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/mock"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	in := fuse.AccessIn{
		InHeader: TestInHeader,
	}

	testcases := []struct {
		in  *fuse.AccessIn
		res *pb.AccessResponse
		err error
		st  fuse.Status
	}{
		{&in, &pb.AccessResponse{Status: &pb.Status{Code: 0}}, nil, fuse.OK},
		{&in, &pb.AccessResponse{Status: &pb.Status{Code: 13}}, nil, fuse.EACCES},
		{&in, &pb.AccessResponse{Status: &pb.Status{Code: 0}}, status.Error(codes.Unimplemented, "Unimplemented"), fuse.ENOSYS},
	}

	for _, testcase := range testcases {
		client.EXPECT().Access(gomock.Any(), gomock.Any()).Return(testcase.res, testcase.err)
		require.Equal(t, testcase.st, fs.Access(nil, testcase.in))
	}
}
