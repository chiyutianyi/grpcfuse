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
	"context"

	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/chiyutianyi/grpcfuse/pb"
)

const (
	defaultName = "grpcfuse"
)

type fileSystem struct {
	fuse.RawFileSystem

	client pb.RawFileSystemClient
	opts   []grpc.CallOption
}

// NewFileSystem creates a new file system.
func NewFileSystem(client pb.RawFileSystemClient, opts ...grpc.CallOption) *fileSystem {
	return &fileSystem{
		RawFileSystem: fuse.NewDefaultRawFileSystem(),
		client:        client,
		opts:          opts,
	}
}

func (fs *fileSystem) String() string {
	res, err := fs.client.String(context.TODO(), &pb.StringRequest{}, fs.opts...)
	if err != nil {
		log.Errorf("String: %v", err)
		return defaultName
	}
	return res.Value
}

func (fs *fileSystem) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) (status fuse.Status) {
	ctx := newContext(cancel, header)
	defer releaseContext(ctx)

	res, err := fs.client.Lookup(ctx, &pb.LookupRequest{
		Header: toPbHeader(header),
		Name:   name,
	}, fs.opts...)
	if err != nil {
		log.Errorf("Lookup: %v", err)
		return fuse.EIO
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}
	toFuseEntryOut(out, res.EntryOut)
	return fuse.OK
}

func (fs *fileSystem) Forget(nodeid, nlookup uint64) {
	_, err := fs.client.Forget(context.TODO(), &pb.ForgetRequest{Nodeid: nodeid, Nlookup: nlookup}, fs.opts...)
	if err != nil {
		log.Errorf("Forget: %v", err)
	}
}

func (fs *fileSystem) GetAttr(cancel <-chan struct{}, in *fuse.GetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	res, err := fs.client.GetAttr(ctx, &pb.GetAttrRequest{
		Header: toPbHeader(&in.InHeader),
	}, fs.opts...)

	if err != nil {
		log.Errorf("GetAttr: %v", err)
		return fuse.EIO
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}
	toFuseAttrOut(out, res.GetAttrOut())
	return fuse.OK
}

func (fs *fileSystem) SetAttr(cancel <-chan struct{}, in *fuse.SetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

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

	if err != nil {
		log.Errorf("SetAttr: %v", err)
		return fuse.EIO
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}
	toFuseAttrOut(out, res.GetAttrOut())
	return 0
}

func (fs *fileSystem) Readlink(cancel <-chan struct{}, header *fuse.InHeader) (out []byte, code fuse.Status) {
	ctx := newContext(cancel, header)
	defer releaseContext(ctx)

	res, err := fs.client.Readlink(ctx, &pb.ReadlinkRequest{
		Header: toPbHeader(header),
	}, fs.opts...)

	if err != nil {
		log.Errorf("Access: %v", err)
		return nil, fuse.EIO
	}

	if res.Status.GetCode() != 0 {
		return nil, fuse.Status(res.Status.GetCode())
	}

	return res.GetOut(), fuse.OK
}

func (fs *fileSystem) Access(cancel <-chan struct{}, input *fuse.AccessIn) (code fuse.Status) {
	ctx := newContext(cancel, &input.InHeader)
	defer releaseContext(ctx)

	res, err := fs.client.Access(ctx, &pb.AccessRequest{
		Header:  toPbHeader(&input.InHeader),
		Mask:    input.Mask,
		Padding: input.Padding,
	}, fs.opts...)

	if err != nil {
		log.Errorf("Access: %v", err)
		return fuse.EIO
	}

	return fuse.Status(res.Status.GetCode())
}

func (fs *fileSystem) Open(cancel <-chan struct{}, in *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	res, err := fs.client.Open(ctx, &pb.OpenRequest{
		OpenIn: &pb.OpenIn{
			Header: toPbHeader(&in.InHeader),
			Flags:  in.Flags,
			Mode:   in.Mode,
		},
	}, fs.opts...)

	if err != nil {
		log.Errorf("Open: %v", err)
		return fuse.EIO
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	out.Fh = res.OpenOut.Fh
	out.OpenFlags = res.OpenOut.OpenFlags
	out.Padding = res.OpenOut.Padding
	return fuse.OK
}
