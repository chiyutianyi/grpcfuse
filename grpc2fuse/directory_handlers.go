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
	"io"

	"github.com/chiyutianyi/grpcfuse/pb"
	log "github.com/sirupsen/logrus"

	"github.com/hanwen/go-fuse/v2/fuse"
)

func (fs *fileSystem) OpenDir(cancel <-chan struct{}, in *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.OpenDir(ctx, &pb.OpenDirRequest{
		OpenIn: &pb.OpenIn{
			Header: toPbHeader(&in.InHeader),
			Flags:  in.Flags,
			Mode:   in.Mode,
		},
	}, fs.opts...)

	if st := dealGrpcError("OpenDir", err); st != fuse.OK {
		return st
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	out.Fh = res.OpenOut.Fh
	out.OpenFlags = res.OpenOut.OpenFlags
	out.Padding = res.OpenOut.Padding
	return fuse.OK
}

func (fs *fileSystem) doReadDir(
	cancel <-chan struct{},
	in *fuse.ReadIn,
	out *fuse.DirEntryList,
	reader func(ctx context.Context, in *pb.ReadDirRequest) RawFileSystem_ReadDirClient,
	funcName string,
) fuse.Status {
	var de fuse.DirEntry
	ctx := newContext(cancel)

	stream := reader(ctx, &pb.ReadDirRequest{ReadIn: toPbReadIn(in)})

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if st := dealGrpcError(funcName, err); st != fuse.OK {
			return st
		}
		if res.Status.GetCode() != 0 {
			return fuse.Status(res.Status.GetCode())
		}
		for _, e := range res.Entries {
			de.Ino = e.Ino
			de.Name = string(e.Name)
			de.Mode = e.Mode
			if !out.AddDirEntry(de) {
				break
			}
		}
	}
	return fuse.OK
}

func (fs *fileSystem) ReadDir(cancel <-chan struct{}, in *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	var (
		stream pb.RawFileSystem_ReadDirClient
		err    error
	)

	reader := func(ctx context.Context, in *pb.ReadDirRequest) RawFileSystem_ReadDirClient {
		stream, err = fs.client.ReadDir(ctx, in, fs.opts...)
		return stream
	}

	if err != nil {
		log.Errorf("ReadDir: %v", err)
		return fuse.EIO
	}

	return fs.doReadDir(cancel, in, out, reader, "ReadDir")
}

func (fs *fileSystem) ReadDirPlus(cancel <-chan struct{}, in *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	var (
		stream pb.RawFileSystem_ReadDirPlusClient
		err    error
	)

	reader := func(ctx context.Context, in *pb.ReadDirRequest) RawFileSystem_ReadDirClient {
		stream, err = fs.client.ReadDirPlus(ctx, in, fs.opts...)
		return stream
	}

	if err != nil {
		log.Errorf("ReadDirPlus: %v", err)
		return fuse.EIO
	}

	return fs.doReadDir(cancel, in, out, reader, "ReadDirPlus")
}

func (fs *fileSystem) ReleaseDir(in *fuse.ReleaseIn) {
	if _, err := fs.client.ReleaseDir(context.TODO(), &pb.ReleaseRequest{
		Header:       toPbHeader(&in.InHeader),
		Fh:           in.Fh,
		Flags:        in.Flags,
		ReleaseFlags: in.ReleaseFlags,
		LockOwner:    in.LockOwner,
	}, fs.opts...); err != nil {
		dealGrpcError("ReleaseDir", err)
	}
}

func (fs *fileSystem) FsyncDir(cancel <-chan struct{}, input *fuse.FsyncIn) (code fuse.Status) {
	ctx := newContext(cancel)

	res, err := fs.client.FsyncDir(ctx, &pb.FsyncRequest{
		Header:     toPbHeader(&input.InHeader),
		Fh:         input.Fh,
		FsyncFlags: input.FsyncFlags,
		Padding:    input.Padding,
	}, fs.opts...)

	if st := dealGrpcError("FsyncDir", err); st != fuse.OK {
		return st
	}

	return fuse.Status(res.Status.GetCode())
}
