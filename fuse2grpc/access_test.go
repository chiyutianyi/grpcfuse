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

func TestAccess(t *testing.T) {
	server, fs := startTestServices(t)
	defer server.Stop()

	client, conn := newRawFileSystemClient(t, serverSocketPath)
	defer conn.Close()

	req := &pb.AccessRequest{
		Header: &pb.InHeader{
			NodeId: 1,
			Caller: &pb.Caller{
				Owner: &pb.Owner{Uid: 1, Gid: 1},
				Pid:   1,
			},
		},
	}

	testcases := []struct {
		req    *pb.AccessRequest
		status fuse.Status
		code   int32
		err    codes.Code
	}{
		{req, fuse.OK, 0, codes.OK},
		{req, fuse.EACCES, 13, codes.OK},
		{req, fuse.ENOSYS, 0, codes.Unimplemented},
	}

	ctx, cancel := Context()
	defer cancel()

	for _, testcase := range testcases {
		fs.EXPECT().Access(gomock.Any(), gomock.Any()).Return(testcase.status)

		resp, err := client.Access(ctx, testcase.req)
		if testcase.err != codes.OK {
			require.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, testcase.err, st.Code())
			continue
		}
		require.NoError(t, err)
		require.Equal(t, testcase.code, resp.Status.GetCode())
	}
}
