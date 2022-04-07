package grpc2fuse

import (
	"context"
	"io"

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

func (fs *fileSystem) Release(cancel <-chan struct{}, in *fuse.ReleaseIn) {
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	if _, err := fs.client.Release(ctx, &pb.ReleaseRequest{
		Header:       toPbHeader(&in.InHeader),
		Fh:           in.Fh,
		Flags:        in.Flags,
		ReleaseFlags: in.ReleaseFlags,
		LockOwner:    in.LockOwner,
	}, fs.opts...); err != nil {
		log.Errorf("Release: %v", err)
	}
}

func (fs *fileSystem) OpenDir(cancel <-chan struct{}, in *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	res, err := fs.client.OpenDir(ctx, &pb.OpenDirRequest{
		OpenIn: &pb.OpenIn{
			Header: toPbHeader(&in.InHeader),
			Flags:  in.Flags,
			Mode:   in.Mode,
		},
	}, fs.opts...)

	if err != nil {
		log.Errorf("OpenDir: %v", err)
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

func (fs *fileSystem) doReadDir(
	cancel <-chan struct{},
	in *fuse.ReadIn,
	out *fuse.DirEntryList,
	reader func(ctx context.Context, in *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirClient, error),
	funcName string,
) fuse.Status {
	var de fuse.DirEntry
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	stream, err := reader(ctx, &pb.ReadDirRequest{
		ReadIn: &pb.ReadIn{
			Header:    toPbHeader(&in.InHeader),
			Fh:        in.Fh,
			ReadFlags: in.ReadFlags,
			Offset:    in.Offset,
			Size:      in.Size,
		},
	}, fs.opts...)

	if err != nil {
		log.Errorf("%s: %v", funcName, err)
		return fuse.EIO
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("%s: %v", funcName, err)
			return fuse.EIO
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
	return fs.doReadDir(cancel, in, out, fs.client.ReadDir, "ReadDir")
}

func (fs *fileSystem) ReadDirPlus(cancel <-chan struct{}, in *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	var de fuse.DirEntry
	ctx := newContext(cancel, &in.InHeader)
	defer releaseContext(ctx)

	stream, err := fs.client.ReadDirPlus(ctx, &pb.ReadDirRequest{
		ReadIn: &pb.ReadIn{
			Header:    toPbHeader(&in.InHeader),
			Fh:        in.Fh,
			ReadFlags: in.ReadFlags,
			Offset:    in.Offset,
			Size:      in.Size,
		},
	}, fs.opts...)

	if err != nil {
		log.Errorf("ReadDirPlus: %v", err)
		return fuse.EIO
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("ReadDirPlus: %v", err)
			return fuse.EIO
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

func (fs *fileSystem) ReleaseDir(in *fuse.ReleaseIn) {
	if _, err := fs.client.ReleaseDir(context.TODO(), &pb.ReleaseRequest{
		Header:       toPbHeader(&in.InHeader),
		Fh:           in.Fh,
		Flags:        in.Flags,
		ReleaseFlags: in.ReleaseFlags,
		LockOwner:    in.LockOwner,
	}, fs.opts...); err != nil {
		log.Errorf("ReleaseDir: %v", err)
	}
}

func (fs *fileSystem) StatFs(cancel <-chan struct{}, in *fuse.InHeader, out *fuse.StatfsOut) (code fuse.Status) {
	ctx := newContext(cancel, in)
	defer releaseContext(ctx)

	res, err := fs.client.StatFs(ctx, &pb.StatfsRequest{
		Input: toPbHeader(in),
	}, fs.opts...)

	if err != nil {
		log.Errorf("SetAttr: %v", err)
		return fuse.EIO
	}

	if res.Status.GetCode() != 0 {
		return fuse.Status(res.Status.GetCode())
	}

	out.Blocks = res.Blocks
	out.Bfree = res.Bfree
	out.Bavail = res.Bavail
	out.Files = res.Files
	out.Ffree = res.Ffree
	out.Bsize = res.Bsize
	out.NameLen = res.NameLen
	out.Frsize = res.Frsize
	out.Padding = res.Padding
	//TODO out.Spare = res.Spare
	return fuse.OK
}
