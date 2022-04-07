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

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/chiyutianyi/grpcfuse/pkg/utils"
)

func main() {
	debug := flag.Bool("debug", false, "print debugging messages.")
	other := flag.Bool("allow-other", false, "mount with -o allowother.")
	ro := flag.Bool("ro", false, "mount read-only")
	loggerLevel := flag.String("logger-level", "info", "log level")
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Usage: %s <mountpath> <fuseserver>")
	}

	log.SetLevel(utils.GetLogLevel(*loggerLevel))
	mp := flag.Arg(0)
	fuseServer := flag.Arg(1)

	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(fuseServer, dialOpts...)
	if err != nil {
		log.Fatal(err)
	}
	cli := pb.NewRawFileSystemClient(conn)
	fs := grpc2fuse.NewFileSystem(cli)

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
	opt.AllowOther = *other
	opt.Debug = *debug
	opt.Options = append(opt.Options, "default_permissions")
	if runtime.GOOS == "darwin" {
		opt.Options = append(opt.Options, "fssubtype=grpcfs")
		opt.Options = append(opt.Options, "volname=grpcfs")
		opt.Options = append(opt.Options, "daemon_timeout=60", "iosize=65536", "novncache")
	}
	if *ro {
		opt.Options = append(opt.Options, "ro")
	}

	srv, err := fuse.NewServer(fs, mp, &opt)
	if err != nil {
		log.Fatalf("New fuse server: %v", err)
	}

	go srv.Serve()

	signal.Ignore(syscall.SIGPIPE)
	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for range sigCh {
		err := srv.Unmount()
		if err != nil {
			log.Fatalf("Unmount: %v", err)
		} else {
			log.Info("Unmounted")
			return
		}
	}
}
