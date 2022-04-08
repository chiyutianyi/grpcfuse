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

const (
	direntSize   = uint32(unsafe.Sizeof(_Dirent{}))
	entryOutSize = uint32(unsafe.Sizeof(fuse.EntryOut{}))
)
