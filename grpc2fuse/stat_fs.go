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
	"github.com/chiyutianyi/grpcfuse/pb"

	"github.com/hanwen/go-fuse/v2/fuse"
)

func (fs *fileSystem) StatFs(cancel <-chan struct{}, in *fuse.InHeader, out *fuse.StatfsOut) (code fuse.Status) {
	ctx := newContext(cancel, in)
	defer releaseContext(ctx)

	res, err := fs.client.StatFs(ctx, &pb.StatfsRequest{
		Input: toPbHeader(in),
	}, fs.opts...)

	if st := dealGrpcError("StatFs", err); st != fuse.OK {
		return st
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	out.Blocks = res.Blocks
	out.Bfree = res.Bfree
	out.Bavail = res.Bavail
	out.Files = res.Files
	out.Ffree = res.Ffree
	out.Bsize = res.Bsize
	out.NameLen = res.NameLen
	out.Frsize = res.Frsize
	out.Padding = res.Padding
	//TODO out.Spare = res.Spare
	return fuse.OK
}
