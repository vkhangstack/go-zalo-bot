# Script Testing Suite

This directory contains automated tests for the build and deployment automation scripts using BATS (Bash Automated Testing System).

## Overview

The test suite provides comprehensive coverage for all automation scripts:

- **test_utils.bats** - Unit tests for utility functions in `utils.sh`
- **test_build.bats** - Integration tests for `build.sh`
- **test_test.bats** - Integration tests for `test.sh`
- **test_lint.bats** - Integration tests for `lint.sh`
- **test_validate.bats** - Integration tests for `validate.sh`
- **test_release.bats** - Integration tests for `release.sh`

## Prerequisites

### Install BATS

BATS (Bash Automated Testing System) must be installed to run the tests.

#### Using npm:
```bash
npm install -g bats
```

#### Using Homebrew (macOS):
```bash
brew install bats-core
```

#### Using apt (Ubuntu/Debian):
```bash
sudo apt-get install bats
```

#### Manual Installation:
```bash
git clone https://github.com/bats-core/bats-core.git
cd bats-core
sudo ./install.sh /usr/local
```

For more information, visit: https://github.com/bats-core/bats-core

## Running Tests

### Run All Tests

```bash
./scripts/test/run_tests.sh
```

### Run Specific Test File

```bash
./scripts/test/run_tests.sh test_utils.bats
```

### Run with Verbose Output

```bash
./scripts/test/run_tests.sh --verbose
```

### Run Tests Directly with BATS

```bash
cd scripts/test
bats test_utils.bats
bats test_*.bats
```

## Test Structure

### Unit Tests (test_utils.bats)

Tests individual utility functions in isolation:

- Logging functions (`log_info`, `log_success`, `log_error`, `log_warning`)
- Command checking (`check_command`, `check_go_version`)
- Semantic version validation (`validate_semver`)
- Git utility functions (`get_latest_tag`, `tag_exists`, `is_git_clean`)

### Integration Tests

Tests complete script workflows:

- **test_build.bats** - Validates build script functionality
  - Help message display
  - Command-line argument parsing
  - Go version verification
  - Dependency verification
  - Package compilation
  - Example compilation

- **test_test.bats** - Validates test script functionality
  - Test execution
  - Coverage report generation
  - Race detector integration
  - Test summary display

- **test_lint.bats** - Validates linting script functionality
  - Code formatting checks
  - Static analysis with go vet
  - golangci-lint integration
  - Godoc comment validation

- **test_validate.bats** - Validates pre-release validation script
  - Example compilation verification
  - CHANGELOG.md format validation
  - README.md validation
  - Godoc coverage checks
  - go.mod version verification

- **test_release.bats** - Validates release automation script
  - Version format validation
  - Duplicate tag checking
  - CHANGELOG updates
  - Git tag creation
  - Dry-run mode

## Test Helpers

### setup_suite.bash

Provides common test utilities:

- `setup_test_environment()` - Creates temporary test directories
- `teardown_test_environment()` - Cleans up test artifacts
- `create_mock_git_repo()` - Creates a mock Git repository for testing
- `create_test_go_module()` - Creates a test Go module

## Writing New Tests

### Test File Template

```bash
#!/usr/bin/env bats

# Description of test file

load setup_suite

setup() {
    setup_test_environment
}

teardown() {
    teardown_test_environment
}

@test "descriptive test name" {
    run command_to_test
    [ "$status" -eq 0 ]
    [[ "$output" =~ "expected output" ]]
}
```

### Best Practices

1. **Use descriptive test names** - Test names should clearly describe what is being tested
2. **Test one thing per test** - Each test should focus on a single behavior
3. **Use setup and teardown** - Clean up test artifacts to avoid side effects
4. **Check exit codes** - Always verify the command exit status
5. **Validate output** - Check that output contains expected messages
6. **Test error cases** - Include tests for invalid inputs and error conditions
7. **Use dry-run mode** - For destructive operations, test with dry-run flags

### BATS Assertions

```bash
# Check exit status
[ "$status" -eq 0 ]          # Success
[ "$status" -ne 0 ]          # Failure

# Check output
[[ "$output" =~ "pattern" ]] # Contains pattern
[ "$output" = "exact" ]      # Exact match
[ -z "$output" ]             # Empty output
[ -n "$output" ]             # Non-empty output

# File checks
[ -f "file.txt" ]            # File exists
[ -d "directory" ]           # Directory exists
[ -x "script.sh" ]           # File is executable
```

## Continuous Integration

The test suite is integrated into the CI/CD pipeline and runs automatically on:

- Pull requests
- Pushes to main branch
- Before releases

See `.github/workflows/ci.yml` for CI configuration.

## Troubleshooting

### Tests Fail Locally

1. Ensure BATS is installed: `bats --version`
2. Check that you're in the project root
3. Verify all scripts are executable: `chmod +x scripts/*.sh`
4. Run tests with verbose output: `./scripts/test/run_tests.sh --verbose`

### Permission Denied Errors

Make the test runner executable:
```bash
chmod +x scripts/test/run_tests.sh
```

### Git-Related Test Failures

Some tests create temporary Git repositories. Ensure Git is installed and configured:
```bash
git config --global user.email "test@example.com"
git config --global user.name "Test User"
```

### Timeout Issues

Some integration tests may take longer on slower systems. Increase timeout if needed.

## Coverage

The test suite aims for comprehensive coverage:

- ✅ All utility functions
- ✅ All script command-line interfaces
- ✅ Success paths for all scripts
- ✅ Error handling and validation
- ✅ Edge cases and invalid inputs

## Contributing

When adding new scripts or modifying existing ones:

1. Add corresponding tests to the appropriate test file
2. Run the full test suite before committing
3. Ensure all tests pass in CI
4. Update this README if adding new test files

## Resources

- [BATS Documentation](https://bats-core.readthedocs.io/)
- [BATS GitHub Repository](https://github.com/bats-core/bats-core)
- [Bash Testing Best Practices](https://github.com/bats-core/bats-core#writing-tests)
