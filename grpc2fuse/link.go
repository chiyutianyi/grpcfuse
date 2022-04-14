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

func (fs *fileSystem) Link(cancel <-chan struct{}, input *fuse.LinkIn, filename string, out *fuse.EntryOut) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Link(ctx, &pb.LinkRequest{
		Header:    toPbHeader(&input.InHeader),
		Oldnodeid: input.Oldnodeid,
		Filename:  filename,
	}, fs.opts...)

	if st := dealGrpcError("Link", err); st != fuse.OK {
		return st
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	toFuseEntryOut(out, res.EntryOut)
	return fuse.OK
}

func (fs *fileSystem) Symlink(cancel <-chan struct{}, header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Symlink(ctx, &pb.SymlinkRequest{
		Header:    toPbHeader(header),
		PointedTo: pointedTo,
		LinkName:  linkName,
	}, fs.opts...)

	if st := dealGrpcError("Symlink", err); st != fuse.OK {
		return st
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	toFuseEntryOut(out, res.EntryOut)
	return fuse.OK
}

func (fs *fileSystem) Readlink(cancel <-chan struct{}, header *fuse.InHeader) (out []byte, code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Readlink(ctx, &pb.ReadlinkRequest{
		Header: toPbHeader(header),
	}, fs.opts...)

	if st := dealGrpcError("Readlink", err); st != fuse.OK {
		return nil, st
	}

	return res.GetOut(), fuse.Status(res.Status.GetCode())
}
