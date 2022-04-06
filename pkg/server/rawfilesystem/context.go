package rawfilesystem

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
