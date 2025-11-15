#!/usr/bin/env bash

# Test runner for script testing suite
# Runs all BATS tests for the build and deployment automation scripts

set -euo pipefail

# Get script directory
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(cd "${TEST_DIR}/.." && pwd)"

# Colors for output
readonly COLOR_RESET='\033[0m'
readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_RED='\033[0;31m'
readonly COLOR_YELLOW='\033[0;33m'
readonly COLOR_BLUE='\033[0;34m'

# Configuration
VERBOSE=false
FILTER=""

# Print usage
show_help() {
    cat << EOF
Usage: $(basename "$0") [options] [test_file]

Run BATS tests for build and deployment automation scripts.

Arguments:
  test_file        Optional: specific test file to run (e.g., test_utils.bats)

Options:
  -v, --verbose    Enable verbose output
  -h, --help       Show this help message

Examples:
  $(basename "$0")                    # Run all tests
  $(basename "$0") test_utils.bats    # Run specific test file
  $(basename "$0") --verbose          # Run with verbose output

Requirements:
  - BATS (Bash Automated Testing System) must be installed
  - Run 'npm install -g bats' or 'brew install bats-core' to install

EOF
}

# Check if BATS is installed
check_bats() {
    if ! command -v bats &> /dev/null; then
        echo -e "${COLOR_RED}ERROR: BATS is not installed${COLOR_RESET}" >&2
        echo ""
        echo "To install BATS, run one of the following:"
        echo ""
        echo "  # Using npm:"
        echo "  npm install -g bats"
        echo ""
        echo "  # Using homebrew (macOS):"
        echo "  brew install bats-core"
        echo ""
        echo "  # Using apt (Ubuntu/Debian):"
        echo "  sudo apt-get install bats"
        echo ""
        echo "  # Manual installation:"
        echo "  git clone https://github.com/bats-core/bats-core.git"
        echo "  cd bats-core"
        echo "  sudo ./install.sh /usr/local"
        echo ""
        echo "For more information, visit: https://github.com/bats-core/bats-core"
        exit 1
    fi
}

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
            *.bats)
                FILTER="$1"
                shift
                ;;
            *)
                echo -e "${COLOR_RED}ERROR: Unknown option: $1${COLOR_RESET}" >&2
                show_help
                exit 2
                ;;
        esac
    done
}

# Run tests
run_tests() {
    local test_files
    local bats_args=""
    
    if [[ "${VERBOSE}" == true ]]; then
        bats_args="--verbose-run"
    fi
    
    # Determine which tests to run
    if [[ -n "${FILTER}" ]]; then
        if [[ -f "${TEST_DIR}/${FILTER}" ]]; then
            test_files="${TEST_DIR}/${FILTER}"
        else
            echo -e "${COLOR_RED}ERROR: Test file not found: ${FILTER}${COLOR_RESET}" >&2
            exit 1
        fi
    else
        test_files="${TEST_DIR}/test_*.bats"
    fi
    
    echo -e "${COLOR_BLUE}Running BATS tests...${COLOR_RESET}"
    echo ""
    
    # Run BATS tests
    if bats ${bats_args} ${test_files}; then
        echo ""
        echo -e "${COLOR_GREEN}✓ All tests passed!${COLOR_RESET}"
        return 0
    else
        echo ""
        echo -e "${COLOR_RED}✗ Some tests failed${COLOR_RESET}"
        return 1
    fi
}

# Main execution
main() {
    parse_args "$@"
    
    # Check prerequisites
    check_bats
    
    # Change to test directory
    cd "${TEST_DIR}"
    
    # Run tests
    if run_tests; then
        exit 0
    else
        exit 1
    fi
}

# Run main function
main "$@"
