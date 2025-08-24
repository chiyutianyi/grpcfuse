# Contributing to Grpcfuse

Thank you for your interest in contributing to Grpcfuse! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style](#code-style)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

This project is governed by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.17 or later
- Git
- Make
- Docker (optional, for containerized development)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/grpcfuse.git
   cd grpcfuse
   git remote add upstream https://github.com/chiyutianyi/grpcfuse.git
   ```

## Development Setup

### Quick Setup

```bash
# Install dependencies and setup development environment
make dev-setup

# Verify setup
make test
```

### Manual Setup

```bash
# Install dependencies
go mod download
go mod tidy

# Generate protobuf code (if protoc is available)
make proto

# Generate mocks
make mock

# Run tests
make test
```

### Docker Setup

```bash
# Build and run services
docker-compose up -d

# Run tests in container
docker-compose run test-runner

# Stop services
docker-compose down
```

## Code Style

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Follow the project's linting rules
- Add comments for exported functions and types
- Use meaningful variable and function names

### File Organization

- Keep files focused and single-purpose
- Group related functionality in packages
- Use descriptive package names
- Follow Go package layout conventions

### Error Handling

- Always check errors
- Use meaningful error messages
- Wrap errors with context when appropriate
- Use custom error types for specific error conditions

### Logging

- Use structured logging with logrus
- Include relevant context in log entries
- Use appropriate log levels
- Avoid logging sensitive information

## Testing

### Test Requirements

- All new code must include tests
- Maintain test coverage above 80%
- Include unit tests, integration tests, and benchmarks
- Test error conditions and edge cases

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-all

# Generate coverage reports
make coverage

# Run benchmarks
make benchmark

# Run specific test package
go test ./grpc2fuse/...

# Run specific test function
go test -run TestFunctionName ./package
```

### Test Guidelines

- Use descriptive test names
- Test both success and failure cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Clean up test resources
- Use `t.Helper()` for helper functions

### Integration Tests

- Use the test filesystem implementation
- Test complete workflows
- Verify error handling and recovery
- Test performance characteristics

## Submitting Changes

### Commit Guidelines

- Use clear, descriptive commit messages
- Follow conventional commit format:
  ```
  type(scope): description
  
  [optional body]
  [optional footer]
  ```
- Types: feat, fix, docs, style, refactor, test, chore
- Keep commits focused and atomic

### Pull Request Process

1. Create a feature branch from `main`
2. Make your changes following the style guide
3. Add tests for new functionality
4. Update documentation as needed
5. Ensure all tests pass: `make pre-commit`
6. Submit a pull request with a clear description

### Pull Request Checklist

- [ ] Code follows style guidelines
- [ ] Tests pass and coverage is maintained
- [ ] Documentation is updated
- [ ] Commit messages are clear and descriptive
- [ ] Changes are focused and atomic
- [ ] No breaking changes (unless documented)

### Review Process

- All changes require review
- Address review comments promptly
- Maintainers may request changes
- Squash commits before merging

## Release Process

### Versioning

- Follow [Semantic Versioning](https://semver.org/)
- Major version for incompatible changes
- Minor version for new features
- Patch version for bug fixes

### Release Steps

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create and push a tag
4. Build and test release artifacts
5. Publish release notes

## Getting Help

### Questions and Discussion

- Use [GitHub Discussions](https://github.com/chiyutianyi/grpcfuse/discussions)
- Ask questions in issues
- Join community channels

### Reporting Issues

- Use the issue template
- Include reproduction steps
- Provide system information
- Include relevant logs

### Feature Requests

- Describe the use case clearly
- Explain the expected behavior
- Consider implementation complexity
- Discuss alternatives

## License

By contributing to Grpcfuse, you agree that your contributions will be licensed under the Apache License 2.0.

## Acknowledgments

Thank you to all contributors who have helped make Grpcfuse better!