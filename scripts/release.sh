#!/usr/bin/env bash

# Release automation script for Go Zalo Bot SDK
# Automates version management, validation, and release creation

set -euo pipefail

# Source utility functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
source "${SCRIPT_DIR}/utils.sh"

# Configuration
VERSION=""
RELEASE_MESSAGE=""
DRY_RUN=false
REMOTE="origin"

# Exit codes
readonly EXIT_SUCCESS=0
readonly EXIT_ERROR=1
readonly EXIT_INVALID_ARGS=2
readonly EXIT_VALIDATION_FAILED=4

# Show usage information
show_help() {
    cat << EOF
Usage: $(basename "$0") <version> [options]

Automate version management and release creation for the Go Zalo Bot SDK.

Arguments:
  version          Semantic version (e.g., v1.2.3)

Options:
  -m, --message    Release message/notes (optional)
  -d, --dry-run    Show what would be done without executing
  -h, --help       Show this help message

Examples:
  $(basename "$0") v1.2.3
  $(basename "$0") v1.2.3 --message "Bug fixes and improvements"
  $(basename "$0") v2.0.0 --dry-run

Exit Codes:
  0 - Release created successfully
  1 - General error
  2 - Invalid arguments
  4 - Validation failed

EOF
}

# Parse command line arguments
parse_args() {
    # Check for help flag first
    for arg in "$@"; do
        if [[ "${arg}" == "-h" ]] || [[ "${arg}" == "--help" ]]; then
            show_help
            exit "${EXIT_SUCCESS}"
        fi
    done
    
    if [[ $# -eq 0 ]]; then
        log_error "Version argument is required" "${EXIT_INVALID_ARGS}"
    fi
    
    # First argument should be the version
    VERSION="$1"
    shift
    
    # Parse options
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--message)
                if [[ $# -lt 2 ]]; then
                    log_error "Option --message requires an argument" "${EXIT_INVALID_ARGS}"
                fi
                RELEASE_MESSAGE="$2"
                shift 2
                ;;
            -d|--dry-run)
                DRY_RUN=true
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

# Validate semantic version format
validate_version_format() {
    log_info "Validating version format..."
    
    if ! validate_semver "${VERSION}"; then
        log_error "Invalid version format: ${VERSION}
Expected format: v*.*.* (e.g., v1.2.3, v2.0.0-beta.1)
Semantic versioning: MAJOR.MINOR.PATCH[-PRERELEASE][+METADATA]" "${EXIT_INVALID_ARGS}"
    fi
    
    # Ensure version starts with 'v'
    if [[ ! "${VERSION}" =~ ^v ]]; then
        log_error "Version must start with 'v' (e.g., v1.2.3)" "${EXIT_INVALID_ARGS}"
    fi
    
    log_success "Version format is valid: ${VERSION}"
}

# Check for duplicate tags
check_duplicate_tag() {
    log_info "Checking for duplicate tags..."
    
    if tag_exists "${VERSION}"; then
        log_error "Tag ${VERSION} already exists. Please use a different version." "${EXIT_ERROR}"
    fi
    
    log_success "No duplicate tag found"
}

# Check Git working directory is clean
check_git_status() {
    log_info "Checking Git working directory..."
    
    if ! is_git_clean; then
        log_warning "Working directory has uncommitted changes"
        git status --short
        echo ""
        log_error "Please commit or stash your changes before creating a release" "${EXIT_ERROR}"
    fi
    
    log_success "Working directory is clean"
}

# Run build validation
run_build_validation() {
    log_info "Running build validation..."
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY RUN] Would run: ${SCRIPT_DIR}/build.sh"
        return 0
    fi
    
    if ! "${SCRIPT_DIR}/build.sh"; then
        log_error "Build validation failed" "${EXIT_VALIDATION_FAILED}"
    fi
    
    echo ""
}

# Run tests
run_tests() {
    log_info "Running tests..."
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY RUN] Would run: ${SCRIPT_DIR}/test.sh --coverage"
        return 0
    fi
    
    if ! "${SCRIPT_DIR}/test.sh" --coverage; then
        log_error "Tests failed" "${EXIT_VALIDATION_FAILED}"
    fi
    
    echo ""
}

# Run linting
run_linting() {
    log_info "Running linting..."
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY RUN] Would run: ${SCRIPT_DIR}/lint.sh"
        return 0
    fi
    
    if ! "${SCRIPT_DIR}/lint.sh"; then
        log_error "Linting failed" "${EXIT_VALIDATION_FAILED}"
    fi
    
    echo ""
}

