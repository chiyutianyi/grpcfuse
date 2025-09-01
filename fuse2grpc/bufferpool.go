// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse2grpc

import (
	"os"
	"sync"
)

// bufferPool implements explicit memory management. It is used for
// minimizing the GC overhead of communicating with the kernel.
type bufferPool struct {
	lock sync.RWMutex

	// For each page size multiple a list of slice pointers.
	buffersBySize []*sync.Pool
}

var pageSize = os.Getpagesize()

func (p *bufferPool) getPool(pageCount int) *sync.Pool {
	// 先尝试读锁获取已存在的pool
	p.lock.RLock()
	if pageCount < len(p.buffersBySize) && p.buffersBySize[pageCount] != nil {
		pool := p.buffersBySize[pageCount]
		p.lock.RUnlock()
		return pool
	}
	p.lock.RUnlock()
	
	// 需要创建新pool时才使用写锁
	p.lock.Lock()
	defer p.lock.Unlock()
	
	// 双重检查，防止在获取写锁期间被其他goroutine创建
	if pageCount < len(p.buffersBySize) && p.buffersBySize[pageCount] != nil {
		return p.buffersBySize[pageCount]
	}
	
	for len(p.buffersBySize) < pageCount+1 {
		p.buffersBySize = append(p.buffersBySize, nil)
	}
	if p.buffersBySize[pageCount] == nil {
		p.buffersBySize[pageCount] = &sync.Pool{
			New: func() interface{} { return make([]byte, pageSize*pageCount) },
		}
	}
	return p.buffersBySize[pageCount]
}

// AllocBuffer creates a buffer of at least the given size. After use,
// it should be deallocated with FreeBuffer().
func (p *bufferPool) AllocBuffer(size uint32) []byte {
	sz := int(size)
	if sz < pageSize {
		sz = pageSize
	}

	if sz%pageSize != 0 {
		sz += pageSize
	}
	pages := sz / pageSize

	b := p.getPool(pages).Get().([]byte)
	return b[:size]
}

// FreeBuffer takes back a buffer if it was allocated through
// AllocBuffer.  It is not an error to call FreeBuffer() on a slice
// obtained elsewhere.
func (p *bufferPool) FreeBuffer(slice []byte) {
	if slice == nil {
		return
	}
	if cap(slice)%pageSize != 0 || cap(slice) == 0 {
		return
	}
	pages := cap(slice) / pageSize
	slice = slice[:cap(slice)]

	p.getPool(pages).Put(slice)
}

// BufferPoolAccessor 提供对bufferPool的测试访问
type BufferPoolAccessor struct {
	pool bufferPool
}

// AllocBuffer 分配缓冲区
func (bpa *BufferPoolAccessor) AllocBuffer(size uint32) []byte {
	return bpa.pool.AllocBuffer(size)
}

// FreeBuffer 释放缓冲区
func (bpa *BufferPoolAccessor) FreeBuffer(buf []byte) {
	bpa.pool.FreeBuffer(buf)
}
