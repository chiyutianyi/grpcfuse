// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse2grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBufferPool(t *testing.T) {
	bp := newBufferPool()
	assert.NotNil(t, bp, "Buffer pool should not be nil")
}

func TestBufferPool_Get(t *testing.T) {
	bp := newBufferPool()
	
	// Test getting buffers
	buf1 := bp.Get()
	assert.NotNil(t, buf1, "Buffer should not be nil")
	assert.GreaterOrEqual(t, len(buf1), 1024, "Buffer should have minimum size")
	
	buf2 := bp.Get()
	assert.NotNil(t, buf2, "Second buffer should not be nil")
	assert.GreaterOrEqual(t, len(buf2), 1024, "Second buffer should have minimum size")
	
	// Buffers should be different instances
	assert.NotSame(t, buf1, buf2, "Buffers should be different instances")
}

func TestBufferPool_Put(t *testing.T) {
	bp := newBufferPool()
	
	// Get a buffer
	buf := bp.Get()
	originalLen := len(buf)
	
	// Put it back
	bp.Put(buf)
	
	// Get another buffer (might be the same one)
	buf2 := bp.Get()
	assert.NotNil(t, buf2, "Buffer should not be nil")
	assert.GreaterOrEqual(t, len(buf2), originalLen, "Buffer should maintain size")
}

func TestBufferPool_Reuse(t *testing.T) {
	bp := newBufferPool()
	
	// Get multiple buffers
	buffers := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		buffers[i] = bp.Get()
	}
	
	// Put them all back
	for _, buf := range buffers {
		bp.Put(buf)
	}
	
	// Get them again
	for i := 0; i < 10; i++ {
		buf := bp.Get()
		assert.NotNil(t, buf, "Buffer should not be nil")
	}
}

func TestBufferPool_Concurrent(t *testing.T) {
	bp := newBufferPool()
	numGoroutines := 100
	operationsPerGoroutine := 100
	
	done := make(chan bool, numGoroutines)
	
	// Test concurrent access
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				buf := bp.Get()
				assert.NotNil(t, buf, "Buffer should not be nil")
				
				// Simulate some work
				time.Sleep(time.Microsecond)
				
				bp.Put(buf)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestBufferPool_Stress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}
	
	bp := newBufferPool()
	numOperations := 10000
	
	start := time.Now()
	
	for i := 0; i < numOperations; i++ {
		buf := bp.Get()
		assert.NotNil(t, buf, "Buffer should not be nil")
		bp.Put(buf)
	}
	
	duration := time.Since(start)
	opsPerSec := float64(numOperations) / duration.Seconds()
	
	t.Logf("Stress test: %d operations in %v (%.2f ops/sec)", 
		numOperations, duration, opsPerSec)
	
	// Ensure reasonable performance
	assert.Greater(t, opsPerSec, 1000.0, "Performance should be at least 1000 ops/sec")
}

func TestBufferPool_MemoryEfficiency(t *testing.T) {
	bp := newBufferPool()
	
	// Get many buffers
	numBuffers := 1000
	buffers := make([][]byte, numBuffers)
	
	for i := 0; i < numBuffers; i++ {
		buffers[i] = bp.Get()
	}
	
	// Put them back
	for _, buf := range buffers {
		bp.Put(buf)
	}
	
	// Get them again to test reuse
	for i := 0; i < numBuffers; i++ {
		buf := bp.Get()
		assert.NotNil(t, buf, "Buffer should not be nil")
	}
}

func TestBufferPool_BufferSizes(t *testing.T) {
	bp := newBufferPool()
	
	// Test different buffer sizes
	sizes := []int{1024, 4096, 8192, 16384}
	
	for _, size := range sizes {
		buf := bp.Get()
		assert.NotNil(t, buf, "Buffer should not be nil")
		assert.GreaterOrEqual(t, len(buf), size, "Buffer should be large enough")
		
		// Resize buffer
		buf = buf[:size]
		assert.Equal(t, size, len(buf), "Buffer should have correct size")
		
		bp.Put(buf)
	}
}

func TestBufferPool_ZeroLength(t *testing.T) {
	bp := newBufferPool()
	
	// Test putting zero-length buffer
	buf := bp.Get()
	buf = buf[:0] // Make it zero length
	bp.Put(buf)
	
	// Should still work
	buf2 := bp.Get()
	assert.NotNil(t, buf2, "Buffer should not be nil")
}

func TestBufferPool_NilBuffer(t *testing.T) {
	bp := newBufferPool()
	
	// Test putting nil buffer (should not panic)
	bp.Put(nil)
	
	// Should still work
	buf := bp.Get()
	assert.NotNil(t, buf, "Buffer should not be nil")
}

// Benchmark tests
func BenchmarkBufferPool_Get(b *testing.B) {
	bp := newBufferPool()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bp.Get()
		_ = buf
	}
}

func BenchmarkBufferPool_GetPut(b *testing.B) {
	bp := newBufferPool()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bp.Get()
		bp.Put(buf)
	}
}

func BenchmarkBufferPool_Concurrent(b *testing.B) {
	bp := newBufferPool()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bp.Get()
			bp.Put(buf)
		}
	})
}

func BenchmarkBufferPool_Reuse(b *testing.B) {
	bp := newBufferPool()
	
	// Pre-allocate some buffers
	buffers := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		buffers[i] = bp.Get()
	}
	
	// Put them back
	for _, buf := range buffers {
		bp.Put(buf)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bp.Get()
		bp.Put(buf)
	}
}

// Test buffer pool with custom implementation
type customBufferPool struct {
	buffers [][]byte
}

func (cbp *customBufferPool) Get() []byte {
	if len(cbp.buffers) > 0 {
		buf := cbp.buffers[len(cbp.buffers)-1]
		cbp.buffers = cbp.buffers[:len(cbp.buffers)-1]
		return buf
	}
	return make([]byte, 1024)
}

func (cbp *customBufferPool) Put(buf []byte) {
	if buf != nil {
		cbp.buffers = append(cbp.buffers, buf)
	}
}

func TestCustomBufferPool(t *testing.T) {
	customBP := &customBufferPool{}
	
	// Test custom implementation
	buf1 := customBP.Get()
	assert.NotNil(t, buf1, "Custom buffer should not be nil")
	
	customBP.Put(buf1)
	
	buf2 := customBP.Get()
	assert.NotNil(t, buf2, "Custom buffer should not be nil")
}

// Test buffer pool edge cases
func TestBufferPool_EdgeCases(t *testing.T) {
	bp := newBufferPool()
	
	// Test with very large buffers
	largeBuf := make([]byte, 1024*1024) // 1MB
	bp.Put(largeBuf)
	
	// Test with small buffers
	smallBuf := make([]byte, 1)
	bp.Put(smallBuf)
	
	// Test with empty slice
	emptyBuf := make([]byte, 0)
	bp.Put(emptyBuf)
	
	// Should still work normally
	buf := bp.Get()
	assert.NotNil(t, buf, "Buffer should not be nil")
}

// Test buffer pool cleanup
func TestBufferPool_Cleanup(t *testing.T) {
	bp := newBufferPool()
	
	// Get many buffers
	numBuffers := 100
	buffers := make([][]byte, numBuffers)
	
	for i := 0; i < numBuffers; i++ {
		buffers[i] = bp.Get()
	}
	
	// Put them back
	for _, buf := range buffers {
		bp.Put(buf)
	}
	
	// Test that we can still get buffers after cleanup
	for i := 0; i < numBuffers; i++ {
		buf := bp.Get()
		assert.NotNil(t, buf, "Buffer should not be nil after cleanup")
	}
}
