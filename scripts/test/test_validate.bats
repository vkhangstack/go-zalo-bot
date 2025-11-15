#!/usr/bin/env bats

# Integration tests for validate.sh

load setup_suite

setup() {
    setup_test_environment
}

teardown() {
    teardown_test_environment
}

@test "validate.sh shows help message" {
    run "${SCRIPTS_DIR}/validate.sh" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
    [[ "$output" =~ "Pre-release validation" ]]
}

@test "validate.sh accepts verbose flag" {
    run "${SCRIPTS_DIR}/validate.sh" --verbose
    [[ "$output" =~ "Starting pre-release validation" ]]
}

@test "validate.sh rejects invalid option" {
    run "${SCRIPTS_DIR}/validate.sh" --invalid-option
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Unknown option" ]]
}

@test "validate.sh validates examples" {
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "Validating example compilation" ]]
}

@test "validate.sh validates CHANGELOG" {
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "Validating CHANGELOG.md" ]]
}

@test "validate.sh validates README" {
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "Validating README.md" ]]
}

@test "validate.sh checks godoc coverage" {
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "Checking godoc comment coverage" ]]
}

@test "validate.sh verifies go.mod version" {
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "Verifying go.mod version" ]]
}

@test "validate.sh generates validation report" {
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "VALIDATION REPORT" ]]
    [[ "$output" =~ "Summary:" ]]
}

@test "validate.sh reports errors, warnings, and info" {
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "Errors:" ]]
    [[ "$output" =~ "Warnings:" ]]
}

@test "validate.sh succeeds in project root" {
    cd "${PROJECT_ROOT}"
    run "${SCRIPTS_DIR}/validate.sh"
    [[ "$output" =~ "VALIDATION REPORT" ]]
}

@test "validate.sh with verbose shows detailed info" {
    run "${SCRIPTS_DIR}/validate.sh" --verbose
    [[ "$output" =~ "Starting pre-release validation" ]]
}
