package fuse2grpc

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/chiyutianyi/grpcfuse/pb"
)

// MockRawFileSystem is a mock implementation of fuse.RawFileSystem
type MockRawFileSystem struct {
	mock.Mock
}

func (m *MockRawFileSystem) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRawFileSystem) Lookup(ctx context.Context, parent *fuse.InHeader, name string, out *fuse.EntryOut) (status syscall.Errno) {
	args := m.Called(ctx, parent, name, out)
	return args.Get(0).(syscall.Errno)
}

func (m *MockRawFileSystem) GetAttr(ctx context.Context, in *fuse.InHeader, out *fuse.AttrOut) (code syscall.Errno) {
	args := m.Called(ctx, in, out)
	return args.Get(0).(syscall.Errno)
}

// Add other methods as needed for testing...

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		fs      fuse.RawFileSystem
		opts    []ServerOption
		wantErr bool
	}{
		{
			name:    "basic server",
			fs:      &MockRawFileSystem{},
			opts:    []ServerOption{},
			wantErr: false,
		},
		{
			name:    "server with custom threshold",
			fs:      &MockRawFileSystem{},
			opts:    []ServerOption{WithMsgSizeThreshold(2 * 1024 * 1024)},
			wantErr: false,
		},
		{
			name:    "server with custom logger",
			fs:      &MockRawFileSystem{},
			opts:    []ServerOption{WithLogger(logrus.New().WithField("test", "custom"))},
			wantErr: false,
		},
		{
			name:    "server with custom buffer pool",
			fs:      &MockRawFileSystem{},
			opts:    []ServerOption{WithBufferPool(&mockBufferPool{})},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.fs, tt.opts...)
			assert.NotNil(t, server)
			assert.Equal(t, tt.fs, server.GetFileSystem())
			
			if len(tt.opts) == 0 {
				assert.Equal(t, defaultMsgSizeThreshold, server.GetMsgSizeThreshold())
			}
		})
	}
}

func TestServer_String(t *testing.T) {
	tests := []struct {
		name           string
		mockReturn     string
		expectedValue  string
		expectedError  bool
	}{
		{
			name:           "successful string",
			mockReturn:     "test-filesystem",
			expectedValue:  "test-filesystem",
			expectedError:  false,
		},
		{
			name:           "empty string",
			mockReturn:     "",
			expectedValue:  "",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockRawFileSystem{}
			mockFS.On("String").Return(tt.mockReturn)

			server := NewServer(mockFS)
			ctx := context.Background()
			req := &pb.StringRequest{}

			resp, err := server.String(ctx, req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedValue, resp.Value)
			}

			mockFS.AssertExpectations(t)
		})
	}
}

func TestServer_SetMsgSizeThreshold(t *testing.T) {
	server := NewServer(&MockRawFileSystem{})
	
	// Test setting threshold
	newThreshold := 5 * 1024 * 1024 // 5MB
	server.SetMsgSizeThreshold(newThreshold)
	
	assert.Equal(t, newThreshold, server.GetMsgSizeThreshold())
	
	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(val int) {
			server.SetMsgSizeThreshold(val)
			done <- true
		}(i)
	}
	
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestServer_GetMsgSizeThreshold(t *testing.T) {
	server := NewServer(&MockRawFileSystem{})
	
	// Default threshold
	assert.Equal(t, defaultMsgSizeThreshold, server.GetMsgSizeThreshold())
	
	// Custom threshold
	customThreshold := 3 * 1024 * 1024
	server.SetMsgSizeThreshold(customThreshold)
	assert.Equal(t, customThreshold, server.GetMsgSizeThreshold())
}

func TestServer_GetFileSystem(t *testing.T) {
	mockFS := &MockRawFileSystem{}
	server := NewServer(mockFS)
	
	assert.Equal(t, mockFS, server.GetFileSystem())
}

func TestServer_GetBufferPool(t *testing.T) {
	server := NewServer(&MockRawFileSystem{})
	
	bp := server.GetBufferPool()
	assert.NotNil(t, bp)
}

func TestServer_Close(t *testing.T) {
	server := NewServer(&MockRawFileSystem{})
	
	// Test close with default buffer pool
	err := server.Close()
	assert.NoError(t, err)
	
	// Test close with custom buffer pool that implements Close
	customBP := &mockBufferPoolWithClose{}
	server = NewServer(&MockRawFileSystem{}, WithBufferPool(customBP))
	
	err = server.Close()
	assert.NoError(t, err)
	assert.True(t, customBP.closed)
}

func TestServer_WithMsgSizeThreshold(t *testing.T) {
	threshold := 10 * 1024 * 1024 // 10MB
	opt := WithMsgSizeThreshold(threshold)
	
	server := NewServer(&MockRawFileSystem{}, opt)
	assert.Equal(t, threshold, server.GetMsgSizeThreshold())
}

func TestServer_WithLogger(t *testing.T) {
	customLogger := logrus.New().WithField("test", "custom")
	opt := WithLogger(customLogger)
	
	server := NewServer(&MockRawFileSystem{}, opt)
	assert.Equal(t, customLogger, server.logger)
}

func TestServer_WithBufferPool(t *testing.T) {
	customBP := &mockBufferPool{}
	opt := WithBufferPool(customBP)
	
	server := NewServer(&MockRawFileSystem{}, opt)
	assert.Equal(t, customBP, server.GetBufferPool())
}

// Benchmark tests
func BenchmarkServer_String(b *testing.B) {
	mockFS := &MockRawFileSystem{}
	mockFS.On("String").Return("benchmark-fs")
	
	server := NewServer(mockFS)
	ctx := context.Background()
	req := &pb.StringRequest{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = server.String(ctx, req)
	}
}

func BenchmarkServer_SetMsgSizeThreshold(b *testing.B) {
	server := NewServer(&MockRawFileSystem{})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.SetMsgSizeThreshold(i)
	}
}

// Mock implementations
type mockBufferPool struct{}

func (m *mockBufferPool) Get() []byte {
	return make([]byte, 1024)
}

func (m *mockBufferPool) Put(buf []byte) {
	// Mock implementation
}

type mockBufferPoolWithClose struct {
	mockBufferPool
	closed bool
}

func (m *mockBufferPoolWithClose) Close() error {
	m.closed = true
	return nil
}

// Integration test helper
func TestServer_Integration(t *testing.T) {
	t.Skip("Integration test - requires real FUSE filesystem")
	
	// This test would demonstrate how to test with a real FUSE filesystem
	// In practice, you might want to use a test filesystem implementation
}

// Test server options
func TestServerOptions(t *testing.T) {
	t.Run("multiple options", func(t *testing.T) {
		mockFS := &MockRawFileSystem{}
		customLogger := logrus.New().WithField("test", "multiple")
		customThreshold := 15 * 1024 * 1024
		
		server := NewServer(mockFS,
			WithLogger(customLogger),
			WithMsgSizeThreshold(customThreshold),
		)
		
		assert.Equal(t, customLogger, server.logger)
		assert.Equal(t, customThreshold, server.GetMsgSizeThreshold())
	})
	
	t.Run("option chaining", func(t *testing.T) {
		mockFS := &MockRawFileSystem{}
		
		// Test that options can be chained
		opt1 := WithMsgSizeThreshold(1024)
		opt2 := WithLogger(logrus.New().WithField("test", "chained"))
		
		server := NewServer(mockFS, opt1, opt2)
		assert.NotNil(t, server)
	})
}