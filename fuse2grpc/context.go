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
	"context"
	"sync"
)

var cancelPool = sync.Pool{
	New: func() interface{} {
		return make(chan struct{})
	},
}

func newCancel(ctx context.Context) chan struct{} {
	cancel := cancelPool.Get().(chan struct{})
	go func() {
		cancel <- <-ctx.Done()
	}()
	return cancel
}

func releaseCancel(ch chan struct{}) {
	cancelPool.Put(ch)
}
