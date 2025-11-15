#!/usr/bin/env bash

# Build validation script for Go Zalo Bot SDK
# Validates compilation and dependencies

set -euo pipefail

# Source utility functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/utils.sh"

# Configuration
MIN_GO_VERSION="1.20"
VERBOSE=false

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--verbose)
                VERBOSE=true
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

Build validation script for Go Zalo Bot SDK.
Validates Go version, dependencies, and compilation.

Options:
  -v, --verbose    Enable verbose output
  -h, --help       Show this help message

Examples:
  $(basename "$0")              # Run build validation
  $(basename "$0") --verbose    # Run with verbose output

EOF
}

# Verify Go version
verify_go_version() {
    log_info "Verifying Go version..."
    
    if ! check_command go; then
        log_error "Go is not installed. Please install Go ${MIN_GO_VERSION} or higher."
    fi
    
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    
    if [[ "$VERBOSE" == true ]]; then
        log_info "Found Go version: ${go_version}"
    fi
    
    check_go_version "$MIN_GO_VERSION"
    log_success "Go version check passed (${go_version} >= ${MIN_GO_VERSION})"
}

# Verify dependencies
verify_dependencies() {
    log_info "Verifying module dependencies..."
    
    if [[ "$VERBOSE" == true ]]; then
        go mod verify
    else
        go mod verify > /dev/null 2>&1
    fi
    
    log_success "Module dependencies verified"
}

# Clean up dependencies
cleanup_dependencies() {
    log_info "Cleaning up dependencies..."
    
    if [[ "$VERBOSE" == true ]]; then
        go mod tidy -v
    else
        go mod tidy
    fi
    
    log_success "Dependencies cleaned up"
}

# Compile all packages
compile_packages() {
    log_info "Compiling all packages..."
    
    # Get all packages excluding examples directory
    local packages
    packages=$(go list ./... 2>/dev/null | grep -v '/examples/' || true)
    
    if [[ -z "$packages" ]]; then
        log_warning "No packages found to compile"
        return 0
    fi
    
    if [[ "$VERBOSE" == true ]]; then
        echo "$packages" | xargs go build -v
    else
        echo "$packages" | xargs go build
    fi
    
    log_success "All packages compiled successfully"
}

# Compile examples
compile_examples() {
    log_info "Compiling examples..."
    
    local examples_dir="examples"
    local example_count=0
    local failed_examples=()
    
    if [[ ! -d "$examples_dir" ]]; then
        log_warning "Examples directory not found, skipping example compilation"
        return 0
    fi
    
    # Find all main.go files in examples directory (e.g., examples/webhook/main.go)
    while IFS= read -r -d '' example_file; do
        local example_dir
        example_dir=$(dirname "$example_file")
        local example_name
        example_name=$(basename "$example_dir")
        
        if [[ "$VERBOSE" == true ]]; then
            log_info "  Compiling example: ${example_name}"
        fi
        
        if [[ "$VERBOSE" == true ]]; then
            if go build -o /dev/null "$example_file"; then
                example_count=$((example_count + 1))
            else
                failed_examples+=("$example_name")
            fi
        else
            if go build -o /dev/null "$example_file" 2>/dev/null; then
                example_count=$((example_count + 1))
            else
                failed_examples+=("$example_name")
            fi
        fi
    done < <(find "$examples_dir" -type f -name "main.go" -print0)
    
    # Find all .go files with main package in examples directory (e.g., examples/advanced/*.go)
    # Each file is compiled individually to avoid "main redeclared" errors
    while IFS= read -r -d '' example_file; do
        # Skip if it's already a main.go (handled above)
        if [[ "$(basename "$example_file")" == "main.go" ]]; then
            continue
        fi
        
        # Check if file contains "package main"
        if ! grep -q "^package main" "$example_file" 2>/dev/null; then
            continue
        fi
        
        local example_name
        example_name=$(basename "$example_file" .go)
        
        if [[ "$VERBOSE" == true ]]; then
            log_info "  Compiling example: ${example_name}"
        fi
        
        # Compile each example file individually to avoid conflicts
        if [[ "$VERBOSE" == true ]]; then
            if go build -o /dev/null "$example_file"; then
                example_count=$((example_count + 1))
            else
                failed_examples+=("$example_name")
            fi
        else
            if go build -o /dev/null "$example_file" 2>/dev/null; then
                example_count=$((example_count + 1))
            else
                failed_examples+=("$example_name")
            fi
        fi
    done < <(find "$examples_dir" -type f -name "*.go" ! -name "main.go" -print0)
    
    if [[ ${#failed_examples[@]} -gt 0 ]]; then
        log_error "Failed to compile examples: ${failed_examples[*]}"
    fi
    
    if [[ $example_count -eq 0 ]]; then
        log_warning "No examples found to compile"
    else
        log_success "All ${example_count} examples compiled successfully"
    fi
}

# Main execution
main() {
    local start_time
    start_time=$(date +%s)
    
    log_info "Starting build validation..."
    echo ""
    
    # Run validation steps
    verify_go_version
    verify_dependencies
    cleanup_dependencies
    compile_packages
    compile_examples
    
    # Calculate duration
    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    log_success "Build validation completed successfully in ${duration}s"
}

# Parse arguments and run
parse_args "$@"
main
