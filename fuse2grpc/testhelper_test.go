package fuse2grpc_test

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/mock"
	"github.com/chiyutianyi/grpcfuse/pb"
)

var (
	serverSocketPath = GetTemporarySocketFileName()
)

func startTestServices(t *testing.T) (*grpc.Server, *mock.MockRawFileSystem) {
	ctl := gomock.NewController(t)

	fs := mock.NewMockRawFileSystem(ctl)

	server := NewTestGrpcServer(t, nil, nil)

	if err := os.RemoveAll(serverSocketPath); err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("unix", serverSocketPath)
	if err != nil {
		t.Fatal("failed to start server")
	}

	pb.RegisterRawFileSystemServer(server, fuse2grpc.NewServer(fs))
	reflection.Register(server)

	go server.Serve(listener)
	return server, fs
}

func newRawFileSystemClient(t *testing.T, serviceSocketPath string) (pb.RawFileSystemClient, *grpc.ClientConn) {
	connOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, _ time.Duration) (net.Conn, error) {
			return net.Dial("unix", addr)
		}),
	}
	conn, err := grpc.Dial(serviceSocketPath, connOpts...)
	if err != nil {
		t.Fatal(err)
	}

	return pb.NewRawFileSystemClient(conn), conn
}

// Context returns a cancellable context.
func Context() (context.Context, func()) {
	return context.WithCancel(context.Background())
}

// NewTestGrpcServer creates a GRPC Server for testing purposes
func NewTestGrpcServer(t *testing.T, streamInterceptors []grpc.StreamServerInterceptor, unaryInterceptors []grpc.UnaryServerInterceptor) *grpc.Server {
	logger := NewTestLogger(t)
	logrusEntry := log.NewEntry(logger).WithField("test", t.Name())
	streamInterceptors = append([]grpc.StreamServerInterceptor{grpc_logrus.StreamServerInterceptor(logrusEntry)}, streamInterceptors...)
	unaryInterceptors = append([]grpc.UnaryServerInterceptor{grpc_logrus.UnaryServerInterceptor(logrusEntry)}, unaryInterceptors...)
	return grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
	)
}

// NewTestLogger created a logrus hook which can be used with testing logs
func NewTestLogger(t *testing.T) *log.Logger {
	logger := log.New()
	logger.Out = ioutil.Discard
	return logger
}

// GetTemporarySocketFileName will return a unique, useable socket file name
func GetTemporarySocketFileName() string {
	tmpfile, err := ioutil.TempFile("", "fuse2grpc.socket")
	if err != nil {
		// No point in handling this error, panic
		panic(err)
	}

	name := tmpfile.Name()
	tmpfile.Close()
	os.Remove(name)

	return name
}
