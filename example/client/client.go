package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/chiyutianyi/grpcfuse/grpcfuse"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Usage: %s <mountpath> <fuseserver>")
	}

	mp := flag.Arg(0)
	fuseServer := flag.Arg(1)

	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(fuseServer, dialOpts...)
	if err != nil {
		log.Fatal(err)
	}
	cli := pb.NewRawFileSystemClient(conn)
	fs := grpcfuse.NewFileSystem(cli)

	var opt fuse.MountOptions
	opt.FsName = "GrpcFS"
	opt.Name = "grpcfs"
	opt.SingleThreaded = false
	opt.MaxBackground = 50
	opt.EnableLocks = true
	opt.IgnoreSecurityLabels = true
	opt.MaxWrite = 1 << 20
	opt.MaxReadAhead = 1 << 20
	opt.DirectMount = true
	opt.AllowOther = os.Getuid() == 0
	opt.Options = append(opt.Options, "default_permissions")
	if runtime.GOOS == "darwin" {
		opt.Options = append(opt.Options, "fssubtype=grpcfs")
		opt.Options = append(opt.Options, "volname=grpcfs")
		opt.Options = append(opt.Options, "daemon_timeout=60", "iosize=65536", "novncache")
	}

	srv, err := fuse.NewServer(fs, mp, &opt)
	if err != nil {
		log.Fatalf("new fuse server: %v", err)
	}

	go srv.Serve()

	signal.Ignore(syscall.SIGPIPE)
	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for range sigCh {
		err := srv.Unmount()
		log.Fatalf("unmount: %v", err)
	}
}
