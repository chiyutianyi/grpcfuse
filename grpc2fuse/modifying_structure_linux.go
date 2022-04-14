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

func (fs *fileSystem) Mknod(cancel <-chan struct{}, input *fuse.MknodIn, name string, out *fuse.EntryOut) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Mknod(ctx, &pb.MknodRequest{
		Header: toPbHeader(&input.InHeader),
		Name:   name,
		Mode:   input.Mode,
		Rdev:   input.Rdev,
		Umask:  input.Umask,
	}, fs.opts...)

	if st := dealGrpcError("Mknod", err); st != fuse.OK {
		return st
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	toFuseEntryOut(out, res.EntryOut)
	return fuse.OK
}
