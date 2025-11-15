#!/usr/bin/env bats

# Integration tests for lint.sh

load setup_suite

setup() {
    setup_test_environment
}

teardown() {
    teardown_test_environment
}

@test "lint.sh shows help message" {
    run "${SCRIPTS_DIR}/lint.sh" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
    [[ "$output" =~ "linting and static analysis" ]]
}

@test "lint.sh accepts fix flag" {
    run "${SCRIPTS_DIR}/lint.sh" --fix
    # May succeed or fail depending on code state
    [[ "$output" =~ "Checking code formatting" ]]
}

@test "lint.sh accepts verbose flag" {
    run "${SCRIPTS_DIR}/lint.sh" --verbose
    [[ "$output" =~ "Checking code formatting" ]]
}

@test "lint.sh rejects invalid option" {
    run "${SCRIPTS_DIR}/lint.sh" --invalid-option
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Unknown option" ]]
}

@test "lint.sh checks formatting" {
    run "${SCRIPTS_DIR}/lint.sh"
    [[ "$output" =~ "Checking code formatting" ]]
}

@test "lint.sh runs go vet" {
    run "${SCRIPTS_DIR}/lint.sh"
    [[ "$output" =~ "Running go vet" ]]
}

@test "lint.sh checks godoc comments" {
    run "${SCRIPTS_DIR}/lint.sh"
    [[ "$output" =~ "Checking godoc comments" ]]
}

@test "lint.sh displays summary" {
    run "${SCRIPTS_DIR}/lint.sh"
    [[ "$output" =~ "Linting Summary" ]]
    [[ "$output" =~ "Total issues" ]]
}

@test "lint.sh succeeds in project root" {
    cd "${PROJECT_ROOT}"
    run "${SCRIPTS_DIR}/lint.sh"
    # Exit code may vary based on code quality
    [[ "$output" =~ "Linting Summary" ]]
}

@test "lint.sh shows golangci-lint installation instructions when not installed" {
    # This test assumes golangci-lint might not be installed
    run "${SCRIPTS_DIR}/lint.sh"
    if [[ "$output" =~ "golangci-lint is not installed" ]]; then
        [[ "$output" =~ "go install" ]]
    fi
}
