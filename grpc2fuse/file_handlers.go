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
	"io"

	"github.com/chiyutianyi/grpcfuse/pb"

	"github.com/hanwen/go-fuse/v2/fuse"
)

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

	if st := dealGrpcError("Open", err); st != fuse.OK {
		return st
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	toFuseOpenOut(out, res.OpenOut)
	return fuse.OK
}

func (fs *fileSystem) Read(cancel <-chan struct{}, input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	ctx := newContext(cancel, &input.InHeader)
	defer releaseContext(ctx)

	stream, err := fs.client.Read(ctx, &pb.ReadRequest{ReadIn: toPbReadIn(input)}, fs.opts...)

	if st := dealGrpcError("Read", err); st != fuse.OK {
		return nil, st
	}

	var rs []byte

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if st := dealGrpcError("Read", err); st != fuse.OK {
			return nil, st
		}
		if res.Status.GetCode() != 0 {
			return nil, fuse.Status(res.Status.GetCode())
		}

		rs = append(rs, res.Buffer...)
	}

	return fuse.ReadResultData(rs), fuse.OK
}

func (fs *fileSystem) Lseek(cancel <-chan struct{}, in *fuse.LseekIn, out *fuse.LseekOut) fuse.Status {
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	res, err := fs.client.Lseek(ctx,
		&pb.LseekRequest{
			Header:  toPbHeader(&in.InHeader),
			Fh:      in.Fh,
			Offset:  in.Offset,
			Whence:  in.Whence,
			Padding: in.Padding,
		}, fs.opts...)

	if st := dealGrpcError("Lseek", err); st != fuse.OK {
		return st
	}
	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}
	out.Offset = res.Offset
	return fuse.OK
}
