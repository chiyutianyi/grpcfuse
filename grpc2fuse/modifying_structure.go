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

func (fs *fileSystem) Mkdir(cancel <-chan struct{}, input *fuse.MkdirIn, name string, out *fuse.EntryOut) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Mkdir(ctx, &pb.MkdirRequest{
		Header: toPbHeader(&input.InHeader),
		Name:   name,
		Mode:   input.Mode,
		Umask:  input.Umask,
	}, fs.opts...)

	if st := dealGrpcError("Mkdir", err); st != fuse.OK {
		return st
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	toFuseEntryOut(out, res.EntryOut)
	return fuse.OK
}

func (fs *fileSystem) Unlink(cancel <-chan struct{}, header *fuse.InHeader, name string) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Unlink(ctx, &pb.UnlinkRequest{
		Header: toPbHeader(header),
		Name:   name,
	}, fs.opts...)

	if st := dealGrpcError("Unlink", err); st != fuse.OK {
		return st
	}
	return fuse.Status(res.Status.GetCode())
}

func (fs *fileSystem) Rmdir(cancel <-chan struct{}, header *fuse.InHeader, name string) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Rmdir(ctx, &pb.RmdirRequest{
		Header: toPbHeader(header),
		Name:   name,
	}, fs.opts...)

	if st := dealGrpcError("Rmdir", err); st != fuse.OK {
		return st
	}
	return fuse.Status(res.Status.GetCode())
}

func (fs *fileSystem) Rename(cancel <-chan struct{}, input *fuse.RenameIn, oldName string, newName string) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.Rename(ctx, &pb.RenameRequest{
		Header:  toPbHeader(&input.InHeader),
		OldName: oldName,
		NewName: newName,
		Newdir:  input.Newdir,
		Flags:   input.Flags,
		Padding: input.Padding,
	}, fs.opts...)

	if st := dealGrpcError("Rename", err); st != fuse.OK {
		return st
	}
	return fuse.Status(res.Status.GetCode())
}
