# Contributing to Go Zalo Bot SDK

Thank you for your interest in contributing to the Go Zalo Bot SDK! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful and constructive in all interactions.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/go-zalo-bot.git
   cd go-zalo-bot
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/vkhangstack/go-zalo-bot.git
   ```
4. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.20 or higher
- Git
- golangci-lint (for linting)
- A Zalo Bot token for testing (optional)

### Install Dependencies

The SDK uses only the Go standard library, so no external dependencies are required:

```bash
go mod download
```

### Install Development Tools

Install golangci-lint for code quality checks:

```bash
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Windows
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Automation Scripts

The project includes automation scripts in the `scripts/` directory to streamline development and release workflows.

### Running Local Scripts

#### Build Script

Validates compilation and dependencies:

```bash
# Run build validation
./scripts/build.sh

# Run with verbose output
./scripts/build.sh --verbose

# Show help
./scripts/build.sh --help
```

The build script:
- Verifies Go version (minimum 1.20)
- Validates module dependencies
- Compiles all packages
- Compiles all examples
- Reports build status and timing

#### Test Script

Executes the comprehensive test suite:

```bash
# Run all tests
./scripts/test.sh

# Run with coverage report
./scripts/test.sh --coverage

# Run with race detector
./scripts/test.sh --race

# Run with verbose output
./scripts/test.sh --verbose

# Combine options
./scripts/test.sh --coverage --race --verbose
```

The test script:
- Runs all unit tests
- Generates coverage reports (HTML and terminal)
- Runs race detector when requested
- Displays coverage percentage
- Fails if coverage drops below 80%

#### Lint Script

Performs static analysis and code quality checks:

```bash
# Run linting
./scripts/lint.sh

# Auto-fix issues where possible
./scripts/lint.sh --fix

# Run with verbose output
./scripts/lint.sh --verbose
```

The lint script checks:
- Code formatting with `go fmt`
- Suspicious constructs with `go vet`
- Comprehensive linting with `golangci-lint`
- Missing godoc comments on exported symbols

#### Validate Script

Pre-release validation checks:

```bash
# Run validation
./scripts/validate.sh

# Run with verbose output
./scripts/validate.sh --verbose
```

The validate script verifies:
- All examples compile successfully
- CHANGELOG.md format is correct
- README.md has installation instructions
- All exported functions have godoc comments
- go.mod version matches latest Git tag

### Creating Releases

Use the release script to automate version management:

```bash
# Create a new release
./scripts/release.sh v1.2.3

# Add release message
./scripts/release.sh v1.2.3 --message "Bug fixes and improvements"

# Dry-run to preview actions
./scripts/release.sh v1.2.3 --dry-run

# Show help
./scripts/release.sh --help
```

The release script:
1. Validates semantic version format (v*.*.*)
2. Checks for duplicate tags
3. Runs full validation suite (build, test, lint, validate)
4. Updates CHANGELOG.md with new version section
5. Creates Git tag with version
6. Pushes tag to remote (triggers CI/CD deployment)
7. Displays next steps and monitoring instructions

**Important**: The release script will automatically trigger the CI/CD deployment workflow when the tag is pushed to GitHub.

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./types
go test ./services
```

### Run Examples

```bash
# Set your bot token
export ZALO_BOT_TOKEN="your-bot-token"

# Run polling example
go run examples/polling/main.go

