package fuse2grpc_test

import (
	"io"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	server, fs := startTestServices(t, 5)
	defer server.Stop()

	client, conn := newRawFileSystemClient(t, serverSocketPath)
	defer conn.Close()

	req := &pb.ReadRequest{
		ReadIn: &pb.ReadIn{
			Header: TestInHeader,
			Fh:     1,
			Offset: 0,
			Size:   100,
		},
	}

	ctx, cancel := Context()
	defer cancel()

	testcases := []struct {
		buffer []byte
		status int32
		err    error
	}{
		{[]byte("hello"), 0, nil},
		{[]byte(" worl"), 0, nil},
		{[]byte("d"), 0, nil},
		{nil, 0, io.EOF},
	}

	fs.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Return(fuse.ReadResultData([]byte("hello world")), fuse.OK)

	stream, err := client.Read(ctx, req)
	require.NoError(t, err)

	idx := 0

	for idx < len(testcases) {
		resp, err := stream.Recv()
		if testcases[idx].err != nil {
			require.Error(t, err)
			require.Equal(t, testcases[idx].err, err)
			idx++
			continue
		}
		require.NoError(t, err)
		require.Equal(t, testcases[idx].status, resp.Status.Code)
		require.Equal(t, testcases[idx].buffer, resp.Buffer)
		idx++
	}
}
