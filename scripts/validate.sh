#!/bin/bash

# Pre-release validation script
# Performs comprehensive validation checks before releasing a new version

set -euo pipefail

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Source utility functions
# shellcheck source=scripts/utils.sh
source "${SCRIPT_DIR}/utils.sh"

# Configuration
VERBOSE=false
EXIT_CODE=0

# Validation counters
ERROR_COUNT=0
WARNING_COUNT=0
INFO_COUNT=0

# Validation results storage
declare -a VALIDATION_ERRORS=()
declare -a VALIDATION_WARNINGS=()
declare -a VALIDATION_INFO=()

# Print usage information
usage() {
    cat << EOF
Usage: $(basename "$0") [options]

Pre-release validation script that performs comprehensive checks before releasing.

Options:
    -v, --verbose    Enable verbose output
    -h, --help       Show this help message

Validation Checks:
    - Example compilation verification
    - CHANGELOG.md format validation
    - README.md validation
    - Godoc comment coverage check
    - go.mod version verification

Exit Codes:
    0 - All validations passed
    1 - General error
    4 - Validation failure

EOF
}

# Add validation issue with severity level
# Usage: add_issue "severity" "message"
add_issue() {
    local severity="$1"
    local message="$2"
    
    case "${severity}" in
        error)
            VALIDATION_ERRORS+=("${message}")
            ERROR_COUNT=$((ERROR_COUNT + 1))
            ;;
        warning)
            VALIDATION_WARNINGS+=("${message}")
            WARNING_COUNT=$((WARNING_COUNT + 1))
            ;;
        info)
            VALIDATION_INFO+=("${message}")
            INFO_COUNT=$((INFO_COUNT + 1))
            ;;
    esac
}

# Validate example compilation
validate_examples() {
    log_info "Validating example compilation..."
    
    local examples_dir="${PROJECT_ROOT}/examples"
    
    if [ ! -d "${examples_dir}" ]; then
        add_issue "warning" "Examples directory not found at ${examples_dir}"
        return
    fi
    
    # Find all example directories with main.go files
    local example_count=0
    local failed_count=0
    
    # Use a temporary file to store find results
    local temp_file
    temp_file=$(mktemp)
    find "${examples_dir}" -name "main.go" -type f > "${temp_file}"
    
    while IFS= read -r example_file; do
        if [ -z "${example_file}" ]; then
            continue
        fi
        
        example_count=$((example_count + 1))
        local example_dir
        example_dir=$(dirname "${example_file}")
        local example_name
        example_name=$(basename "${example_dir}")
        
        if [ "${VERBOSE}" = true ]; then
            log_info "Compiling example: ${example_name}"
        fi
        
        # Try to compile the example
        local relative_example_path
        relative_example_path=$(realpath --relative-to="${PROJECT_ROOT}" "${example_dir}")
        
        if go build -o /dev/null "./${relative_example_path}" 2>/dev/null; then
            if [ "${VERBOSE}" = true ]; then
                add_issue "info" "Example '${example_name}' compiled successfully"
            fi
        else
            add_issue "error" "Example '${example_name}' failed to compile"
            failed_count=$((failed_count + 1))
        fi
    done < "${temp_file}"
    
    rm -f "${temp_file}"
    
    if [ "${example_count}" -eq 0 ]; then
        add_issue "warning" "No example files found in ${examples_dir}"
    elif [ "${failed_count}" -eq 0 ]; then
        log_success "All ${example_count} examples compiled successfully"
    else
        add_issue "error" "${failed_count} out of ${example_count} examples failed to compile"
    fi
}

# Validate CHANGELOG.md format
validate_changelog() {
    log_info "Validating CHANGELOG.md format..."
    
    local changelog="${PROJECT_ROOT}/CHANGELOG.md"
    
    if [ ! -f "${changelog}" ]; then
        add_issue "error" "CHANGELOG.md not found at ${changelog}"
        return
    fi
    
    # Check for required header
    if ! grep -q "^# Changelog" "${changelog}"; then
        add_issue "error" "CHANGELOG.md missing '# Changelog' header"
    fi
    
    # Check for Keep a Changelog reference
    if ! grep -q "Keep a Changelog" "${changelog}"; then
        add_issue "warning" "CHANGELOG.md should reference Keep a Changelog format"
    fi
    
    # Check for Semantic Versioning reference
    if ! grep -q "Semantic Versioning" "${changelog}"; then
        add_issue "warning" "CHANGELOG.md should reference Semantic Versioning"
    fi
    
    # Check for version sections (format: ## [X.Y.Z] - YYYY-MM-DD)
    if ! grep -qE "^## \[[0-9]+\.[0-9]+\.[0-9]+\]" "${changelog}"; then
        add_issue "error" "CHANGELOG.md missing properly formatted version sections (## [X.Y.Z] - YYYY-MM-DD)"
    fi
    
    # Check for standard sections (Added, Changed, Fixed, etc.)
    local has_sections=false
    for section in "Added" "Changed" "Fixed" "Removed" "Deprecated" "Security"; do
        if grep -q "^### ${section}" "${changelog}"; then
            has_sections=true
            break
        fi
    done
    
    if [ "${has_sections}" = false ]; then
        add_issue "warning" "CHANGELOG.md should include standard sections (Added, Changed, Fixed, etc.)"
    fi
    
    # Check for [Unreleased] section
    if ! grep -q "^## \[Unreleased\]" "${changelog}"; then
        add_issue "info" "CHANGELOG.md could include an [Unreleased] section for upcoming changes"
    fi
    
    log_success "CHANGELOG.md format validation completed"
}

