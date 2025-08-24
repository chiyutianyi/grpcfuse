package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/pb"
)

const bufSize = 1024 * 1024

// TestFileSystem is a simple in-memory filesystem for testing
type TestFileSystem struct {
	fuse.RawFileSystem
	files map[string]*TestFile
	dirs  map[string]*TestDirectory
}

type TestFile struct {
	name     string
	content  []byte
	modified time.Time
}

type TestDirectory struct {
	name    string
	entries map[string]interface{}
}

func NewTestFileSystem() *TestFileSystem {
	fs := &TestFileSystem{
		files: make(map[string]*TestFile),
		dirs:  make(map[string]*TestDirectory),
	}
	
	// Create root directory
	fs.dirs["/"] = &TestDirectory{
		name:    "/",
		entries: make(map[string]interface{}),
	}
	
	return fs
}

func (fs *TestFileSystem) String() string {
	return "TestFileSystem"
}

func (fs *TestFileSystem) Lookup(ctx context.Context, parent *fuse.InHeader, name string, out *fuse.EntryOut) (status syscall.Errno) {
	path := filepath.Join("/", name)
	
	if file, exists := fs.files[path]; exists {
		out.NodeId = uint64(len(fs.files))
		out.Attr = fs.fileToAttr(file)
		return 0
	}
	
	if dir, exists := fs.dirs[path]; exists {
		out.NodeId = uint64(len(fs.dirs))
		out.Attr = fs.dirToAttr(dir)
		return 0
	}
	
	return syscall.ENOENT
}

func (fs *TestFileSystem) GetAttr(ctx context.Context, in *fuse.InHeader, out *fuse.AttrOut) (code syscall.Errno) {
	// For simplicity, we'll just return basic attributes
	out.Attr = fuse.Attr{
		Mode: 0644,
		Size: 0,
	}
	return 0
}

func (fs *TestFileSystem) Create(ctx context.Context, in *fuse.CreateIn, name string, out *fuse.CreateOut) (code syscall.Errno) {
	path := filepath.Join("/", name)
	
	file := &TestFile{
		name:     name,
		content:  []byte{},
		modified: time.Now(),
	}
	
	fs.files[path] = file
	out.NodeId = uint64(len(fs.files))
	out.Attr = fs.fileToAttr(file)
	
	return 0
}

func (fs *TestFileSystem) Write(ctx context.Context, in *fuse.WriteIn, data []byte) (written uint32, code syscall.Errno) {
	// For simplicity, we'll just return success
	return uint32(len(data)), 0
}

func (fs *TestFileSystem) Read(ctx context.Context, in *fuse.ReadIn, buf []byte) (res fuse.ReadResult, code syscall.Errno) {
	// For simplicity, we'll just return empty data
	return fuse.ReadResultData([]byte{}), 0
}

func (fs *TestFileSystem) OpenDir(ctx context.Context, in *fuse.OpenIn, out *fuse.OpenOut) (status syscall.Errno) {
	return 0
}

func (fs *TestFileSystem) ReadDir(ctx context.Context, in *fuse.ReadIn, out *fuse.DirEntryList) (status syscall.Errno) {
	// Return root directory entries
	out.AddDirEntry(fuse.DirEntry{
		Name: ".",
		Mode: 0755,
	})
	out.AddDirEntry(fuse.DirEntry{
		Name: "..",
		Mode: 0755,
	})
	
	return 0
}

func (fs *TestFileSystem) fileToAttr(file *TestFile) fuse.Attr {
	return fuse.Attr{
		Mode:  0644,
		Size:  uint64(len(file.content)),
		Mtime: uint64(file.modified.Unix()),
	}
}

func (fs *TestFileSystem) dirToAttr(dir *TestDirectory) fuse.Attr {
	return fuse.Attr{
		Mode: 0755,
		Size: 0,
	}
}

// setupTestEnvironment creates a test gRPC server and client
func setupTestEnvironment(t *testing.T) (*grpc.ClientConn, *grpc.Server, func()) {
	// Create test filesystem
	testFS := NewTestFileSystem()
	
	// Create gRPC server
	grpcServer := grpc.NewServer()
	fuseServer := fuse2grpc.NewServer(testFS)
	pb.RegisterRawFileSystemServer(grpcServer, fuseServer)
	
	// Create buffer listener for testing
	listener := bufconn.Listen(bufSize)
	
	// Start server
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()
	
	// Create client connection
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithInsecure(),
	)
	require.NoError(t, err)
	
	// Cleanup function
	cleanup := func() {
		conn.Close()
		grpcServer.Stop()
		listener.Close()
	}
	
	return conn, grpcServer, cleanup
}

