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

	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func (fs *fileSystem) doReadDir(
	cancel <-chan struct{},
	in *fuse.ReadIn,
	out *fuse.DirEntryList,
	reader func(ctx context.Context, in *pb.ReadDirRequest) RawFileSystem_ReadDirClient,
	funcName string,
) fuse.Status {
	var de fuse.DirEntry
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	stream := reader(ctx, &pb.ReadDirRequest{
		ReadIn: &pb.ReadIn{
			Header:    toPbHeader(&in.InHeader),
			Fh:        in.Fh,
			ReadFlags: in.ReadFlags,
			Offset:    in.Offset,
			Size:      in.Size,
		},
	})

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
