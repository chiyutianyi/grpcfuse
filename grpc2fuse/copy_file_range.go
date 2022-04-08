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

func (fs *fileSystem) CopyFileRange(cancel <-chan struct{}, input *fuse.CopyFileRangeIn) (written uint32, code fuse.Status) {
	ctx := newContext(cancel, &input.InHeader)
	defer releaseContext(ctx)

	res, err := fs.client.CopyFileRange(ctx, &pb.CopyFileRangeRequest{
		Header:    toPbHeader(&input.InHeader),
		FhIn:      input.FhIn,
		OffIn:     input.OffIn,
		NodeIdOut: input.NodeIdOut,
		FhOut:     input.FhOut,
		OffOut:    input.OffOut,
		Len:       input.Len,
		Flags:     input.Flags,
	}, fs.opts...)

	if st := dealError("CopyFileRange", err); st != fuse.OK {
		return 0, st
	}

	return uint32(res.Written), fuse.Status(res.Status.GetCode())
}
