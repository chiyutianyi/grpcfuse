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

package grpc2fuse

import (
	"context"
	"sync"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type fuseContext struct {
	context.Context
	header   *fuse.InHeader
	canceled bool
	cancel   <-chan struct{}
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return &fuseContext{}
	},
}

func newContext(cancel <-chan struct{}, header *fuse.InHeader) *fuseContext {
	ctx := contextPool.Get().(*fuseContext)
	ctx.Context = context.Background()
	ctx.canceled = false
	ctx.cancel = cancel
	ctx.header = header
	return ctx
}

func releaseContext(ctx *fuseContext) {
	contextPool.Put(ctx)
}

func (c *fuseContext) Cancel() {
	c.canceled = true
}

func (c *fuseContext) Canceled() bool {
	if c.canceled {
		return true
	}
	select {
	case <-c.cancel:
		return true
	default:
		return false
	}
}
