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
	"github.com/hanwen/go-fuse/v2/fuse"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func toFuseInHeader(in *pb.InHeader, out *fuse.InHeader) {
	out.Length = in.Length
	out.Opcode = in.Opcode
	out.Unique = in.Unique
	out.NodeId = in.NodeId
	out.Uid = in.Caller.Owner.Uid
	out.Gid = in.Caller.Owner.Gid
	out.Pid = in.Caller.Pid
	out.Padding = in.Padding
}

// modeToType converts a file *type* (as used in _Dirent.Typ)
// to a file *mode* (as used in syscall.Stat_t.Mode).
func typeToMode(typ uint32) uint32 {
	return (typ << 12) & 0170000
}
