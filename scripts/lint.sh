#!/bin/bash

# Linting and static analysis script for Go Zalo Bot SDK
# Performs code quality checks including formatting, vetting, and linting

set -euo pipefail

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Source utility functions
# shellcheck source=utils.sh
source "${SCRIPT_DIR}/utils.sh"

# Configuration
VERBOSE=false
AUTO_FIX=false

# Exit codes
readonly EXIT_SUCCESS=0
readonly EXIT_LINT_FAILURE=1
readonly EXIT_INVALID_ARGS=2

# Counters for issues
TOTAL_ISSUES=0
FMT_ISSUES=0
VET_ISSUES=0
LINT_ISSUES=0
GODOC_ISSUES=0

# Show usage information
show_help() {
    cat << EOF
Usage: $(basename "$0") [options]

Perform linting and static analysis on the Go Zalo Bot SDK.

Options:
  -f, --fix        Auto-fix issues where possible (formatting)
  -v, --verbose    Enable verbose output
  -h, --help       Show this help message

Examples:
  $(basename "$0")              # Run all linting checks
  $(basename "$0") --fix        # Run checks and auto-fix formatting
  $(basename "$0") --verbose    # Run with detailed output

Exit Codes:
  0 - All checks passed
  1 - Linting issues found
  2 - Invalid arguments

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--fix)
                AUTO_FIX=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                show_help
                exit "${EXIT_SUCCESS}"
                ;;
            *)
                log_error "Unknown option: $1" "${EXIT_INVALID_ARGS}"
                ;;
        esac
    done
}

# Check if golangci-lint is installed
check_golangci_lint() {
    if ! command -v golangci-lint &> /dev/null; then
        log_warning "golangci-lint is not installed"
        echo ""
        echo "To install golangci-lint, run one of the following:"
        echo ""
        echo "  # Using go install (recommended):"
        echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        echo ""
        echo "  # Using homebrew (macOS):"
        echo "  brew install golangci-lint"
        echo ""
        echo "  # Using script (Linux/macOS):"
        echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
        echo ""
        echo "For more installation options, visit: https://golangci-lint.run/usage/install/"
        echo ""
        return 1
    fi
    return 0
}

# Check Go formatting
check_formatting() {
    log_info "Checking code formatting with go fmt..."
    
    local unformatted_files
    unformatted_files=$(gofmt -l . 2>&1 | grep -v "^vendor/" || true)
    
    if [ -n "${unformatted_files}" ]; then
        FMT_ISSUES=$(echo "${unformatted_files}" | wc -l | tr -d ' ')
        TOTAL_ISSUES=$((TOTAL_ISSUES + FMT_ISSUES))
        
        log_warning "Found ${FMT_ISSUES} file(s) with formatting issues:"
        echo "${unformatted_files}" | while read -r file; do
            echo "  - ${file}"
        done
        
        if [ "${AUTO_FIX}" = true ]; then
            log_info "Auto-fixing formatting issues..."
            echo "${unformatted_files}" | xargs -I {} gofmt -w {}
            log_success "Formatting issues fixed"
            FMT_ISSUES=0
            TOTAL_ISSUES=$((TOTAL_ISSUES - FMT_ISSUES))
        fi
    else
        log_success "All files are properly formatted"
    fi
}

# Run go vet for suspicious constructs
check_vet() {
    log_info "Running go vet for suspicious constructs..."
    
    local vet_output
    local vet_exit_code=0
    
    if [ "${VERBOSE}" = true ]; then
        vet_output=$(go vet ./... 2>&1) || vet_exit_code=$?
    else
        vet_output=$(go vet ./... 2>&1) || vet_exit_code=$?
    fi
    
    if [ ${vet_exit_code} -ne 0 ]; then
        VET_ISSUES=$(echo "${vet_output}" | grep -c "^" || echo "0")
        TOTAL_ISSUES=$((TOTAL_ISSUES + VET_ISSUES))
        
        log_warning "go vet found ${VET_ISSUES} issue(s):"
        echo "${vet_output}"
    else
        log_success "go vet found no issues"
    fi
}