# Run pre-release validation
run_validation() {
    log_info "Running pre-release validation..."
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY RUN] Would run: ${SCRIPT_DIR}/validate.sh"
        return 0
    fi
    
    if ! "${SCRIPT_DIR}/validate.sh"; then
        log_error "Pre-release validation failed" "${EXIT_VALIDATION_FAILED}"
    fi
    
    echo ""
}

# Update CHANGELOG.md with new version
update_changelog() {
    log_info "Updating CHANGELOG.md..."
    
    local changelog="${PROJECT_ROOT}/CHANGELOG.md"
    
    if [[ ! -f "${changelog}" ]]; then
        log_error "CHANGELOG.md not found at ${changelog}" "${EXIT_ERROR}"
    fi
    
    # Get current date in YYYY-MM-DD format
    local release_date
    release_date=$(date +%Y-%m-%d)
    
    # Extract version number without 'v' prefix
    local version_number="${VERSION#v}"
    
    # Check if [Unreleased] section exists
    if ! grep -q "^## \[Unreleased\]" "${changelog}"; then
        log_warning "No [Unreleased] section found in CHANGELOG.md"
        log_info "Creating new version section..."
    fi
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY RUN] Would update CHANGELOG.md:"
        log_info "  - Add version section: ## [${version_number}] - ${release_date}"
        if [[ -n "${RELEASE_MESSAGE}" ]]; then
            log_info "  - Add release message: ${RELEASE_MESSAGE}"
        fi
        return 0
    fi
    
    # Create a temporary file
    local temp_file
    temp_file=$(mktemp)
    
    # Flag to track if we've added the new version
    local version_added=false
    
    # Read the changelog line by line
    while IFS= read -r line; do
        echo "${line}" >> "${temp_file}"
        
        # If we find the [Unreleased] section and haven't added version yet
        if [[ "${line}" =~ ^\#\#[[:space:]]\[Unreleased\] ]] && [[ "${version_added}" == false ]]; then
            # Add a blank line after [Unreleased]
            echo "" >> "${temp_file}"
            
            # Add the new version section
            echo "## [${version_number}] - ${release_date}" >> "${temp_file}"
            echo "" >> "${temp_file}"
            
            # Add release message if provided
            if [[ -n "${RELEASE_MESSAGE}" ]]; then
                echo "### Release Notes" >> "${temp_file}"
                echo "" >> "${temp_file}"
                echo "${RELEASE_MESSAGE}" >> "${temp_file}"
                echo "" >> "${temp_file}"
            fi
            
            version_added=true
        fi
    done < "${changelog}"
    
    # If we didn't find [Unreleased] section, add version after the header
    if [[ "${version_added}" == false ]]; then
        # Reset temp file
        rm -f "${temp_file}"
        temp_file=$(mktemp)
        
        local header_found=false
        while IFS= read -r line; do
            echo "${line}" >> "${temp_file}"
            
            # Add after the main header and description
            if [[ "${line}" =~ ^#[[:space:]]Changelog ]] && [[ "${header_found}" == false ]]; then
                # Skip until we find a blank line or next section
                while IFS= read -r next_line; do
                    echo "${next_line}" >> "${temp_file}"
                    if [[ -z "${next_line}" ]] || [[ "${next_line}" =~ ^\#\# ]]; then
                        break
                    fi
                done
                
                # Add new version section
                echo "" >> "${temp_file}"
                echo "## [${version_number}] - ${release_date}" >> "${temp_file}"
                echo "" >> "${temp_file}"
                
                if [[ -n "${RELEASE_MESSAGE}" ]]; then
                    echo "### Release Notes" >> "${temp_file}"
                    echo "" >> "${temp_file}"
                    echo "${RELEASE_MESSAGE}" >> "${temp_file}"
                    echo "" >> "${temp_file}"
                fi
                
                header_found=true
                
                # Continue with rest of file
                while IFS= read -r remaining_line; do
                    echo "${remaining_line}" >> "${temp_file}"
                done
                break
            fi
        done < "${changelog}"
    fi
    
    # Replace original changelog with updated version
    mv "${temp_file}" "${changelog}"
    
    log_success "CHANGELOG.md updated with version ${VERSION}"
}

# Create Git tag
create_git_tag() {
    log_info "Creating Git tag..."
    
    local tag_message="Release ${VERSION}"
    if [[ -n "${RELEASE_MESSAGE}" ]]; then
        tag_message="${tag_message}: ${RELEASE_MESSAGE}"
    fi
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY RUN] Would create tag: ${VERSION}"
        log_info "[DRY RUN] Tag message: ${tag_message}"
        return 0
    fi
    
    # Create annotated tag
    if ! git tag -a "${VERSION}" -m "${tag_message}"; then
        log_error "Failed to create Git tag" "${EXIT_ERROR}"
    fi
    
    log_success "Git tag ${VERSION} created"
}

# Push tag to remote
push_tag() {
    log_info "Pushing tag to remote..."
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY RUN] Would push tag to ${REMOTE}: git push ${REMOTE} ${VERSION}"
        return 0
    fi
    
    # Push the tag
    if ! git push "${REMOTE}" "${VERSION}"; then
        log_error "Failed to push tag to remote. You may need to push manually:
  git push ${REMOTE} ${VERSION}" "${EXIT_ERROR}"
    fi
    
    log_success "Tag ${VERSION} pushed to ${REMOTE}"
}

# Display next steps
display_next_steps() {
    echo ""
    echo "======================================"
    echo "       RELEASE COMPLETED"
    echo "======================================"
    echo ""
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_info "This was a dry run. No changes were made."
        echo ""
        log_info "To create the release for real, run:"
        echo "  $(basename "$0") ${VERSION}"
        if [[ -n "${RELEASE_MESSAGE}" ]]; then
            echo "  with message: ${RELEASE_MESSAGE}"
        fi
    else
        log_success "Release ${VERSION} has been created successfully!"
        echo ""
        log_info "Next steps:"
        echo ""
        echo "1. Monitor the CI/CD pipeline:"
        echo "   https://github.com/$(git config --get remote.${REMOTE}.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions"
        echo ""
        echo "2. Verify the release on GitHub:"
        echo "   https://github.com/$(git config --get remote.${REMOTE}.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/releases/tag/${VERSION}"
        echo ""
        echo "3. Check pkg.go.dev indexing (may take a few minutes):"
        echo "   https://pkg.go.dev/$(grep '^module ' "${PROJECT_ROOT}/go.mod" | awk '{print $2}')@${VERSION}"
        echo ""
        echo "4. Announce the release to users and contributors"
        echo ""
    fi
    
    echo "======================================"
}

# Main execution
main() {
    # Change to project root
    cd "${PROJECT_ROOT}"
    
    log_info "Go Zalo Bot SDK - Release Automation"
    echo ""
    
    if [[ "${DRY_RUN}" == true ]]; then
        log_warning "DRY RUN MODE - No changes will be made"
        echo ""
    fi
    
    # Pre-flight checks
    log_info "Running pre-flight checks..."
    echo ""
    
    validate_version_format
    check_duplicate_tag
    check_git_status
    
    echo ""
    log_info "Starting validation suite..."
    echo ""
    
    # Run full validation suite
    run_build_validation
    run_tests
    run_linting
    run_validation
    
    log_success "All validation checks passed!"
    echo ""
    
    # Update changelog and create release
    log_info "Creating release..."
    echo ""
    
    update_changelog
    
    # Commit changelog changes if not dry run
    if [[ "${DRY_RUN}" == false ]]; then
        log_info "Committing CHANGELOG.md changes..."
        git add "${PROJECT_ROOT}/CHANGELOG.md"
        git commit -m "chore: update CHANGELOG for ${VERSION}"
        log_success "CHANGELOG.md committed"
        echo ""
    else
        log_info "[DRY RUN] Would commit CHANGELOG.md changes"
        echo ""
    fi
    
    create_git_tag
    push_tag
    
    # Display next steps
    display_next_steps
}

# Parse arguments and run
parse_args "$@"
main

