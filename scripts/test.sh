#!/usr/bin/env bash

# Test execution script for Go Zalo Bot SDK
# Runs tests with optional coverage, race detection, and verbose output

set -euo pipefail

# Source utility functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/utils.sh"

# Default values
COVERAGE=false
VERBOSE=false
RACE=false
COVERAGE_THRESHOLD=80

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -c|--coverage)
                COVERAGE=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -r|--race)
                RACE=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 2
                ;;
        esac
    done
}

# Show help message
show_help() {
    cat << EOF
Usage: $(basename "$0") [options]

Run all tests for the Go Zalo Bot SDK with optional coverage and race detection.

Options:
  -c, --coverage   Generate coverage report (HTML and terminal output)
  -v, --verbose    Enable verbose test output
  -r, --race       Enable race detector
  -h, --help       Show this help message

Examples:
  $(basename "$0")                    # Run all tests
  $(basename "$0") --coverage         # Run tests with coverage report
  $(basename "$0") --race --verbose   # Run tests with race detector and verbose output
  $(basename "$0") -c -r              # Run tests with coverage and race detector

EOF
}

# Run tests
run_tests() {
    log_info "Running tests..."
    
    # Build test command
    local test_cmd="go test"
    local test_args="./..."
    
    # Add verbose flag if requested
    if [[ "$VERBOSE" == true ]]; then
        test_cmd="$test_cmd -v"
    fi
    
    # Add race detector if requested
    if [[ "$RACE" == true ]]; then
        test_cmd="$test_cmd -race"
        log_info "Race detector enabled"
    fi
    
    # Add coverage if requested
    if [[ "$COVERAGE" == true ]]; then
        test_cmd="$test_cmd -coverprofile=coverage.out -covermode=atomic"
        log_info "Coverage reporting enabled"
    fi
    
    # Record start time
    local start_time=$(date +%s)
    
    # Run tests and capture output
    local test_output
    local test_exit_code=0
    
    if test_output=$(eval "$test_cmd $test_args" 2>&1); then
        test_exit_code=0
    else
        test_exit_code=$?
    fi
    
    # Record end time
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    # Display test output
    echo "$test_output"
    
    # Check if tests passed
    if [[ $test_exit_code -ne 0 ]]; then
        log_error "Tests failed (exit code: $test_exit_code)"
        return $test_exit_code
    fi
    
    log_success "All tests passed in ${duration}s"
    
    return 0
}

# Generate and display coverage report
generate_coverage_report() {
    if [[ "$COVERAGE" != true ]]; then
        return 0
    fi
    
    if [[ ! -f coverage.out ]]; then
        log_warning "Coverage file not found, skipping coverage report"
        return 0
    fi
    
    log_info "Generating coverage report..."
    
    # Generate HTML coverage report
    if go tool cover -html=coverage.out -o coverage.html 2>/dev/null; then
        log_success "HTML coverage report generated: coverage.html"
    else
        log_warning "Failed to generate HTML coverage report"
    fi
    
    # Calculate coverage percentage
    local coverage_output
    if coverage_output=$(go tool cover -func=coverage.out 2>/dev/null); then
        # Extract total coverage percentage
        local coverage_percent=$(echo "$coverage_output" | grep "total:" | awk '{print $3}' | sed 's/%//')
        
        if [[ -n "$coverage_percent" ]]; then
            # Display coverage summary
            echo ""
            log_info "Coverage Summary:"
            echo "$coverage_output" | tail -n 1
            echo ""
            
            # Check coverage threshold
            local coverage_int=${coverage_percent%.*}
            if [[ $coverage_int -lt $COVERAGE_THRESHOLD ]]; then
                log_error "Coverage ${coverage_percent}% is below threshold ${COVERAGE_THRESHOLD}%"
                return 4
            else
                log_success "Coverage ${coverage_percent}% meets threshold ${COVERAGE_THRESHOLD}%"
            fi
        else
            log_warning "Could not extract coverage percentage"
        fi
    else
        log_warning "Failed to generate coverage summary"
    fi
    
    return 0
}

# Display summary
display_summary() {
    echo ""
    log_info "Test Summary:"
    echo "  Coverage: $(if [[ "$COVERAGE" == true ]]; then echo "enabled"; else echo "disabled"; fi)"
    echo "  Race detector: $(if [[ "$RACE" == true ]]; then echo "enabled"; else echo "disabled"; fi)"
    echo "  Verbose: $(if [[ "$VERBOSE" == true ]]; then echo "enabled"; else echo "disabled"; fi)"
    echo ""
}

# Main execution
main() {
    log_info "Go Zalo Bot SDK - Test Execution"
    echo ""
    
    # Parse arguments
    parse_args "$@"
    
    # Check Go installation
    check_command "go" "Go is not installed. Please install Go 1.20 or higher."
    check_go_version "1.20"
    
    # Display configuration
    display_summary
    
    # Run tests
    if ! run_tests; then
        exit 5
    fi
    
    # Generate coverage report if requested
    if ! generate_coverage_report; then
        exit 4
    fi
    
    echo ""
    log_success "Test execution completed successfully!"
}

# Run main function
main "$@"
