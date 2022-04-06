package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func main() {
	debug := flag.Bool("debug", false, "print debugging messages.")
	other := flag.Bool("allow-other", false, "mount with -o allowother.")
	quiet := flag.Bool("q", false, "quiet")
	ro := flag.Bool("ro", false, "mount read-only")
	flag.Parse()

	if flag.NArg() < 1 {
		logrus.Fatal("Usage: %s <ORIGINAL>")
	}

	orig := flag.Arg(0)

	l, err := net.Listen("tcp", "127.0.0.1:8760")
	if err != nil {
		logrus.Fatal(err)
	}
	s := grpc.NewServer()

	loopbackRoot, err := fs.NewLoopbackRoot(orig)

	sec := time.Second
	opts := &fs.Options{
		// These options are to be compatible with libfuse defaults,
		// making benchmarking easier.
		AttrTimeout:  &sec,
		EntryTimeout: &sec,
	}
	opts.Debug = *debug
	opts.AllowOther = *other
	if opts.AllowOther {
		// Make the kernel check file permissions for us
		opts.MountOptions.Options = append(opts.MountOptions.Options, "default_permissions")
	}
	if *ro {
		opts.MountOptions.Options = append(opts.MountOptions.Options, "ro")
	}
	// First column in "df -T": original dir
	opts.MountOptions.Options = append(opts.MountOptions.Options, "fsname="+orig)
	// Second column in "df -T" will be shown as "fuse." + Name
	opts.MountOptions.Name = "loopback"
	// Leave file permissions on "000" files as-is
	opts.NullPermissions = true
	// Enable diagnostics logging
	if !*quiet {
		opts.Logger = log.New(os.Stderr, "", 0)
	}

	rawFS := fs.NewNodeFS(loopbackRoot, opts)

	srv := fuse2grpc.NewServer(rawFS)

	pb.RegisterRawFileSystemServer(s, srv)
	go s.Serve(l)

	logrus.Infof("Listen on %s for dir %s", l.Addr(), orig)

	signal.Ignore(syscall.SIGPIPE)
	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for range sigCh {
		logrus.Fatalf("unmount: %v", err)
	}
}
