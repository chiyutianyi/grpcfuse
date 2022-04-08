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

func (fs *fileSystem) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, dest []byte) (sz uint32, code fuse.Status) {
	ctx := newContext(cancel, header)
	defer releaseContext(ctx)

	res, err := fs.client.GetXAttr(ctx, &pb.GetXAttrRequest{
		Header: toPbHeader(header),
		Attr:   attr,
		Dest:   dest,
	}, fs.opts...)

	if st := dealGrpcError("GetXAttr", err); st != fuse.OK {
		return 0, st
	}
	return res.Size, fuse.Status(res.Status.Code)
}

func (fs *fileSystem) ListXAttr(cancel <-chan struct{}, header *fuse.InHeader, dest []byte) (uint32, fuse.Status) {
	ctx := newContext(cancel, header)
	defer releaseContext(ctx)

	res, err := fs.client.ListXAttr(ctx, &pb.ListXAttrRequest{
		Header: toPbHeader(header),
		Dest:   dest,
	}, fs.opts...)

	if st := dealGrpcError("ListXAttr", err); st != fuse.OK {
		return 0, st
	}
	return res.Size, fuse.Status(res.Status.Code)
}

func (fs *fileSystem) RemoveXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string) (code fuse.Status) {
	ctx := newContext(cancel, header)
	defer releaseContext(ctx)

	res, err := fs.client.RemoveXAttr(ctx, &pb.RemoveXAttrRequest{
		Header: toPbHeader(header),
		Attr:   attr,
	}, fs.opts...)

	if st := dealGrpcError("RemoveXAttr", err); st != fuse.OK {
		return st
	}
	return fuse.Status(res.Status.Code)
}
