#!/usr/bin/env bats

# Integration tests for release.sh

load setup_suite

setup() {
    setup_test_environment
    # Create a test repository for release testing
    TEST_REPO="${TEST_TEMP_DIR}/release_test_repo"
    create_mock_git_repo "${TEST_REPO}"
    cd "${TEST_REPO}"
    create_test_go_module "${TEST_REPO}" "github.com/test/release"
    
    # Create required files
    cat > CHANGELOG.md << 'EOF'
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release

EOF
    
    git add .
    git commit -m "Add initial files"
}

teardown() {
    teardown_test_environment
}

@test "release.sh shows help message" {
    run "${SCRIPTS_DIR}/release.sh" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
    [[ "$output" =~ "version management" ]]
}

@test "release.sh requires version argument" {
    run "${SCRIPTS_DIR}/release.sh"
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Version argument is required" ]]
}

@test "release.sh accepts valid semantic version" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Version format is valid" ]]
}

@test "release.sh rejects invalid version format" {
    run "${SCRIPTS_DIR}/release.sh" 1.0 --dry-run
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Invalid version format" ]]
}

@test "release.sh requires v prefix" {
    run "${SCRIPTS_DIR}/release.sh" 1.0.0 --dry-run
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Version must start with 'v'" ]]
}

@test "release.sh accepts message option" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --message "Test release" --dry-run
    [ "$status" -eq 0 ]
}

@test "release.sh dry-run mode shows actions without executing" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "DRY RUN MODE" ]]
    [[ "$output" =~ "Would" ]]
}

@test "release.sh checks for duplicate tags" {
    git tag v1.0.0
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [ "$status" -eq 1 ]
    [[ "$output" =~ "already exists" ]]
}

@test "release.sh validates version format" {
    run "${SCRIPTS_DIR}/release.sh" v1.2.3 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Validating version format" ]]
}

@test "release.sh checks git status" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [[ "$output" =~ "Checking Git working directory" ]]
}

@test "release.sh in dry-run shows validation steps" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Would run" ]]
}

@test "release.sh updates CHANGELOG in dry-run" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Would update CHANGELOG.md" ]]
}

@test "release.sh creates tag in dry-run" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Would create tag" ]]
}

@test "release.sh displays next steps" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --dry-run
    [ "$status" -eq 0 ]
    [[ "$output" =~ "RELEASE COMPLETED" ]]
}

@test "release.sh rejects unknown option" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --invalid-option
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Unknown option" ]]
}

@test "release.sh message option requires argument" {
    run "${SCRIPTS_DIR}/release.sh" v1.0.0 --message
    [ "$status" -eq 2 ]
    [[ "$output" =~ "requires an argument" ]]
}
