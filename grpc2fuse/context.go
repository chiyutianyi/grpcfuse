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
	log "github.com/sirupsen/logrus"
)

type fuseContext struct {
	context.Context
	header *fuse.InHeader
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return &fuseContext{}
	},
}

func newContext(cancel <-chan struct{}, header *fuse.InHeader) *fuseContext {
	ctx := contextPool.Get().(*fuseContext)

	c, cancelFunc := context.WithCancel(context.Background())

	ctx.Context = c
	ctx.header = header
	go func() {
		<-cancel
		cancelFunc()
		log.Debugf("ino %v context done", header.NodeId)
	}()
	return ctx
}

func releaseContext(ctx *fuseContext) {
	contextPool.Put(ctx)
}
