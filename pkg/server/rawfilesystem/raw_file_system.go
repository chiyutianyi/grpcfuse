package rawfilesystem

import (
	"context"

	"github.com/hanwen/go-fuse/v2/fuse"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func (s *server) String(context.Context, *pb.StringRequest) (*pb.StringResponse, error) {
	return &pb.StringResponse{Value: s.fs.String()}, nil
}
func (s *server) Lookup(ctx context.Context, req *pb.LookupRequest) (*pb.LookupResponse, error) {
	var out fuse.EntryOut

	ch := newCancel(ctx)
	defer releaseCancel(ch)

	status := s.fs.Lookup(ch, toFuseInHeader(req.Header), req.Name, &out)
	if status != fuse.OK {
		return &pb.LookupResponse{Status: &pb.Status{Code: int32(status)}}, nil
	}
	return &pb.LookupResponse{
		EntryOut: &pb.EntryOut{
			NodeId:     out.NodeId,
			Generation: out.Generation,
			Attr:       toPbAttr(&out.Attr),
		},
		Status: &pb.Status{Code: 0},
	}, nil
}
func (s *server) Forget(ctx context.Context, req *pb.ForgetRequest) (*emptypb.Empty, error) {
	s.fs.Forget(req.Nodeid, req.Nlookup)
	return &emptypb.Empty{}, nil
}

func (s *server) GetAttr(context.Context, *pb.GetAttrRequest) (*pb.GetAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAttr not implemented")
}
func (s *server) SetAttr(context.Context, *pb.SetAttrRequest) (*pb.SetAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAttr not implemented")
}
func (s *server) Mknod(context.Context, *pb.MknodRequest) (*pb.MknodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mknod not implemented")
}
func (s *server) Mkdir(context.Context, *pb.MkdirRequest) (*pb.MkdirResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mkdir not implemented")
}
func (s *server) Unlink(context.Context, *pb.UnlinkRequest) (*pb.UnlinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlink not implemented")
}
func (s *server) Rmdir(context.Context, *pb.RmdirRequest) (*pb.RmdirResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rmdir not implemented")
}
func (s *server) Rename(context.Context, *pb.RenameRequest) (*pb.RenameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rename not implemented")
}
func (s *server) Link(context.Context, *pb.LinkRequest) (*pb.LinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Link not implemented")
}
func (s *server) Symlink(context.Context, *pb.SymlinkRequest) (*pb.SymlinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Symlink not implemented")
}
func (s *server) Readlink(context.Context, *pb.ReadlinkRequest) (*pb.ReadlinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Readlink not implemented")
}
func (s *server) Access(context.Context, *pb.AccessRequest) (*pb.AccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Access not implemented")
}
func (s *server) GetXAttr(context.Context, *pb.GetXAttrRequest) (*pb.GetXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetXAttr not implemented")
}
func (s *server) ListXAttr(context.Context, *pb.ListXAttrRequest) (*pb.ListXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListXAttr not implemented")
}
func (s *server) SetXAttr(context.Context, *pb.SetXAttrRequest) (*pb.SetXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetXAttr not implemented")
}
func (s *server) RemoveXAttr(context.Context, *pb.RemoveXAttrRequest) (*pb.RemoveXAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveXAttr not implemented")
}
func (s *server) Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (s *server) Open(context.Context, *pb.OpenRequest) (*pb.OpenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Open not implemented")
}
func (s *server) Read(req *pb.ReadRequest, stream pb.RawFileSystem_ReadServer) error {
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (s *server) LSeek(context.Context, *pb.LSeekRequest) (*pb.LSeekResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LSeek not implemented")
}
func (s *server) GetLk(context.Context, *pb.LkRequest) (*pb.GetLkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLk not implemented")
}
func (s *server) SetLk(context.Context, *pb.LkRequest) (*pb.SetLkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLk not implemented")
}
func (s *server) SetLkw(context.Context, *pb.LkRequest) (*pb.SetLkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLkw not implemented")
}
func (s *server) Release(context.Context, *pb.ReleaseRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Release not implemented")
}
func (s *server) Write(context.Context, *pb.WriteRequest) (*pb.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (s *server) CopyFileRange(context.Context, *pb.CopyFileRangeRequest) (*pb.CopyFileRangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CopyFileRange not implemented")
}
func (s *server) Flush(context.Context, *pb.FlushRequest) (*pb.FlushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
}
func (s *server) Fsync(context.Context, *pb.FsyncRequest) (*pb.FsyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fsync not implemented")
}
func (s *server) Fallocate(context.Context, *pb.FallocateRequest) (*pb.FallocateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fallocate not implemented")
}
func (s *server) OpenDir(context.Context, *pb.OpenDirRequest) (*pb.OpenDirResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenDir not implemented")
}
func (s *server) ReadDir(*pb.ReadDirRequest, pb.RawFileSystem_ReadDirServer) error {
	return status.Errorf(codes.Unimplemented, "method ReadDir not implemented")
}
func (s *server) ReadDirPlus(*pb.ReadDirRequest, pb.RawFileSystem_ReadDirPlusServer) error {
	return status.Errorf(codes.Unimplemented, "method ReadDirPlus not implemented")
}
func (s *server) ReleaseDir(context.Context, *pb.ReleaseRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseDir not implemented")
}
func (s *server) FsyncDir(context.Context, *pb.FsyncRequest) (*pb.FsyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FsyncDir not implemented")
}
func (s *server) StatFs(ctx context.Context, req *pb.StatfsRequest) (*pb.StatfsResponse, error) {
	var out fuse.StatfsOut
	ch := newCancel(ctx)
	defer releaseCancel(ch)
	status := s.fs.StatFs(ch, toFuseInHeader(req.Input), &out)
	if status != fuse.OK {
		return &pb.StatfsResponse{Status: &pb.Status{Code: int32(status)}}, nil
	}
	return &pb.StatfsResponse{
		Blocks: out.Blocks,
		Status: &pb.Status{Code: 0},
	}, nil
}