# Run golangci-lint
run_golangci_lint() {
    if ! check_golangci_lint; then
        log_warning "Skipping golangci-lint checks (not installed)"
        return
    fi
    
    log_info "Running golangci-lint..."
    
    local lint_args="run"
    if [ "${VERBOSE}" = true ]; then
        lint_args="${lint_args} -v"
    fi
    
    local lint_output
    local lint_exit_code=0
    
    lint_output=$(golangci-lint ${lint_args} ./... 2>&1) || lint_exit_code=$?
    
    if [ ${lint_exit_code} -ne 0 ]; then
        # Count issues (each issue typically has a line with file:line:col)
        LINT_ISSUES=$(echo "${lint_output}" | grep -c ":[0-9]*:[0-9]*:" || echo "0")
        TOTAL_ISSUES=$((TOTAL_ISSUES + LINT_ISSUES))
        
        log_warning "golangci-lint found ${LINT_ISSUES} issue(s):"
        echo "${lint_output}"
    else
        log_success "golangci-lint found no issues"
    fi
}

# Check for missing godoc comments on exported symbols
check_godoc_comments() {
    log_info "Checking godoc comments for exported symbols..."
    
    local missing_docs=0
    local files_with_issues=""
    
    # Find all Go files except test files and vendor
    while IFS= read -r file; do
        # Skip vendor and hidden directories
        if [[ "${file}" == *"/vendor/"* ]] || [[ "${file}" == *"/."* ]]; then
            continue
        fi
        
        # Check for exported symbols without comments
        local file_issues
        file_issues=$(awk '
            /^(func|type|const|var) [A-Z]/ {
                # Check if previous line is a comment
                if (prev !~ /^\/\//) {
                    print FILENAME ":" NR ": exported " $1 " " $2 " should have comment"
                    count++
                }
            }
            { prev = $0 }
            END { if (count > 0) exit 1 }
        ' "${file}" 2>&1 || true)
        
        if [ -n "${file_issues}" ]; then
            if [ -z "${files_with_issues}" ]; then
                files_with_issues="${file_issues}"
            else
                files_with_issues="${files_with_issues}"$'\n'"${file_issues}"
            fi
            missing_docs=$((missing_docs + $(echo "${file_issues}" | wc -l | tr -d ' ')))
        fi
    done < <(find "${PROJECT_ROOT}" -name "*.go" -not -name "*_test.go" -type f)
    
    if [ ${missing_docs} -gt 0 ]; then
        GODOC_ISSUES=${missing_docs}
        TOTAL_ISSUES=$((TOTAL_ISSUES + GODOC_ISSUES))
        
        log_warning "Found ${GODOC_ISSUES} exported symbol(s) without godoc comments:"
        echo "${files_with_issues}"
    else
        log_success "All exported symbols have godoc comments"
    fi
}

# Main execution
main() {
    parse_args "$@"
    
    cd "${PROJECT_ROOT}"
    
    log_info "Starting linting and static analysis..."
    echo ""
    
    # Check Go installation
    check_go_version "1.20"
    echo ""
    
    # Run all checks
    check_formatting
    echo ""
    
    check_vet
    echo ""
    
    run_golangci_lint
    echo ""
    
    check_godoc_comments
    echo ""
    
    # Summary
    echo "========================================"
    echo "Linting Summary"
    echo "========================================"
    echo "Formatting issues:  ${FMT_ISSUES}"
    echo "go vet issues:      ${VET_ISSUES}"
    echo "golangci-lint:      ${LINT_ISSUES}"
    echo "Missing godoc:      ${GODOC_ISSUES}"
    echo "----------------------------------------"
    echo "Total issues:       ${TOTAL_ISSUES}"
    echo "========================================"
    echo ""
    
    if [ ${TOTAL_ISSUES} -eq 0 ]; then
        log_success "All linting checks passed!"
        exit "${EXIT_SUCCESS}"
    else
        log_error "Linting found ${TOTAL_ISSUES} issue(s)" "${EXIT_LINT_FAILURE}"
    fi
}

# Run main function
main "$@"
