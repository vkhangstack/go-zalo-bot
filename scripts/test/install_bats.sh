#!/usr/bin/env bash

# BATS installation script
# Installs BATS (Bash Automated Testing System) for script testing

set -euo pipefail

# Colors for output
readonly COLOR_RESET='\033[0m'
readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_RED='\033[0;31m'
readonly COLOR_BLUE='\033[0;34m'
readonly COLOR_YELLOW='\033[0;33m'

# Print colored message
log_info() {
    echo -e "${COLOR_BLUE}ℹ INFO:${COLOR_RESET} $*"
}

log_success() {
    echo -e "${COLOR_GREEN}✓ SUCCESS:${COLOR_RESET} $*"
}

log_error() {
    echo -e "${COLOR_RED}✗ ERROR:${COLOR_RESET} $*" >&2
}

log_warning() {
    echo -e "${COLOR_YELLOW}⚠ WARNING:${COLOR_RESET} $*"
}

# Check if BATS is already installed
check_bats_installed() {
    if command -v bats &> /dev/null; then
        local version
        version=$(bats --version | head -n 1)
        log_success "BATS is already installed: ${version}"
        return 0
    else
        return 1
    fi
}

# Detect operating system
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        echo "windows"
    else
        echo "unknown"
    fi
}

# Install BATS on Linux
install_linux() {
    log_info "Installing BATS on Linux..."
    
    # Check if apt is available
    if command -v apt-get &> /dev/null; then
        log_info "Using apt-get to install BATS..."
        sudo apt-get update
        sudo apt-get install -y bats
        return 0
    fi
    
    # Check if yum is available
    if command -v yum &> /dev/null; then
        log_info "Using yum to install BATS..."
        sudo yum install -y bats
        return 0
    fi
    
    # Fallback to manual installation
    log_warning "Package manager not found, installing manually..."
    install_manual
}

# Install BATS on macOS
install_macos() {
    log_info "Installing BATS on macOS..."
    
    # Check if Homebrew is available
    if command -v brew &> /dev/null; then
        log_info "Using Homebrew to install BATS..."
        brew install bats-core
        return 0
    else
        log_error "Homebrew is not installed. Please install Homebrew first:"
        echo "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
        return 1
    fi
}

# Manual installation
install_manual() {
    log_info "Installing BATS manually..."
    
    local install_dir="/usr/local"
    local temp_dir
    temp_dir=$(mktemp -d)
    
    log_info "Cloning BATS repository..."
    if ! git clone https://github.com/bats-core/bats-core.git "${temp_dir}/bats-core"; then
        log_error "Failed to clone BATS repository"
        rm -rf "${temp_dir}"
        return 1
    fi
    
    log_info "Installing BATS to ${install_dir}..."
    cd "${temp_dir}/bats-core"
    
    if ! sudo ./install.sh "${install_dir}"; then
        log_error "Failed to install BATS"
        rm -rf "${temp_dir}"
        return 1
    fi
    
    # Clean up
    rm -rf "${temp_dir}"
    
    log_success "BATS installed successfully"
}

# Install using npm
install_npm() {
    log_info "Installing BATS using npm..."
    
    if ! command -v npm &> /dev/null; then
        log_error "npm is not installed"
        return 1
    fi
    
    if npm install -g bats; then
        log_success "BATS installed successfully via npm"
        return 0
    else
        log_error "Failed to install BATS via npm"
        return 1
    fi
}

# Main installation function
main() {
    echo "======================================"
    echo "  BATS Installation Script"
    echo "======================================"
    echo ""
    
    # Check if already installed
    if check_bats_installed; then
        echo ""
        log_info "No installation needed"
        exit 0
    fi
    
    # Detect OS
    local os
    os=$(detect_os)
    log_info "Detected operating system: ${os}"
    echo ""
    
    # Install based on OS
    case "${os}" in
        linux)
            if install_linux; then
                echo ""
                log_success "BATS installation completed!"
            else
                log_error "BATS installation failed"
                exit 1
            fi
            ;;
        macos)
            if install_macos; then
                echo ""
                log_success "BATS installation completed!"
            else
                log_error "BATS installation failed"
                exit 1
            fi
            ;;
        windows)
            log_warning "Windows detected. Please install BATS manually:"
            echo ""
            echo "  1. Install Git Bash or WSL"
            echo "  2. Follow Linux installation instructions"
            echo ""
            exit 1
            ;;
        *)
            log_warning "Unknown operating system. Attempting manual installation..."
            if install_manual; then
                echo ""
                log_success "BATS installation completed!"
            else
                log_error "BATS installation failed"
                exit 1
            fi
            ;;
    esac
    
    # Verify installation
    echo ""
    if check_bats_installed; then
        echo ""
        log_info "You can now run the test suite:"
        echo "  ./scripts/test/run_tests.sh"
    else
        log_error "BATS installation verification failed"
        exit 1
    fi
}

# Run main function
main "$@"
