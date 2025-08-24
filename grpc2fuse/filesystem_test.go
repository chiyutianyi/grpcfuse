package grpc2fuse

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

// MockRawFileSystemClient is a mock implementation of pb.RawFileSystemClient
type MockRawFileSystemClient struct {
	mock.Mock
}

func (m *MockRawFileSystemClient) String(ctx context.Context, in *pb.StringRequest, opts ...grpc.CallOption) (*pb.StringResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.StringResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Add other methods as needed for testing...

func TestNewFileSystem(t *testing.T) {
	tests := []struct {
		name    string
		client  pb.RawFileSystemClient
		opts    []grpc.CallOption
		wantErr bool
	}{
		{
			name:    "valid client",
			client:  &MockRawFileSystemClient{},
			opts:    []grpc.CallOption{},
			wantErr: false,
		},
		{
			name:    "nil client",
			client:  nil,
			opts:    []grpc.CallOption{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.client == nil {
				// Test that we can handle nil client gracefully
				fs := NewFileSystem(tt.client, tt.opts...)
				assert.NotNil(t, fs)
				assert.Nil(t, fs.GetClient())
			} else {
				fs := NewFileSystem(tt.client, tt.opts...)
				assert.NotNil(t, fs)
				assert.Equal(t, tt.client, fs.GetClient())
				assert.Equal(t, tt.opts, fs.GetOptions())
			}
		})
	}
}

func TestFileSystem_String(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   *pb.StringResponse
		mockError     error
		expectedValue string
		expectedLogs  bool
	}{
		{
			name:           "successful response",
			mockResponse:   &pb.StringResponse{Value: "test-fs"},
			mockError:     nil,
			expectedValue: "test-fs",
			expectedLogs:  false,
		},
		{
			name:           "grpc error",
			mockResponse:   nil,
			mockError:     status.Error(codes.Internal, "internal error"),
			expectedValue: defaultName,
			expectedLogs:  true,
		},
		{
			name:           "nil response",
			mockResponse:   nil,
			mockError:     nil,
			expectedValue: defaultName,
			expectedLogs:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockRawFileSystemClient{}
			mockClient.On("String", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			fs := NewFileSystem(mockClient)
			result := fs.String()

			assert.Equal(t, tt.expectedValue, result)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestFileSystem_SetLogger(t *testing.T) {
	fs := NewFileSystem(&MockRawFileSystemClient{})
	
	customLogger := logrus.New().WithField("test", "custom")
	fs.SetLogger(customLogger)
	
	assert.Equal(t, customLogger, fs.logger)
}

func TestFileSystem_handleGRPCError(t *testing.T) {
	fs := NewFileSystem(&MockRawFileSystemClient{})
	
	tests := []struct {
		name     string
		err      error
		expected syscall.Errno
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: 0,
		},
		{
			name:     "not found error",
			err:      status.Error(codes.NotFound, "not found"),
			expected: syscall.ENOENT,
		},
		{
			name:     "permission denied error",
			err:      status.Error(codes.PermissionDenied, "permission denied"),
			expected: syscall.EACCES,
		},
		{
			name:     "invalid argument error",
			err:      status.Error(codes.InvalidArgument, "invalid argument"),
			expected: syscall.EINVAL,
		},
		{
			name:     "resource exhausted error",
			err:      status.Error(codes.ResourceExhausted, "resource exhausted"),
			expected: syscall.ENOMEM,
		},
		{
			name:     "unavailable error",
			err:      status.Error(codes.Unavailable, "unavailable"),
			expected: syscall.EAGAIN,
		},
		{
			name:     "deadline exceeded error",
			err:      status.Error(codes.DeadlineExceeded, "deadline exceeded"),
			expected: syscall.ETIMEDOUT,
		},
		{
			name:     "unknown error",
			err:      status.Error(codes.Unknown, "unknown"),
			expected: syscall.EIO,
		},
		{
			name:     "non-grpc error",
			err:      assert.AnError,
			expected: syscall.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fs.handleGRPCError(tt.err, "test-operation")
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileSystem_GetClient(t *testing.T) {
	mockClient := &MockRawFileSystemClient{}
	fs := NewFileSystem(mockClient)
	
	assert.Equal(t, mockClient, fs.GetClient())
}

func TestFileSystem_GetOptions(t *testing.T) {
	opts := []grpc.CallOption{grpc.MaxCallRecvMsgSize(1024)}
	fs := NewFileSystem(&MockRawFileSystemClient{}, opts...)
	
	assert.Equal(t, opts, fs.GetOptions())
}

// Benchmark tests
func BenchmarkFileSystem_String(b *testing.B) {
	mockClient := &MockRawFileSystemClient{}
	mockClient.On("String", mock.Anything, mock.Anything, mock.Anything).Return(&pb.StringResponse{Value: "test-fs"}, nil)
	
	fs := NewFileSystem(mockClient)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fs.String()
	}
}

// Integration test helper
func TestFileSystem_Integration(t *testing.T) {
	// This test demonstrates how to test with a real gRPC server
	// In practice, you might want to use a test server or mock
	t.Skip("Integration test - requires real gRPC server")
	
	// Example of how you might set up an integration test:
	/*
	ctx := context.Background()
	
	// Start a test gRPC server
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	
	grpcServer := grpc.NewServer()
	testFS := &testFileSystem{name: "integration-test-fs"}
	pb.RegisterRawFileSystemServer(grpcServer, fuse2grpc.NewServer(testFS))
	
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()
	
	// Connect client
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	
	client := pb.NewRawFileSystemClient(conn)
	fs := NewFileSystem(client)
	
	// Test the filesystem
	result := fs.String()
	assert.Equal(t, "integration-test-fs", result)
	*/
}

// testFileSystem is a simple implementation for testing
type testFileSystem struct {
	fuse.RawFileSystem
	name string
}

func (fs *testFileSystem) String() string {
	return fs.name
}