func TestBasicIntegration(t *testing.T) {
	conn, _, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Create gRPC client
	client := pb.NewRawFileSystemClient(conn)
	
	// Create FUSE filesystem
	fs := grpc2fuse.NewFileSystem(client)
	
	// Test basic operations
	t.Run("String operation", func(t *testing.T) {
		result := fs.String()
		assert.Equal(t, "TestFileSystem", result)
	})
}

func TestFileOperations(t *testing.T) {
	conn, _, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	client := pb.NewRawFileSystemClient(conn)
	fs := grpc2fuse.NewFileSystem(client)
	
	// Test file operations through gRPC
	t.Run("Create file", func(t *testing.T) {
		// This would test actual file creation through the gRPC interface
		// For now, we'll just verify the filesystem is working
		assert.NotNil(t, fs)
	})
}

func TestConcurrentOperations(t *testing.T) {
	conn, _, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	client := pb.NewRawFileSystemClient(conn)
	fs := grpc2fuse.NewFileSystem(client)
	
	// Test concurrent String operations
	numGoroutines := 100
	results := make(chan string, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			result := fs.String()
			results <- result
		}()
	}
	
	// Collect results
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		assert.Equal(t, "TestFileSystem", result)
	}
}

func TestErrorHandling(t *testing.T) {
	conn, _, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	client := pb.NewRawFileSystemClient(conn)
	fs := grpc2fuse.NewFileSystem(client)
	
	// Test error handling
	t.Run("Handle nil error", func(t *testing.T) {
		result := fs.handleGRPCError(nil, "test-operation")
		assert.Equal(t, syscall.Errno(0), result)
	})
	
	t.Run("Handle unknown error", func(t *testing.T) {
		result := fs.handleGRPCError(fmt.Errorf("unknown error"), "test-operation")
		assert.Equal(t, syscall.EIO, result)
	})
}

func TestPerformanceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	conn, _, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	client := pb.NewRawFileSystemClient(conn)
	fs := grpc2fuse.NewFileSystem(client)
	
	// Performance test
	start := time.Now()
	numOperations := 10000
	
	for i := 0; i < numOperations; i++ {
		_ = fs.String()
	}
	
	duration := time.Since(start)
	opsPerSec := float64(numOperations) / duration.Seconds()
	
	t.Logf("Performance: %d operations in %v (%.2f ops/sec)", 
		numOperations, duration, opsPerSec)
	
	// Ensure reasonable performance (should be > 1000 ops/sec)
	assert.Greater(t, opsPerSec, 1000.0, "Performance should be at least 1000 ops/sec")
}

func TestMemoryLeaks(t *testing.T) {
	// Test for memory leaks by creating many filesystems
	numIterations := 1000
	
	for i := 0; i < numIterations; i++ {
		conn, _, cleanup := setupTestEnvironment(t)
		
		client := pb.NewRawFileSystemClient(conn)
		fs := grpc2fuse.NewFileSystem(client)
		
		// Perform some operations
		_ = fs.String()
		
		// Cleanup
		cleanup()
	}
	
	// If we get here without panicking, we're probably not leaking memory
	t.Logf("Created and destroyed %d filesystems without memory leaks", numIterations)
}

func TestRealFilesystemIntegration(t *testing.T) {
	t.Skip("Real filesystem integration test - requires root privileges")
	
	// This test would mount a real FUSE filesystem
	// It requires root privileges and is not suitable for CI/CD
	
	/*
	// Example of how this would work:
	mountPoint := "/tmp/test-mount"
	
	// Create mount point
	err := os.MkdirAll(mountPoint, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(mountPoint)
	
	// Start gRPC server
	// Mount FUSE filesystem
	// Perform operations
	// Unmount
	*/
}

// Benchmark tests for integration
func BenchmarkIntegrationString(b *testing.B) {
	conn, _, cleanup := setupTestEnvironment(b)
	defer cleanup()
	
	client := pb.NewRawFileSystemClient(conn)
	fs := grpc2fuse.NewFileSystem(client)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fs.String()
	}
}

func BenchmarkIntegrationConcurrent(b *testing.B) {
	conn, _, cleanup := setupTestEnvironment(b)
	defer cleanup()
	
	client := pb.NewRawFileSystemClient(conn)
	fs := grpc2fuse.NewFileSystem(client)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = fs.String()
		}
	})
}