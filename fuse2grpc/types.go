package fuse2grpc

import (
	"unsafe"

	"github.com/hanwen/go-fuse/v2/fuse"
)

// _Dirent is what we send to the kernel, but we offer DirEntry and
// DirEntryList to the user.
type _Dirent struct {
	Ino     uint64
	Off     uint64
	NameLen uint32
	Typ     uint32
}

const (
	direntSize   = uint32(unsafe.Sizeof(_Dirent{}))
	entryOutSize = uint32(unsafe.Sizeof(fuse.EntryOut{}))
)
