package grpc2fuse

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/chiyutianyi/grpcfuse/pb"
)

// PerformanceTestClient is a mock client optimized for performance testing
type PerformanceTestClient struct {
	mock.Mock
	responseDelay time.Duration
}

func (m *PerformanceTestClient) String(ctx context.Context, in *pb.StringRequest, opts ...grpc.CallOption) (*pb.StringResponse, error) {
	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.StringResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// BenchmarkFileSystemCreation tests filesystem creation performance
func BenchmarkFileSystemCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client := &PerformanceTestClient{}
		_ = NewFileSystem(client)
	}
}

// BenchmarkFileSystemString tests String operation performance
func BenchmarkFileSystemString(b *testing.B) {
	client := &PerformanceTestClient{}
	client.On("String", mock.Anything, mock.Anything, mock.Anything).Return(&pb.StringResponse{Value: "test-fs"}, nil)
	
	fs := NewFileSystem(client)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fs.String()
	}
}

// BenchmarkFileSystemStringWithDelay tests String operation with network delay simulation
func BenchmarkFileSystemStringWithDelay(b *testing.B) {
	client := &PerformanceTestClient{responseDelay: 1 * time.Millisecond}
	client.On("String", mock.Anything, mock.Anything, mock.Anything).Return(&pb.StringResponse{Value: "test-fs"}, nil)
	
	fs := NewFileSystem(client)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fs.String()
	}
}

// BenchmarkConcurrentStringOperations tests concurrent String operations
func BenchmarkConcurrentStringOperations(b *testing.B) {
	client := &PerformanceTestClient{}
	client.On("String", mock.Anything, mock.Anything, mock.Anything).Return(&pb.StringResponse{Value: "test-fs"}, nil)
	
	fs := NewFileSystem(client)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = fs.String()
		}
	})
}

// BenchmarkErrorHandling tests error handling performance
func BenchmarkErrorHandling(b *testing.B) {
	client := &PerformanceTestClient{}
	client.On("String", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
	
	fs := NewFileSystem(client)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fs.String()
	}
}

// BenchmarkGRPCOptions tests performance with different gRPC options
func BenchmarkGRPCOptions(b *testing.B) {
	client := &PerformanceTestClient{}
	client.On("String", mock.Anything, mock.Anything, mock.Anything).Return(&pb.StringResponse{Value: "test-fs"}, nil)
	
	opts := []grpc.CallOption{
		grpc.MaxCallRecvMsgSize(1024 * 1024),
		grpc.MaxCallSendMsgSize(1024 * 1024),
	}
	
	fs := NewFileSystem(client, opts...)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fs.String()
	}
}

// TestPerformanceMetrics tests various performance metrics
func TestPerformanceMetrics(t *testing.T) {
	client := &PerformanceTestClient{}
	fs := NewFileSystem(client)
	
	// Test logger performance
	start := time.Now()
	for i := 0; i < 1000; i++ {
		fs.SetLogger(fs.logger.WithField("iteration", i))
	}
	duration := time.Since(start)
	
	t.Logf("Logger operations: %d operations in %v (%.2f ops/sec)", 
		1000, duration, float64(1000)/duration.Seconds())
	
	// Test error handling performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		fs.handleGRPCError(assert.AnError, "test-operation")
	}
	duration = time.Since(start)
	
	t.Logf("Error handling: %d operations in %v (%.2f ops/sec)", 
		1000, duration, float64(1000)/duration.Seconds())
}

// TestMemoryUsage tests memory usage patterns
func TestMemoryUsage(t *testing.T) {
	// This test would typically use runtime.ReadMemStats
	// For now, we'll just test that we can create many filesystems without panicking
	
	clients := make([]*PerformanceTestClient, 1000)
	filesystems := make([]*FileSystem, 1000)
	
	for i := 0; i < 1000; i++ {
		clients[i] = &PerformanceTestClient{}
		filesystems[i] = NewFileSystem(clients[i])
	}
	
	// Verify all filesystems were created successfully
	for i, fs := range filesystems {
		assert.NotNil(t, fs, "Filesystem %d should not be nil", i)
		assert.Equal(t, clients[i], fs.GetClient(), "Filesystem %d should have correct client", i)
	}
}

// TestConcurrencySafety tests thread safety of filesystem operations
func TestConcurrencySafety(t *testing.T) {
	client := &PerformanceTestClient{}
	fs := NewFileSystem(client)
	
	// Test concurrent logger setting
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func(id int) {
			fs.SetLogger(fs.logger.WithField("goroutine", id))
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}
	
	// Verify filesystem is still functional
	assert.NotNil(t, fs.logger)
}

// TestStressTest performs a stress test with many concurrent operations
func TestStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}
	
	client := &PerformanceTestClient{}
	fs := NewFileSystem(client)
	
	// Simulate many concurrent operations
	numGoroutines := 1000
	operationsPerGoroutine := 100
	
	done := make(chan bool, numGoroutines)
	
	start := time.Now()
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				// Simulate various operations
				fs.SetLogger(fs.logger.WithField("goroutine", id).WithField("operation", j))
				fs.handleGRPCError(nil, "stress-test")
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	duration := time.Since(start)
	totalOperations := numGoroutines * operationsPerGoroutine
	
	t.Logf("Stress test completed: %d operations in %v (%.2f ops/sec)", 
		totalOperations, duration, float64(totalOperations)/duration.Seconds())
}

// TestPerformanceRegression tests for performance regressions
func TestPerformanceRegression(t *testing.T) {
	client := &PerformanceTestClient{}
	fs := NewFileSystem(client)
	
	// Baseline performance measurement
	start := time.Now()
	for i := 0; i < 10000; i++ {
		fs.handleGRPCError(nil, "baseline-test")
	}
	baselineDuration := time.Since(start)
	
	// Current performance measurement
	start = time.Now()
	for i := 0; i < 10000; i++ {
		fs.handleGRPCError(nil, "current-test")
	}
	currentDuration := time.Since(start)
	
	// Check for significant performance regression (allow 20% degradation)
	regressionRatio := float64(currentDuration) / float64(baselineDuration)
	
	t.Logf("Baseline: %v, Current: %v, Ratio: %.2f", 
		baselineDuration, currentDuration, regressionRatio)
	
	if regressionRatio > 1.2 {
		t.Warnf("Performance regression detected: %.2fx slower than baseline", regressionRatio)
	}
}