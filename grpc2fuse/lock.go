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

func (fs *fileSystem) GetLk(cancel <-chan struct{}, input *fuse.LkIn, out *fuse.LkOut) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.GetLk(ctx, &pb.LkRequest{
		Header: toPbHeader(&input.InHeader),
		Fh:     input.Fh,
		Owner:  input.Owner,
		Lk: &pb.FileLock{
			Start: input.Lk.Start,
			End:   input.Lk.End,
			Type:  input.Lk.Typ,
			Pid:   input.Lk.Pid,
		},
		LkFlags: input.LkFlags,
		Padding: input.Padding,
	}, fs.opts...)

	if st := dealGrpcError("GetLk", err); st != fuse.OK {
		return st
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}
	out.Lk.Start = res.Lk.Start
	out.Lk.End = res.Lk.End
	out.Lk.Typ = res.Lk.Type
	out.Lk.Pid = res.Lk.Pid
	return fuse.Status(res.Status.GetCode())
}

func (fs *fileSystem) SetLk(cancel <-chan struct{}, input *fuse.LkIn) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.SetLk(ctx, &pb.LkRequest{
		Header: toPbHeader(&input.InHeader),
		Fh:     input.Fh,
		Owner:  input.Owner,
		Lk: &pb.FileLock{
			Start: input.Lk.Start,
			End:   input.Lk.End,
			Type:  input.Lk.Typ,
			Pid:   input.Lk.Pid,
		},
		LkFlags: input.LkFlags,
		Padding: input.Padding,
	}, fs.opts...)

	if st := dealGrpcError("SetLk", err); st != fuse.OK {
		return st
	}

	return fuse.Status(res.Status.GetCode())
}

func (fs *fileSystem) SetLkw(cancel <-chan struct{}, input *fuse.LkIn) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.SetLkw(ctx, &pb.LkRequest{
		Header: toPbHeader(&input.InHeader),
		Fh:     input.Fh,
		Owner:  input.Owner,
		Lk: &pb.FileLock{
			Start: input.Lk.Start,
			End:   input.Lk.End,
			Type:  input.Lk.Typ,
			Pid:   input.Lk.Pid,
		},
		LkFlags: input.LkFlags,
		Padding: input.Padding,
	}, fs.opts...)

	if st := dealGrpcError("SetLkw", err); st != fuse.OK {
		return st
	}

	return fuse.Status(res.Status.GetCode())
}
