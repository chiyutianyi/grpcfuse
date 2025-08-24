# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive test suite with unit, integration, and performance tests
- Enhanced error handling with gRPC status code mapping
- Improved logging with structured logrus integration
- Docker support for containerized development and testing
- Performance benchmarks and stress tests
- Code quality tools (golangci-lint, go vet, go fmt)
- Enhanced Makefile with multiple build targets
- Integration test framework with mock filesystem
- Memory leak detection tests
- Concurrent operation safety tests

### Changed
- Refactored FileSystem struct with better encapsulation
- Improved Server implementation with functional options pattern
- Enhanced error handling and logging throughout the codebase
- Updated build system with better dependency management
- Improved documentation and code comments

### Fixed
- Better error handling for gRPC operations
- Improved thread safety in server operations
- Enhanced buffer pool management
- Better cleanup of resources

## [0.1.0] - 2022-01-01

### Added
- Initial release of Grpcfuse
- Basic FUSE to gRPC server implementation
- Basic gRPC to FUSE client implementation
- Protocol buffer definitions for filesystem operations
- Example client and loopback server
- Basic test coverage

### Features
- File operations (create, read, write, delete)
- Directory operations (create, list, remove)
- Extended attributes support
- File locking support
- Basic error handling

## [0.0.1] - 2022-01-01

### Added
- Project initialization
- Basic project structure
- License and documentation setup