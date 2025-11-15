#!/usr/bin/env bats

# Unit tests for utils.sh functions

# Load test setup
load setup_suite

# Source the utils script
setup() {
    source "${SCRIPTS_DIR}/utils.sh"
    setup_test_environment
}

teardown() {
    teardown_test_environment
}

# Test logging functions

@test "log_info prints informational message" {
    run log_info "Test message"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "INFO: Test message" ]]
}

@test "log_success prints success message" {
    run log_success "Operation completed"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "SUCCESS: Operation completed" ]]
}

@test "log_warning prints warning message" {
    run log_warning "This is a warning"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "WARNING: This is a warning" ]]
}

@test "log_error exits with default code 1" {
    run log_error "Error occurred"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "ERROR: Error occurred" ]]
}

@test "log_error exits with custom code" {
    run log_error "Custom error" 42
    [ "$status" -eq 42 ]
    [[ "$output" =~ "ERROR: Custom error" ]]
}

# Test command checking functions

@test "check_command succeeds for existing command" {
    run check_command "bash"
    [ "$status" -eq 0 ]
}

@test "check_command fails for non-existing command" {
    run check_command "nonexistent_command_xyz"
    [ "$status" -eq 3 ]
    [[ "$output" =~ "nonexistent_command_xyz is not installed" ]]
}

@test "check_command shows custom installation message" {
    run check_command "nonexistent_cmd" "Install from https://example.com"
    [ "$status" -eq 3 ]
    [[ "$output" =~ "Install from https://example.com" ]]
}

@test "check_go_version succeeds with valid Go version" {
    skip "Requires Go to be installed"
    run check_go_version "1.20"
    [ "$status" -eq 0 ]
}

# Test semantic version validation

@test "validate_semver accepts valid version with v prefix" {
    run validate_semver "v1.2.3"
    [ "$status" -eq 0 ]
}

@test "validate_semver accepts valid version without v prefix" {
    run validate_semver "1.2.3"
    [ "$status" -eq 0 ]
}

@test "validate_semver accepts version with prerelease" {
    run validate_semver "v1.2.3-beta.1"
    [ "$status" -eq 0 ]
}

@test "validate_semver accepts version with metadata" {
    run validate_semver "v1.2.3+build.123"
    [ "$status" -eq 0 ]
}

@test "validate_semver accepts version with prerelease and metadata" {
    run validate_semver "v1.2.3-alpha.1+build.456"
    [ "$status" -eq 0 ]
}

@test "validate_semver rejects invalid version format" {
    run validate_semver "1.2"
    [ "$status" -eq 1 ]
}

@test "validate_semver rejects version with letters" {
    run validate_semver "v1.2.x"
    [ "$status" -eq 1 ]
}

@test "validate_semver rejects empty version" {
    run validate_semver ""
    [ "$status" -eq 1 ]
}

# Test Git utility functions

@test "get_latest_tag returns empty string when no tags exist" {
    local test_repo="${TEST_TEMP_DIR}/test_repo"
    create_mock_git_repo "${test_repo}"
    
    cd "${test_repo}"
    run get_latest_tag
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}

@test "get_latest_tag returns latest tag" {
    local test_repo="${TEST_TEMP_DIR}/test_repo"
    create_mock_git_repo "${test_repo}"
    
    cd "${test_repo}"
    git tag v1.0.0
    git tag v1.1.0
    
    run get_latest_tag
    [ "$status" -eq 0 ]
    [ "$output" = "v1.1.0" ]
}

@test "tag_exists returns 0 for existing tag" {
    local test_repo="${TEST_TEMP_DIR}/test_repo"
    create_mock_git_repo "${test_repo}"
    
    cd "${test_repo}"
    git tag v1.0.0
    
    run tag_exists "v1.0.0"
    [ "$status" -eq 0 ]
}

@test "tag_exists returns 1 for non-existing tag" {
    local test_repo="${TEST_TEMP_DIR}/test_repo"
    create_mock_git_repo "${test_repo}"
    
    cd "${test_repo}"
    
    run tag_exists "v1.0.0"
    [ "$status" -eq 1 ]
}

@test "get_current_branch returns branch name" {
    local test_repo="${TEST_TEMP_DIR}/test_repo"
    create_mock_git_repo "${test_repo}"
    
    cd "${test_repo}"
    run get_current_branch
    [ "$status" -eq 0 ]
    [[ "$output" =~ ^(main|master)$ ]]
}

@test "is_git_clean returns 0 for clean working directory" {
    local test_repo="${TEST_TEMP_DIR}/test_repo"
    create_mock_git_repo "${test_repo}"
    
    cd "${test_repo}"
    run is_git_clean
    [ "$status" -eq 0 ]
}

@test "is_git_clean returns 1 for dirty working directory" {
    local test_repo="${TEST_TEMP_DIR}/test_repo"
    create_mock_git_repo "${test_repo}"
    
    cd "${test_repo}"
    echo "new content" > newfile.txt
    
    run is_git_clean
    [ "$status" -eq 1 ]
}