# Validate README.md
validate_readme() {
    log_info "Validating README.md..."
    
    local readme="${PROJECT_ROOT}/README.md"
    
    if [ ! -f "${readme}" ]; then
        add_issue "error" "README.md not found at ${readme}"
        return
    fi
    
    # Check for installation instructions
    if ! grep -qi "installation" "${readme}"; then
        add_issue "error" "README.md missing installation section"
    fi
    
    # Check for go get command
    if ! grep -q "go get" "${readme}"; then
        add_issue "error" "README.md missing 'go get' installation command"
    fi
    
    # Check for usage examples
    if ! grep -qi "usage\|example\|quick start" "${readme}"; then
        add_issue "warning" "README.md should include usage examples or quick start guide"
    fi
    
    # Check for code blocks
    if ! grep -q '```' "${readme}"; then
        add_issue "warning" "README.md should include code examples in code blocks"
    fi
    
    # Check for project description
    local first_line
    first_line=$(head -n 1 "${readme}")
    if [[ ! "${first_line}" =~ ^#[[:space:]] ]]; then
        add_issue "warning" "README.md should start with a project title (# Title)"
    fi
    
    # Check for license information
    if ! grep -qi "license" "${readme}"; then
        add_issue "info" "README.md could include license information"
    fi
    
    log_success "README.md validation completed"
}

# Check godoc comment coverage for exported symbols
validate_godoc_coverage() {
    log_info "Checking godoc comment coverage..."
    
    local missing_docs=0
    local total_exported=0
    
    # Find all Go files (excluding test files and examples)
    local temp_file
    temp_file=$(mktemp)
    find "${PROJECT_ROOT}" -name "*.go" -type f > "${temp_file}"
    
    while IFS= read -r go_file; do
        if [ -z "${go_file}" ]; then
            continue
        fi
        
        # Skip if file is in vendor, examples, or test files
        if [[ "${go_file}" =~ vendor/ ]] || [[ "${go_file}" =~ examples/ ]] || [[ "${go_file}" =~ _test\.go$ ]]; then
            continue
        fi
        
        # Extract exported symbols (functions, types, constants, variables)
        # Look for lines starting with: func, type, const, var followed by uppercase letter
        local exported_symbols
        exported_symbols=$(grep -nE '^(func|type|const|var)[[:space:]]+[A-Z]' "${go_file}" || true)
        
        if [ -z "${exported_symbols}" ]; then
            continue
        fi
        
        # Check each exported symbol for documentation
        while IFS= read -r symbol_line; do
            if [ -z "${symbol_line}" ]; then
                continue
            fi
            
            total_exported=$((total_exported + 1))
            
            local line_num
            line_num=$(echo "${symbol_line}" | cut -d: -f1)
            local symbol_def
            symbol_def=$(echo "${symbol_line}" | cut -d: -f2-)
            
            # Extract symbol name
            local symbol_name
            symbol_name=$(echo "${symbol_def}" | awk '{print $2}' | cut -d'(' -f1)
            
            # Check if there's a comment on the line before
            local prev_line=$((line_num - 1))
            if [ "${prev_line}" -gt 0 ]; then
                local comment
                comment=$(sed -n "${prev_line}p" "${go_file}")
                
                # Check if comment exists and starts with // and mentions the symbol name
                if [[ ! "${comment}" =~ ^//[[:space:]]*${symbol_name} ]]; then
                    local relative_path
                    relative_path=$(realpath --relative-to="${PROJECT_ROOT}" "${go_file}")
                    add_issue "warning" "Missing godoc comment for exported symbol '${symbol_name}' in ${relative_path}:${line_num}"
                    missing_docs=$((missing_docs + 1))
                    
                    if [ "${VERBOSE}" = true ]; then
                        log_warning "  ${symbol_def}"
                    fi
                fi
            fi
        done <<< "${exported_symbols}"
        
    done < "${temp_file}"
    
    rm -f "${temp_file}"
    
    if [ "${total_exported}" -eq 0 ]; then
        add_issue "warning" "No exported symbols found in the project"
    else
        local coverage_percent
        coverage_percent=$(awk "BEGIN {printf \"%.1f\", (1 - ${missing_docs}/${total_exported}) * 100}")
        
        if [ "${missing_docs}" -eq 0 ]; then
            log_success "All ${total_exported} exported symbols have godoc comments (100% coverage)"
        else
            log_warning "Godoc coverage: ${coverage_percent}% (${missing_docs}/${total_exported} symbols missing documentation)"
            
            if [ "${missing_docs}" -gt 10 ]; then
                add_issue "error" "Too many exported symbols without documentation (${missing_docs})"
            fi
        fi
    fi
}

# Verify go.mod version against latest Git tag
validate_gomod_version() {
    log_info "Verifying go.mod version..."
    
    local gomod="${PROJECT_ROOT}/go.mod"
    
    if [ ! -f "${gomod}" ]; then
        add_issue "error" "go.mod not found at ${gomod}"
        return
    fi
    
    # Get module path from go.mod
    local module_path
    module_path=$(grep "^module " "${gomod}" | awk '{print $2}')
    
    if [ -z "${module_path}" ]; then
        add_issue "error" "Could not extract module path from go.mod"
        return
    fi
    
    log_info "Module path: ${module_path}"
    
    # Get latest Git tag
    local latest_tag
    latest_tag=$(get_latest_tag)
    
    if [ -z "${latest_tag}" ]; then
        add_issue "info" "No Git tags found - this might be the first release"
        return
    fi
    
    log_info "Latest Git tag: ${latest_tag}"
    
    # Validate tag format
    if ! validate_semver "${latest_tag}"; then
        add_issue "warning" "Latest Git tag '${latest_tag}' does not follow semantic versioning format"
    fi
    
    # Check if go.mod has version suffix
    if [[ "${module_path}" =~ /v[0-9]+$ ]]; then
        local gomod_version
        gomod_version=$(echo "${module_path}" | grep -oE '/v[0-9]+$')
        local tag_major
        tag_major=$(echo "${latest_tag}" | grep -oE '^v[0-9]+' | sed 's/v//')
        
        if [ "${tag_major}" -ge 2 ]; then
            local expected_suffix="/v${tag_major}"
            if [ "${gomod_version}" != "${expected_suffix}" ]; then
                add_issue "error" "go.mod module path should have suffix '${expected_suffix}' for version ${latest_tag}"
            else
                log_success "go.mod version suffix matches Git tag major version"
            fi
        fi
    else
        # For v0 and v1, no suffix is required
        local tag_major
        tag_major=$(echo "${latest_tag}" | grep -oE '^v[0-9]+' | sed 's/v//')
        
        if [ "${tag_major}" -ge 2 ]; then
            add_issue "error" "go.mod module path should have /v${tag_major} suffix for version ${latest_tag}"
        else
            log_success "go.mod version is appropriate for ${latest_tag}"
        fi
    fi
}

# Generate validation report
generate_report() {
    echo ""
    echo "======================================"
    echo "       VALIDATION REPORT"
    echo "======================================"
    echo ""
    
    # Summary
    echo "Summary:"
    echo "  Errors:   ${ERROR_COUNT}"
    echo "  Warnings: ${WARNING_COUNT}"
    echo "  Info:     ${INFO_COUNT}"
    echo ""
    
    # Errors
    if [ "${ERROR_COUNT}" -gt 0 ]; then
        echo -e "${COLOR_RED}Errors:${COLOR_RESET}"
        for error in "${VALIDATION_ERRORS[@]}"; do
            echo -e "  ${COLOR_RED}✗${COLOR_RESET} ${error}"
        done
        echo ""
    fi
    
    # Warnings
    if [ "${WARNING_COUNT}" -gt 0 ]; then
        echo -e "${COLOR_YELLOW}Warnings:${COLOR_RESET}"
        for warning in "${VALIDATION_WARNINGS[@]}"; do
            echo -e "  ${COLOR_YELLOW}⚠${COLOR_RESET} ${warning}"
        done
        echo ""
    fi
    
    # Info
    if [ "${INFO_COUNT}" -gt 0 ] && [ "${VERBOSE}" = true ]; then
        echo -e "${COLOR_CYAN}Info:${COLOR_RESET}"
        for info in "${VALIDATION_INFO[@]}"; do
            echo -e "  ${COLOR_CYAN}ℹ${COLOR_RESET} ${info}"
        done
        echo ""
    fi
    
    # Final result
    echo "======================================"
    if [ "${ERROR_COUNT}" -eq 0 ]; then
        if [ "${WARNING_COUNT}" -eq 0 ]; then
            echo -e "${COLOR_GREEN}✓ All validations passed!${COLOR_RESET}"
        else
            echo -e "${COLOR_YELLOW}⚠ Validation passed with ${WARNING_COUNT} warning(s)${COLOR_RESET}"
        fi
        echo "======================================"
        return 0
    else
        echo -e "${COLOR_RED}✗ Validation failed with ${ERROR_COUNT} error(s)${COLOR_RESET}"
        echo "======================================"
        return 4
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
                usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1" 2
                ;;
        esac
    done
}

# Main function
main() {
    parse_args "$@"
    
    log_info "Starting pre-release validation..."
    echo ""
    
    # Change to project root
    cd "${PROJECT_ROOT}"
    
    # Run all validation checks
    validate_examples
    validate_changelog
    validate_readme
    validate_godoc_coverage
    validate_gomod_version
    
    # Generate and display report
    echo ""
    if ! generate_report; then
        exit 4
    fi
}

# Run main function
main "$@"
