package grpcfuse

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
