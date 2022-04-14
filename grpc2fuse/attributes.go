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

func (fs *fileSystem) GetAttr(cancel <-chan struct{}, in *fuse.GetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.GetAttr(ctx, &pb.GetAttrRequest{
		Header: toPbHeader(&in.InHeader),
	}, fs.opts...)

	if st := dealGrpcError("GetAttr", err); st != fuse.OK {
		return st
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}
	toFuseAttrOut(out, res.GetAttrOut())
	return fuse.OK
}

func (fs *fileSystem) SetAttr(cancel <-chan struct{}, in *fuse.SetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.SetAttr(ctx, &pb.SetAttrRequest{
		Header:    toPbHeader(&in.InHeader),
		Valid:     in.Valid,
		Padding:   in.Padding,
		Fh:        in.Fh,
		Size:      in.Size,
		LockOwner: in.LockOwner,
		Atime:     in.Atime,
		Mtime:     in.Mtime,
		Ctime:     in.Ctime,
		Atimensec: in.Atimensec,
		Mtimensec: in.Mtimensec,
		Ctimensec: in.Ctimensec,
		Mode:      in.Mode,
		Unused4:   in.Unused4,
		Owner: &pb.Owner{
			Uid: in.Uid,
			Gid: in.Gid,
		},
		Unused5: in.Unused5,
	}, fs.opts...)

	if st := dealGrpcError("SetAttr", err); st != fuse.OK {
		return st
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}
	toFuseAttrOut(out, res.GetAttrOut())
	return fuse.OK
}
