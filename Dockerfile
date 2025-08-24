# Multi-stage build for grpcfuse
FROM golang:1.17-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    make \
    protobuf \
    protobuf-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make example

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    fuse3 \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 grpcfuse && \
    adduser -D -s /bin/sh -u 1000 -G grpcfuse grpcfuse

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/* /app/bin/

# Create mount point directory
RUN mkdir -p /mnt/grpcfuse && \
    chown -R grpcfuse:grpcfuse /mnt/grpcfuse

# Switch to non-root user
USER grpcfuse

# Expose default gRPC port
EXPOSE 8760

# Set default command
CMD ["/app/bin/loopback", "/mnt/grpcfuse"]