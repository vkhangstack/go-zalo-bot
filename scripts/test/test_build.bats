#!/usr/bin/env bats

# Integration tests for build.sh

load setup_suite

setup() {
    setup_test_environment
}

teardown() {
    teardown_test_environment
}

@test "build.sh shows help message" {
    run "${SCRIPTS_DIR}/build.sh" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
    [[ "$output" =~ "Build validation script" ]]
}

@test "build.sh accepts verbose flag" {
    run "${SCRIPTS_DIR}/build.sh" --verbose
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Starting build validation" ]]
}

@test "build.sh rejects invalid option" {
    run "${SCRIPTS_DIR}/build.sh" --invalid-option
    [ "$status" -eq 2 ]
    [[ "$output" =~ "Unknown option" ]]
}

@test "build.sh verifies Go version" {
    run "${SCRIPTS_DIR}/build.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Verifying Go version" ]]
}

@test "build.sh verifies dependencies" {
    run "${SCRIPTS_DIR}/build.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Verifying module dependencies" ]]
}

@test "build.sh compiles packages" {
    run "${SCRIPTS_DIR}/build.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Compiling all packages" ]]
}

@test "build.sh compiles examples" {
    run "${SCRIPTS_DIR}/build.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Compiling examples" ]]
}

@test "build.sh reports completion time" {
    run "${SCRIPTS_DIR}/build.sh"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "completed successfully in" ]]
    [[ "$output" =~ "s" ]]
}

@test "build.sh succeeds in project root" {
    cd "${PROJECT_ROOT}"
    run "${SCRIPTS_DIR}/build.sh"
    [ "$status" -eq 0 ]
}

@test "build.sh with verbose shows detailed output" {
    run "${SCRIPTS_DIR}/build.sh" --verbose
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Go version" ]]
}
