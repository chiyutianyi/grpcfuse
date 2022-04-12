package grpc2fuse_test

import (
	"github.com/hanwen/go-fuse/v2/fuse"
)

var (
	TestOwner = fuse.Owner{Uid: 1, Gid: 1}

	TestCaller = fuse.Caller{
		Owner: TestOwner,
		Pid:   1,
	}

	TestInHeader = fuse.InHeader{
		NodeId: 1,
		Caller: TestCaller,
	}
)
