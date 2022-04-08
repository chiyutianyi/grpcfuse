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
	"github.com/chiyutianyi/grpcfuse/pb"

	"github.com/hanwen/go-fuse/v2/fuse"
)

func toPbAttr(in *fuse.Attr) *pb.Attr {
	return &pb.Attr{
		Ino:       in.Ino,
		Size:      in.Size,
		Blocks:    in.Blocks,
		Atime:     in.Atime,
		Mtime:     in.Mtime,
		Ctime:     in.Ctime,
		Atimensec: in.Atimensec,
		Mtimensec: in.Mtimensec,
		Ctimensec: in.Ctimensec,
		Mode:      in.Mode,
		Nlink:     in.Nlink,
		Owner: &pb.Owner{
			Uid: in.Uid,
			Gid: in.Gid,
		},
		Rdev:    in.Rdev,
		Blksize: in.Blksize,
		Padding: in.Padding,
	}
}
