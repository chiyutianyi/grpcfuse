# Grpcfuse

[![CI](https://github.com/chiyutianyi/grpcfuse/actions/workflows/ci.yml/badge.svg)](https://github.com/chiyutianyi/grpcfuse/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chiyutianyi/grpcfuse)](https://goreportcard.com/report/github.com/chiyutianyi/grpcfuse)
[![GoDoc](https://godoc.org/github.com/chiyutianyi/grpcfuse?status.svg)](https://godoc.org/github.com/chiyutianyi/grpcfuse)

高性能远程文件系统，基于gRPC和FUSE实现。服务端和客户端均使用纯Go语言开发。

## 🚀 主要特性

- **高性能**: 优化的缓冲池管理，减少GC压力
- **可靠性**: 内置重试机制和超时控制
- **可配置**: 灵活的配置系统，支持自定义超时和性能参数
- **全面测试**: 高测试覆盖率，包含单元测试、集成测试和基准测试
- **跨平台**: 支持Linux和macOS

## 🏗️ 架构

Grpcfuse 由两个主要组件构成：

1. **GRPC服务器** (`fuse2grpc`): 将FUSE操作转换为gRPC调用
2. **GRPC客户端** (`grpc2fuse`): 将gRPC调用转换为FUSE操作

两个组件都遵循 [github.com/hanwen/go-fuse/fuse#RawFileSystem](https://pkg.go.dev/github.com/hanwen/go-fuse/fuse#RawFileSystem) 接口，支持多种服务端实现（如 [pathfs#FileSystem](https://pkg.go.dev/github.com/hanwen/go-fuse/fuse/pathfs#FileSystem)、[nodefs#Node](https://pkg.go.dev/github.com/hanwen/go-fuse/fuse/nodefs#Node) 或推荐的 [fs](https://pkg.go.dev/github.com/hanwen/go-fuse/v2/fs)）。

## ⚡ 性能优化

- **优化的缓冲池**: 使用读写锁和双重检查锁定，减少锁竞争
- **智能重试**: 仅对可重试的错误进行重试（如网络不可用、资源耗尽）
- **配置化超时**: 可自定义的超时设置，避免长时间阻塞
- **内存管理**: 显式内存管理减少GC开销

## 🚀 快速开始

### 基本使用

1. **构建项目**:
```bash
make all
```

2. **启动服务器**:
```bash
# 启动环回服务器，将 /some/directory 作为远程文件系统
./bin/loopback /some/directory
```

3. **启动客户端**:
```bash
# 将远程文件系统挂载到 /tmp/mountpoint
./bin/client /tmp/mountpoint 127.0.0.1:8760
```

### 高级配置

可以通过编程方式创建带自定义配置的文件系统：

```go
import "github.com/chiyutianyi/grpcfuse/grpc2fuse"

config := &grpc2fuse.Config{
    DefaultTimeout: time.Second * 30,
    ForgetTimeout:  time.Second * 10,
    MaxRetries:     5,
    RetryInterval:  time.Millisecond * 200,
    MaxReadAhead:   2 << 20, // 2MB
    MaxWriteSize:   2 << 20, // 2MB
}

fs := grpc2fuse.NewFileSystemWithConfig(client, config)
```

## 📊 性能指标

经过优化后的性能表现：

- **缓冲池**: >100K 操作/秒，<0.1 分配/操作
- **内存效率**: 显著减少GC压力
- **并发性能**: 优化的读写锁减少锁竞争
- **网络弹性**: 智能重试机制提高可靠性

## 🧪 测试

运行所有测试：
```bash
make test
```

运行基准测试：
```bash
go test -bench=. ./...
```

运行集成测试：
```bash
go test -tags=integration ./...
```

## Bugs

Yes, probably.  Report them through
https://github.com/chiyutianyi/grpcfuse/issues

## Disclaimer

This is not an official Alibaba product.

## License

This library is distributed under Apache License 2.0, see [LICENSE](LICENSE)
