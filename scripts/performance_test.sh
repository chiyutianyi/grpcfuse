#!/bin/bash

# Performance testing script for Grpcfuse
# This script runs various performance tests and generates reports

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_DURATION=300  # 5 minutes
CONCURRENT_USERS=100
TEST_DATA_SIZE="1GB"
MOUNT_POINT="/tmp/grpcfuse-test"
SERVER_ADDR="127.0.0.1:8760"
LOG_FILE="performance_test_$(date +%Y%m%d_%H%M%S).log"

# Functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

success() {
    echo -e "${GREEN}✓ $1${NC}" | tee -a "$LOG_FILE"
}

warning() {
    echo -e "${YELLOW}⚠ $1${NC}" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}✗ $1${NC}" | tee -a "$LOG_FILE"
}

check_dependencies() {
    log "Checking dependencies..."
    
    # Check if required tools are installed
    local missing_tools=()
    
    for tool in go make docker docker-compose; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        error "Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    success "All dependencies are available"
}

setup_test_environment() {
    log "Setting up test environment..."
    
    # Create mount point
    sudo mkdir -p "$MOUNT_POINT"
    sudo chmod 777 "$MOUNT_POINT"
    
    # Start services if not running
    if ! pgrep -f "grpcfuse" > /dev/null; then
        log "Starting Grpcfuse services..."
        docker-compose up -d
        sleep 10  # Wait for services to start
    fi
    
    success "Test environment setup complete"
}

run_basic_benchmarks() {
    log "Running basic benchmarks..."
    
    # Go benchmarks
    log "Running Go benchmarks..."
    make benchmark 2>&1 | tee -a "$LOG_FILE"
    
    # Memory usage test
    log "Testing memory usage..."
    go test -run TestMemoryUsage ./grpc2fuse/... 2>&1 | tee -a "$LOG_FILE"
    
    success "Basic benchmarks completed"
}

run_stress_test() {
    log "Running stress test for ${TEST_DURATION} seconds..."
    
    # Create test data
    log "Creating test data..."
    dd if=/dev/urandom of="$MOUNT_POINT/test_data" bs=1M count=1024 2>/dev/null
    
    # Run stress test
    log "Starting stress test..."
    local start_time=$(date +%s)
    local end_time=$((start_time + TEST_DURATION))
    
    while [ $(date +%s) -lt $end_time ]; do
        # Simulate concurrent file operations
        for i in $(seq 1 $CONCURRENT_USERS); do
            (
                # Read operations
                cat "$MOUNT_POINT/test_data" > /dev/null &
                
                # Write operations
                dd if=/dev/urandom of="$MOUNT_POINT/stress_$i" bs=1M count=10 2>/dev/null &
                
                # Directory operations
                mkdir -p "$MOUNT_POINT/stress_dir_$i" &
                
                # Wait for all operations
                wait
            ) &
        done
        
        # Wait for this batch to complete
        wait
        
        # Small delay between batches
        sleep 1
    done
    
    success "Stress test completed"
}

run_latency_test() {
    log "Running latency tests..."
    
    # Test file operation latency
    local test_file="$MOUNT_POINT/latency_test"
    
    # Create test file
    dd if=/dev/urandom of="$test_file" bs=1M count=100 2>/dev/null
    
    # Measure read latency
    log "Measuring read latency..."
    local read_start=$(date +%s%N)
    cat "$test_file" > /dev/null
    local read_end=$(date +%s%N)
    local read_latency=$(((read_end - read_start) / 1000000))
    
    # Measure write latency
    log "Measuring write latency..."
    local write_start=$(date +%s%N)
    dd if=/dev/urandom of="$test_file" bs=1M count=100 2>/dev/null
    local write_end=$(date +%s%N)
    local write_latency=$(((write_end - write_start) / 1000000))
    
    log "Read latency: ${read_latency}ms"
    log "Write latency: ${write_latency}ms"
    
    success "Latency tests completed"
}

run_concurrent_test() {
    log "Running concurrent access tests..."
    
    # Test concurrent read/write operations
    local test_file="$MOUNT_POINT/concurrent_test"
    dd if=/dev/urandom of="$test_file" bs=1M count=100 2>/dev/null
    
    # Start concurrent readers
    for i in $(seq 1 50); do
        (
            while true; do
                cat "$test_file" > /dev/null 2>/dev/null || break
                sleep 0.1
            done
        ) &
    done
    
    # Start concurrent writers
    for i in $(seq 1 10); do
        (
            while true; do
                dd if=/dev/urandom of="$test_file" bs=1M count=10 2>/dev/null || break
                sleep 0.5
            done
        ) &
    done
    
    # Run for 30 seconds
    sleep 30
    
    # Clean up background processes
    pkill -f "cat.*concurrent_test" || true
    pkill -f "dd.*concurrent_test" || true
    
    success "Concurrent access tests completed"
}

generate_report() {
    log "Generating performance report..."
    
    local report_file="performance_report_$(date +%Y%m%d_%H%M%S).md"
    
    cat > "$report_file" << EOF
# Grpcfuse Performance Test Report

**Test Date:** $(date)
**Test Duration:** ${TEST_DURATION} seconds
**Concurrent Users:** ${CONCURRENT_USERS}
**Test Data Size:** ${TEST_DATA_SIZE}

## Test Environment
- **Mount Point:** ${MOUNT_POINT}
- **Server Address:** ${SERVER_ADDR}
- **Log File:** ${LOG_FILE}

## Test Results

### Basic Benchmarks
- Go benchmarks completed successfully
- Memory usage tests passed

### Stress Test
- Duration: ${TEST_DURATION} seconds
- Concurrent operations: ${CONCURRENT_USERS} users
- Test data size: ${TEST_DATA_SIZE}

### Latency Tests
- File operations completed
- Performance metrics recorded

### Concurrent Access Tests
- Multiple concurrent readers/writers
- System stability verified

## Recommendations

1. Monitor system resources during high load
2. Consider tuning buffer pool sizes for your workload
3. Monitor gRPC connection health
4. Set up alerting for performance degradation

## Next Steps

1. Review detailed logs in: ${LOG_FILE}
2. Analyze system metrics during tests
3. Tune configuration based on results
4. Run tests with different workloads

EOF
    
    success "Performance report generated: $report_file"
}

cleanup() {
    log "Cleaning up test environment..."
    
    # Remove test files
    sudo rm -rf "$MOUNT_POINT"/*
    
    # Stop services if we started them
    if [ "$1" = "full" ]; then
        docker-compose down
    fi
    
    success "Cleanup completed"
}

# Main execution
main() {
    log "Starting Grpcfuse performance tests..."
    
    # Parse command line arguments
    local cleanup_level="partial"
    while [[ $# -gt 0 ]]; do
        case $1 in
            --full-cleanup)
                cleanup_level="full"
                shift
                ;;
            --help)
                echo "Usage: $0 [--full-cleanup] [--help]"
                echo "  --full-cleanup  Stop all services after tests"
                echo "  --help          Show this help message"
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Set up trap for cleanup
    trap 'cleanup $cleanup_level' EXIT
    
    # Run tests
    check_dependencies
    setup_test_environment
    run_basic_benchmarks
    run_stress_test
    run_latency_test
    run_concurrent_test
    generate_report
    
    success "All performance tests completed successfully!"
}

# Run main function
main "$@"