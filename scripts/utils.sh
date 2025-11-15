#!/bin/bash

# Utility functions for build and deployment automation scripts
# This file provides common functions for logging, validation, and Git operations

set -euo pipefail

# Color codes for output
readonly COLOR_RESET='\033[0m'
readonly COLOR_RED='\033[0;31m'
readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_YELLOW='\033[0;33m'
readonly COLOR_BLUE='\033[0;34m'
readonly COLOR_CYAN='\033[0;36m'

# Logging functions

# Print informational message
# Usage: log_info "message"
log_info() {
    echo -e "${COLOR_BLUE}ℹ INFO:${COLOR_RESET} $*"
}

# Print success message
# Usage: log_success "message"
log_success() {
    echo -e "${COLOR_GREEN}✓ SUCCESS:${COLOR_RESET} $*"
}

# Print error message and exit
# Usage: log_error "message" [exit_code]
log_error() {
    local message="$1"
    local exit_code="${2:-1}"
    echo -e "${COLOR_RED}✗ ERROR:${COLOR_RESET} ${message}" >&2
    exit "${exit_code}"
}

# Print warning message
# Usage: log_warning "message"
log_warning() {
    echo -e "${COLOR_YELLOW}⚠ WARNING:${COLOR_RESET} $*"
}

# Command checking functions

# Check if a command exists
# Usage: check_command "command_name" ["installation_instructions"]
check_command() {
    local cmd="$1"
    local install_msg="${2:-Please install ${cmd} to continue.}"
    
    if ! command -v "${cmd}" &> /dev/null; then
        log_error "${cmd} is not installed. ${install_msg}" 3
    fi
}

# Check Go version meets minimum requirement
# Usage: check_go_version "minimum_version"
check_go_version() {
    local min_version="$1"
    
    check_command "go" "Please install Go from https://golang.org/dl/"
    
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    
    # Extract major and minor versions
    local go_major
    local go_minor
    go_major=$(echo "${go_version}" | cut -d. -f1)
    go_minor=$(echo "${go_version}" | cut -d. -f2)
    
    local min_major
    local min_minor
    min_major=$(echo "${min_version}" | cut -d. -f1)
    min_minor=$(echo "${min_version}" | cut -d. -f2)
    
    if [ "${go_major}" -lt "${min_major}" ] || \
       ([ "${go_major}" -eq "${min_major}" ] && [ "${go_minor}" -lt "${min_minor}" ]); then
        log_error "Go version ${min_version} or higher is required. Found: ${go_version}" 3
    fi
    
    log_info "Go version: ${go_version}"
}

# Git utility functions

# Get the latest Git tag
# Usage: get_latest_tag
get_latest_tag() {
    check_command "git"
    
    local latest_tag
    latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    
    if [ -z "${latest_tag}" ]; then
        echo ""
    else
        echo "${latest_tag}"
    fi
}

# Validate semantic version format
# Usage: validate_semver "version"
# Returns: 0 if valid, 1 if invalid
validate_semver() {
    local version="$1"
    
    # Check if version matches semantic versioning pattern: v*.*.* or *.*.*
    if [[ "${version}" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$ ]]; then
        return 0
    else
        return 1
    fi
}

# Check if a Git tag exists
# Usage: tag_exists "tag_name"
# Returns: 0 if exists, 1 if not
tag_exists() {
    local tag="$1"
    
    check_command "git"
    
    if git rev-parse "${tag}" >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Get current Git branch name
# Usage: get_current_branch
get_current_branch() {
    check_command "git"
    
    git rev-parse --abbrev-ref HEAD
}

# Check if working directory is clean
# Usage: is_git_clean
# Returns: 0 if clean, 1 if dirty
is_git_clean() {
    check_command "git"
    
    if [ -z "$(git status --porcelain)" ]; then
        return 0
    else
        return 1
    fi
}
