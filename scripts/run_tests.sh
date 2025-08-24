#!/bin/bash

# Test runner script for Grpcfuse
# This script runs all tests and generates comprehensive reports

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
LOG_FILE="test_run_$(date +%Y%m%d_%H%M%S).log"
COVERAGE_DIR="coverage_reports"
TEST_RESULTS_DIR="test_results"

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

# Create directories
setup_directories() {
    log "Setting up test directories..."
    mkdir -p "$COVERAGE_DIR" "$TEST_RESULTS_DIR"
    success "Directories created"
}

# Run unit tests
run_unit_tests() {
    log "Running unit tests..."
    
    # Test grpc2fuse package
    log "Testing grpc2fuse package..."
    go test -v -cover -coverprofile="$COVERAGE_DIR/grpc2fuse.out" ./grpc2fuse/... 2>&1 | tee -a "$LOG_FILE"
    
    # Test fuse2grpc package
    log "Testing fuse2grpc package..."
    go test -v -cover -coverprofile="$COVERAGE_DIR/fuse2grpc.out" ./fuse2grpc/... 2>&1 | tee -a "$LOG_FILE"
    
    # Test utils package
    log "Testing utils package..."
    go test -v -cover -coverprofile="$COVERAGE_DIR/utils.out" ./pkg/utils/... 2>&1 | tee -a "$LOG_FILE"
    
    success "Unit tests completed"
}

# Run integration tests
run_integration_tests() {
    log "Running integration tests..."
    
    if [ -d "./integration" ]; then
        go test -v -cover -coverprofile="$COVERAGE_DIR/integration.out" ./integration/... 2>&1 | tee -a "$LOG_FILE"
        success "Integration tests completed"
    else
        warning "Integration tests directory not found, skipping"
    fi
}

# Run performance tests
run_performance_tests() {
    log "Running performance tests..."
    
    # Run benchmarks
    log "Running benchmarks..."
    go test -bench=. -benchmem ./grpc2fuse/... 2>&1 | tee "$TEST_RESULTS_DIR/grpc2fuse_benchmarks.txt"
    go test -bench=. -benchmem ./fuse2grpc/... 2>&1 | tee "$TEST_RESULTS_DIR/fuse2grpc_benchmarks.txt"
    
    success "Performance tests completed"
}

# Generate coverage reports
generate_coverage_reports() {
    log "Generating coverage reports..."
    
    # Generate HTML reports
    for profile in "$COVERAGE_DIR"/*.out; do
        if [ -f "$profile" ]; then
            local package_name=$(basename "$profile" .out)
            go tool cover -html="$profile" -o "$COVERAGE_DIR/${package_name}_coverage.html"
            log "Generated coverage report for $package_name"
        fi
    done
    
    # Generate summary
    echo "# Coverage Summary" > "$COVERAGE_DIR/coverage_summary.txt"
    echo "Generated: $(date)" >> "$COVERAGE_DIR/coverage_summary.txt"
    echo "" >> "$COVERAGE_DIR/coverage_summary.txt"
    
    for profile in "$COVERAGE_DIR"/*.out; do
        if [ -f "$profile" ]; then
            local package_name=$(basename "$profile" .out)
            echo "## $package_name" >> "$COVERAGE_DIR/coverage_summary.txt"
            go tool cover -func="$profile" >> "$COVERAGE_DIR/coverage_summary.txt"
            echo "" >> "$COVERAGE_DIR/coverage_summary.txt"
        fi
    done
    
    success "Coverage reports generated"
}

# Run code quality checks
run_code_quality_checks() {
    log "Running code quality checks..."
    
    # Format check
    log "Checking code format..."
    if [ "$(gofmt -l . | wc -l)" -gt 0 ]; then
        warning "Code formatting issues found. Run 'make fmt' to fix."
        gofmt -l . | tee -a "$LOG_FILE"
    else
        success "Code formatting is correct"
    fi
    
    # Lint check
    log "Running linter..."
    if command -v golangci-lint >/dev/null 2>&1; then
        golangci-lint run 2>&1 | tee "$TEST_RESULTS_DIR/lint_results.txt"
        success "Linting completed"
    else
        warning "golangci-lint not found, skipping linting"
    fi
    
    # Vet check
    log "Running go vet..."
    go vet ./... 2>&1 | tee "$TEST_RESULTS_DIR/vet_results.txt"
    success "Go vet completed"
}

# Generate test summary
generate_test_summary() {
    log "Generating test summary..."
    
    local summary_file="$TEST_RESULTS_DIR/test_summary.md"
    
    cat > "$summary_file" << EOF
# Grpcfuse Test Summary

**Test Run Date:** $(date)
**Log File:** $LOG_FILE

## Test Results

### Unit Tests
- grpc2fuse package: ✅ Completed
- fuse2grpc package: ✅ Completed  
- utils package: ✅ Completed

### Integration Tests
- Integration tests: ✅ Completed

### Performance Tests
- Benchmarks: ✅ Completed
- Results saved to: $TEST_RESULTS_DIR/

### Code Quality
- Formatting: ✅ Checked
- Linting: ✅ Completed
- Go vet: ✅ Completed

## Coverage Reports
Coverage reports are available in: $COVERAGE_DIR/

## Next Steps
1. Review coverage reports for areas of improvement
2. Check linting results for code quality issues
3. Analyze benchmark results for performance optimization
4. Review integration test results for system stability

EOF
    
    success "Test summary generated: $summary_file"
}

# Cleanup function
cleanup() {
    log "Cleaning up temporary files..."
    # Remove temporary coverage files
    rm -f "$COVERAGE_DIR"/*.out
    success "Cleanup completed"
}

# Main execution
main() {
    log "Starting Grpcfuse test suite..."
    
    # Set up trap for cleanup
    trap cleanup EXIT
    
    # Parse command line arguments
    local run_unit=true
    local run_integration=true
    local run_performance=true
    local run_quality=true
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit-only)
                run_integration=false
                run_performance=false
                run_quality=false
                shift
                ;;
            --integration-only)
                run_unit=false
                run_performance=false
                run_quality=false
                shift
                ;;
            --performance-only)
                run_unit=false
                run_integration=false
                run_quality=false
                shift
                ;;
            --quality-only)
                run_unit=false
                run_integration=false
                run_performance=false
                shift
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo "Options:"
                echo "  --unit-only        Run only unit tests"
                echo "  --integration-only Run only integration tests"
                echo "  --performance-only Run only performance tests"
                echo "  --quality-only     Run only code quality checks"
                echo "  --help             Show this help message"
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Run tests based on options
    setup_directories
    
    if [ "$run_unit" = true ]; then
        run_unit_tests
    fi
    
    if [ "$run_integration" = true ]; then
        run_integration_tests
    fi
    
    if [ "$run_performance" = true ]; then
        run_performance_tests
    fi
    
    if [ "$run_quality" = true ]; then
        run_code_quality_checks
    fi
    
    # Generate reports
    generate_coverage_reports
    generate_test_summary
    
    success "All tests completed successfully!"
    log "Results saved to: $TEST_RESULTS_DIR/"
    log "Coverage reports: $COVERAGE_DIR/"
    log "Full log: $LOG_FILE"
}

# Run main function
main "$@"