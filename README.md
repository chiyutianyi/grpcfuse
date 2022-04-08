# Grpcfuse

Remote filesystem based on grpc and fuse. The server and client were implemented with pure go.

Grpcfuse consists of two parts:
1. GRPC Server
2. GRPC Client

They all follow [github.com/hanwen/go-fuse/fuse#RawFileSystem](https://pkg.go.dev/github.com/hanwen/go-fuse/fuse#RawFileSystem), so you can choose from multiple server-side implementations (e.g. [pathfs#FileSystem](https://pkg.go.dev/github.com/hanwen/go-fuse/fuse/pathfs#FileSystem), [nodefs#Node](https://pkg.go.dev/github.com/hanwen/go-fuse/fuse/nodefs#Node) or sugguested [fs](https://pkg.go.dev/github.com/hanwen/go-fuse/v2/fs) )and convert to RawFileSystem.

## Examples

- `example/client/client.go` contains a grpc client filesystem. A binary to run is in example/loopback/. For example
```
example/client/client /tmp/mountpoint 127.0.0.1:8760
```
- `example/loopback/server.go` contains a grpc server which mounts another piece of the filesystem. Functionally, it is similar to a symlink. A binary to run is in example/loopback/ . For example
```
example/loopback/loopback /some/other/directory
```

## Bugs

Yes, probably.  Report them through
https://github.com/chiyutianyi/grpcfuse/issues

## Disclaimer

This is not an official Alibaba product.

## License

This library is distributed under Apache License 2.0, see [LICENSE](LICENSE)
