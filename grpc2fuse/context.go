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
)

type fuseContext struct {
	context.Context
	cancel <-chan struct{}
}

func newContext(cancel <-chan struct{}) *fuseContext {
	return &fuseContext{Context: context.Background(), cancel: cancel}
}

func (ctx *fuseContext) Done() <-chan struct{} { return ctx.cancel }

func (ctx *fuseContext) Err() error { return context.Canceled }