# Run webhook example
export ZALO_WEBHOOK_SECRET="your-webhook-secret"
go run examples/webhook/main.go
```

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug Fixes**: Fix issues and improve stability
- **New Features**: Add new functionality to the SDK
- **Documentation**: Improve or add documentation
- **Examples**: Create new examples or improve existing ones
- **Tests**: Add or improve test coverage
- **Performance**: Optimize code for better performance
- **Refactoring**: Improve code quality and maintainability

### Before You Start

1. Check existing issues and pull requests to avoid duplicates
2. For major changes, open an issue first to discuss your proposal
3. Make sure your changes align with the project's goals and design

## Coding Standards

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html) guidelines.

### Code Formatting

- Use `gofmt` to format your code:
  ```bash
  gofmt -w .
  ```
- Use `go vet` to check for common mistakes:
  ```bash
  go vet ./...
  ```

### Naming Conventions

- Use descriptive names for variables, functions, and types
- Follow Go naming conventions (camelCase for unexported, PascalCase for exported)
- Use meaningful package names (lowercase, no underscores)

### Code Organization

- Keep functions small and focused
- Group related functionality together
- Use interfaces for abstraction where appropriate
- Avoid circular dependencies

### Error Handling

- Always handle errors explicitly
- Use typed errors (`ZaloBotError`) for SDK-specific errors
- Provide descriptive error messages
- Use error wrapping with `fmt.Errorf` and `%w`

### Comments and Documentation

- Add GoDoc comments for all exported types, functions, and methods
- Use complete sentences in comments
- Explain the "why" not just the "what"
- Include examples in documentation where helpful

Example:

```go
// SendMessage sends a text message to the specified chat.
// It validates the message configuration and returns the sent message
// or an error if the operation fails.
//
// The message text supports Unicode characters including Vietnamese.
// The chat ID must be a valid user or group identifier.
//
// Example:
//
//	message, err := bot.SendMessage(types.MessageConfig{
//	    ChatID: "user123",
//	    Text:   "Hello, World!",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
func (b *BotAPI) SendMessage(config types.MessageConfig) (*types.Message, error) {
    // Implementation
}
```

## Testing Guidelines

### Writing Tests

- Write tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Test both success and failure cases
- Test edge cases and boundary conditions

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "invalid input",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Coverage

- Aim for at least 80% test coverage
- Focus on testing critical paths and error handling
- Use `go test -cover` to check coverage

## Documentation

### Types of Documentation

1. **Code Documentation**: GoDoc comments for all exported items
2. **README**: High-level overview and quick start guide
3. **Examples**: Working code examples in the `examples/` directory
4. **Package Documentation**: Package-level documentation in `doc.go`
5. **CHANGELOG**: Document all changes in CHANGELOG.md

### Documentation Standards

- Keep documentation up-to-date with code changes
- Use clear, concise language
- Include code examples where helpful
- Document parameters, return values, and errors
- Explain complex logic or algorithms

## Pull Request Process

### Before Submitting

1. Run the validation scripts:
   ```bash
   # Run build validation
   ./scripts/build.sh
   
   # Run tests with coverage
   ./scripts/test.sh --coverage
   
   # Run linting (auto-fix if needed)
   ./scripts/lint.sh --fix
   
   # Run pre-release validation
   ./scripts/validate.sh
   ```

2. Alternatively, run individual commands:
   ```bash
   # Ensure all tests pass
   go test ./...
   
   # Format your code
   gofmt -w .
   
   # Check for common issues
   go vet ./...
   ```

3. Update documentation if needed
4. Add or update tests for your changes
5. Update CHANGELOG.md with your changes

### Submitting a Pull Request

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
2. Open a pull request on GitHub
3. Fill out the pull request template
4. Link any related issues
5. Wait for review and address feedback

### Pull Request Guidelines

- Keep pull requests focused and atomic
- Write clear, descriptive commit messages
- Reference related issues in the PR description
- Respond to review comments promptly
- Be open to feedback and suggestions

### Commit Message Format

Use clear, descriptive commit messages:

```
type: brief description

Detailed explanation of the changes (if needed)

Fixes #123
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

Example:
```
feat: add support for video messages

Implement SendVideo method to support sending video messages
with MIME type validation and proper error handling.

Fixes #45
```

## Reporting Issues

### Before Reporting

1. Check if the issue already exists
2. Verify you're using the latest version
3. Gather relevant information (Go version, OS, error messages)

### Issue Template

When reporting an issue, include:

- **Description**: Clear description of the problem
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Expected Behavior**: What you expected to happen
- **Actual Behavior**: What actually happened
- **Environment**: Go version, OS, SDK version
- **Code Sample**: Minimal code to reproduce the issue
- **Error Messages**: Full error messages and stack traces

### Feature Requests

For feature requests, include:

- **Use Case**: Why this feature is needed
- **Proposed Solution**: How you envision the feature working
- **Alternatives**: Other solutions you've considered
- **Additional Context**: Any other relevant information

## CI/CD Workflows

The project uses GitHub Actions for continuous integration and deployment.

### Continuous Integration (CI)

The CI workflow runs automatically on:
- Pull requests to the main branch
- Pushes to the main branch

**What the CI workflow does:**
- Tests against multiple Go versions (1.20, 1.21, 1.22)
- Tests on multiple operating systems (Ubuntu, macOS, Windows)
- Runs all tests with race detector and coverage
- Performs linting with golangci-lint
- Runs validation checks

**Required Status Checks:**
All CI jobs must pass before a pull request can be merged. If any job fails:
1. Review the error logs in the GitHub Actions tab
2. Fix the issues locally using the automation scripts
3. Push the fixes to your branch
4. The CI will automatically re-run

