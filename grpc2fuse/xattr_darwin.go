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

func (fs *fileSystem) SetXAttr(cancel <-chan struct{}, input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	ctx := newContext(cancel, &input.InHeader)
	defer releaseContext(ctx)

	res, err := fs.client.SetXAttr(ctx, &pb.SetXAttrRequest{
		Header:   toPbHeader(&input.InHeader),
		Attr:     attr,
		Data:     data,
		Size:     input.Size,
		Flags:    input.Flags,
		Position: input.Position,
		Padding:  input.Padding,
	}, fs.opts...)

	if st := dealGrpcError("SetXAttr", err); st != fuse.OK {
		return st
	}
	return fuse.Status(res.Status.Code)
}
