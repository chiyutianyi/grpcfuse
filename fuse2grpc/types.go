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

// DirEntryList holds the return value for READDIR and READDIRPLUS
// opcodes.
type DirEntryList struct {
	buf []byte
	// capacity of the underlying buffer
	size int
	// offset is the requested location in the directory. go-fuse
	// currently counts in number of directory entries, but this is an
	// implementation detail and may change in the future.
	// If `offset` and `fs.fileEntry.dirOffset` disagree, then a
	// directory seek has taken place.
	offset uint64
	// pointer to the last serialized _Dirent. Used by FixMode().
	lastDirent *_Dirent
}

const (
	direntSize   = uint32(unsafe.Sizeof(_Dirent{}))
	entryOutSize = uint32(unsafe.Sizeof(fuse.EntryOut{}))
)
