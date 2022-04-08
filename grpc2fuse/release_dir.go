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

	"github.com/chiyutianyi/grpcfuse/pb"

	"github.com/hanwen/go-fuse/v2/fuse"
)

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
