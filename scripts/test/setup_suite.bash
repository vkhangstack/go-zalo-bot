#!/usr/bin/env bash

# Setup suite for BATS tests
# This file is sourced before running test suites

# Export test directories
export TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export SCRIPTS_DIR="$(cd "${TEST_DIR}/.." && pwd)"
export PROJECT_ROOT="$(cd "${SCRIPTS_DIR}/.." && pwd)"
export TEST_TEMP_DIR="${TEST_DIR}/tmp"

# Create temporary directory for tests
setup_test_environment() {
    mkdir -p "${TEST_TEMP_DIR}"
}

# Clean up temporary directory
teardown_test_environment() {
    if [ -d "${TEST_TEMP_DIR}" ]; then
        rm -rf "${TEST_TEMP_DIR}"
    fi
}

# Mock Git repository for testing
create_mock_git_repo() {
    local repo_dir="$1"
    mkdir -p "${repo_dir}"
    cd "${repo_dir}"
    git init
    git config user.email "test@example.com"
    git config user.name "Test User"
    echo "# Test" > README.md
    git add README.md
    git commit -m "Initial commit"
}

# Create a test Go module
create_test_go_module() {
    local module_dir="$1"
    local module_name="${2:-example.com/test}"
    
    mkdir -p "${module_dir}"
    cd "${module_dir}"
    
    cat > go.mod << EOF
module ${module_name}

go 1.20
EOF
    
    cat > main.go << EOF
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF
}