### Release and Deployment (CD)

The CD workflow runs automatically when a version tag is pushed (e.g., `v1.2.3`).

**What the CD workflow does:**
1. Validates the release (runs tests, validation, build)
2. Creates a GitHub Release with changelog notes
3. Notifies pkg.go.dev to index the new version
4. Verifies the package appears on pkg.go.dev

**To trigger a release:**
```bash
# Use the release script (recommended)
./scripts/release.sh v1.2.3 --message "Release notes"

# Or manually create and push a tag
git tag v1.2.3
git push origin v1.2.3
```

**Monitoring releases:**
- Check the GitHub Actions tab for workflow status
- Verify the release appears in the Releases section
- Confirm the package is indexed on pkg.go.dev (may take a few minutes)

### Workflow Requirements

For successful CI/CD execution:
- All tests must pass
- Code coverage must be at least 80%
- No linting errors
- All examples must compile
- CHANGELOG.md must be properly formatted
- Version tags must follow semantic versioning (v*.*.*)

## Troubleshooting

### Common Script Issues

#### Build Script Fails

**Issue**: Go version too old
```
Error: Go version 1.19 is below minimum required version 1.20
```
**Solution**: Upgrade Go to version 1.20 or higher

**Issue**: Module dependencies out of sync
```
Error: go.mod file is not tidy
```
**Solution**: Run `go mod tidy` to clean up dependencies

#### Test Script Fails

**Issue**: Coverage below threshold
```
Error: Coverage 75.5% is below threshold 80%
```
**Solution**: Add more tests to increase coverage

**Issue**: Race condition detected
```
WARNING: DATA RACE
```
**Solution**: Fix the race condition in your code. Use `go run -race` to reproduce locally.

#### Lint Script Fails

**Issue**: golangci-lint not installed
```
Error: golangci-lint is not installed
```
**Solution**: Install golangci-lint following the instructions in Development Setup

**Issue**: Formatting issues
```
Error: Files are not formatted with gofmt
```
**Solution**: Run `./scripts/lint.sh --fix` to auto-format

**Issue**: Missing godoc comments
```
Warning: Exported function 'FunctionName' is missing godoc comment
```
**Solution**: Add godoc comments to all exported functions, types, and methods

#### Validate Script Fails

**Issue**: Example compilation fails
```
Error: Example 'examples/polling/main.go' failed to compile
```
**Solution**: Fix the example code to ensure it compiles

**Issue**: CHANGELOG.md format invalid
```
Error: CHANGELOG.md does not contain version section
```
**Solution**: Ensure CHANGELOG.md follows the Keep a Changelog format

#### Release Script Fails

**Issue**: Invalid version format
```
Error: Version must follow semantic versioning format (v*.*.*)
```
**Solution**: Use proper semantic version format (e.g., v1.2.3, not 1.2.3 or v1.2)

**Issue**: Duplicate tag
```
Error: Tag v1.2.3 already exists
```
**Solution**: Use a different version number or delete the existing tag if it was created in error

**Issue**: Validation fails before release
```
Error: Pre-release validation failed
```
**Solution**: Fix all validation issues before creating a release. Run `./scripts/validate.sh --verbose` for details.

### CI/CD Issues

#### CI Workflow Fails

**Issue**: Tests fail on specific OS or Go version
```
Error: Test failed on windows-latest with Go 1.20
```
**Solution**: Test locally with the failing configuration or use GitHub Actions to debug

**Issue**: Timeout during test execution
```
Error: The job running on runner has exceeded the maximum execution time
```
**Solution**: Optimize slow tests or increase timeout in the workflow file

#### CD Workflow Fails

**Issue**: pkg.go.dev not indexing
```
Error: Package not found on pkg.go.dev after 30 seconds
```
**Solution**: Wait a few more minutes. pkg.go.dev indexing can take time. Verify manually at https://pkg.go.dev/github.com/vkhangstack/go-zalo-bot

**Issue**: GitHub Release creation fails
```
Error: Resource not accessible by integration
```
**Solution**: Ensure the GitHub token has proper permissions. Check repository settings.

### Getting Help

If you encounter issues not covered here:
1. Check the script output for detailed error messages
2. Run scripts with `--verbose` flag for more information
3. Review the GitHub Actions logs for CI/CD issues
4. Search existing issues on GitHub
5. Open a new issue with details about the problem

## Questions and Support

- **Documentation**: Check the README and examples first
- **Issues**: Open an issue for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in the project's README and release notes.

Thank you for contributing to the Go Zalo Bot SDK! 🎉
