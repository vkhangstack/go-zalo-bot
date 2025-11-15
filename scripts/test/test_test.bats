#!/usr/bin/env bats

# Integration tests for test.sh

load setup_suite

setup() {
    setup_test_environment
}

teardown() {
    teardown_test_environment
}

@test "test.sh shows help message" {
    run "${SCRIPTS_DIR}/test.sh" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
    [[ "$output" =~ "Run all tests" ]]
}

@test "test.sh accepts coverage flag" {
    run "${SCRIPTS_DIR}/test.sh" --coverage
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Coverage reporting enabled" ]]
}

@test "test.sh accepts verbose flag" {
    run "${SCRIPTS_DIR}/test.sh" --verbose
    [ "$status" -eq 0 ]
}

@test "test.sh accepts race flag" {
    run "${SCRIPTS_DIR}/test.sh" --race
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Race detector enabled" ]]
}

@test "test.sh accepts combined flags" {
    run "${SCRIPTS_DIR}/test.sh" --coverage --race --verbose
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Coverage reporting enabled" ]]
    [[ "$output" =~ "Race detector enabled" ]]
}

@test "test.sh rejects invalid option" {
    run "${SCRIPTS_DIR}/test.sh" --invalid-option
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Unknown option" ]]
}

@test "test.sh runs tests successfully" {
    run "${SCRIPTS_DIR}/test.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Running tests" ]]
}

@test "test.sh displays test summary" {
    run "${SCRIPTS_DIR}/test.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Test Summary" ]]
}

@test "test.sh with coverage generates report" {
    cd "${PROJECT_ROOT}"
    run "${SCRIPTS_DIR}/test.sh" --coverage
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Coverage" ]]
    [ -f "${PROJECT_ROOT}/coverage.out" ]
}

@test "test.sh reports test duration" {
    run "${SCRIPTS_DIR}/test.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "passed in" ]]
    [[ "$output" =~ "s" ]]
}

@test "test.sh succeeds in project root" {
    cd "${PROJECT_ROOT}"
    run "${SCRIPTS_DIR}/test.sh"
    [ "$status" -eq 0 ]
}
