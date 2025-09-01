/*
 * Copyright 2022 Han Xin, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fuse2grpc

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufferPoolConcurrency(t *testing.T) {
	bp := bufferPool{}
	const numGoroutines = 100
	const numOperations = 1000
	
	var wg sync.WaitGroup
	
	// 并发测试分配和释放
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			for j := 0; j < numOperations; j++ {
				size := uint32(1024 + (id*j)%4096) // 变化的大小
				buf := bp.AllocBuffer(size)
				
				require.NotNil(t, buf)
				require.GreaterOrEqual(t, len(buf), int(size))
				
				// 写入数据确保buffer可用
				for k := 0; k < len(buf); k++ {
					buf[k] = byte(k % 256)
				}
				
				bp.FreeBuffer(buf)
			}
		}(i)
	}
	
	wg.Wait()
}

func TestBufferPoolSizes(t *testing.T) {
	bp := bufferPool{}
	
	testSizes := []uint32{1, 100, 1024, 4096, 8192, 65536}
	
	for _, size := range testSizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			buf := bp.AllocBuffer(size)
			
			assert.NotNil(t, buf)
			assert.Equal(t, int(size), len(buf))
			assert.GreaterOrEqual(t, cap(buf), int(size))
			
			// 测试buffer内容可以修改
			if len(buf) > 0 {
				buf[0] = 0xFF
				assert.Equal(t, byte(0xFF), buf[0])
			}
			
			bp.FreeBuffer(buf)
		})
	}
}

func TestBufferPoolReuse(t *testing.T) {
	bp := bufferPool{}
	size := uint32(4096)
	
	// 分配一个buffer
	buf1 := bp.AllocBuffer(size)
	assert.NotNil(t, buf1)
	
	// 标记buffer
	if len(buf1) > 0 {
		buf1[0] = 0xFF
	}
	
	// 释放buffer
	bp.FreeBuffer(buf1)
	
	// 再次分配相同大小的buffer
	buf2 := bp.AllocBuffer(size)
	assert.NotNil(t, buf2)
	
	// 检查是否重用了同一个底层数组
	assert.Equal(t, cap(buf1), cap(buf2))
}

func TestBufferPoolEdgeCases(t *testing.T) {
	bp := bufferPool{}
	
	t.Run("zero size", func(t *testing.T) {
		buf := bp.AllocBuffer(0)
		assert.NotNil(t, buf)
		assert.Equal(t, 0, len(buf))
		bp.FreeBuffer(buf)
	})
	
	t.Run("nil buffer", func(t *testing.T) {
		// 不应该panic
		assert.NotPanics(t, func() {
			bp.FreeBuffer(nil)
		})
	})
	
	t.Run("invalid capacity buffer", func(t *testing.T) {
		// 创建一个不是页面大小倍数的slice
		invalidBuf := make([]byte, 100, 100)
		
		// 不应该panic
		assert.NotPanics(t, func() {
			bp.FreeBuffer(invalidBuf)
		})
	})
}

func BenchmarkBufferPoolAlloc(b *testing.B) {
	bp := bufferPool{}
	size := uint32(4096)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bp.AllocBuffer(size)
			bp.FreeBuffer(buf)
		}
	})
}

func BenchmarkBufferPoolAllocLarge(b *testing.B) {
	bp := bufferPool{}
	size := uint32(1024 * 1024) // 1MB
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bp.AllocBuffer(size)
		bp.FreeBuffer(buf)
	}
}

func BenchmarkBufferPoolConcurrentAlloc(b *testing.B) {
	bp := bufferPool{}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 测试不同大小的并发分配
			sizes := []uint32{1024, 2048, 4096, 8192}
			for _, size := range sizes {
				buf := bp.AllocBuffer(size)
				bp.FreeBuffer(buf)
			}
		}
	})
}