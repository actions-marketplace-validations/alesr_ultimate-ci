# Contributing to Ultimate Go CI

This document covers development setup, testing, and contribution guidelines.

## Development Setup

### Prerequisites

- Go 1.23 or later
- golangci-lint (for linting)
- gomodernize (for code modernization checks)
- govulncheck (for security scanning)
- act (optional, for local GitHub Actions testing)

### Installation

Clone the repository and install dependencies:

```bash
# Clone
git clone https://github.com/alesr/ultimate-ci.git
cd ultimate-ci

# Download module dependencies
go mod download
cd testdata && go mod download && cd ..

# Install tools (optional but recommended for local testing)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/golangci/gomodernize/cmd/gomodernize@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## Local Testing

### Using act

Simulate GitHub Actions locally:

```bash
act -j test-subdirectory
```

### Manual Testing

Run checks manually from the testdata directory:

```bash
cd testdata

# Run tests with race detector and coverage
go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Check coverage
go tool cover -func=coverage.out

# Run linter
golangci-lint run ./...

# Run modernization check
gomodernize check ./...

# Run security scan
govulncheck ./...

# Verify module tidiness
go mod tidy
git diff --exit-code go.mod go.sum
```

## Project Structure

```
ultimate-ci/
├── action.yml              # Main action definition
├── .golangci.yml           # Linting configuration
├── README.md               # User documentation
├── CONTRIBUTING.md         # This file
├── LICENSE                 # License
├── testdata/               # Test fixture Go module
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── main_test.go
└── .github/
    └── workflows/
        └── test.yml        # Self-test workflow
```

## Modifying the Action

### action.yml Changes

1. Test locally with `act`:
   ```bash
   act -j test-subdirectory
   ```

2. Validate YAML syntax:
   ```bash
   yamllint action.yml
   ```
   Or use an online validator.

3. Verify variable and command substitution syntax:
   - Use `${VAR}` for variable substitution
   - Use `$(...)` for command substitution
   - Quote strings properly
   - Ensure working-directory context is applied correctly

4. Test with different input combinations by modifying the test workflow.

### .golangci.yml Changes

1. Test against testdata:
   ```bash
   cd testdata
golangci-lint run --config=../.golangci.yml ./...
   ```

2. Update README.md if adding or removing linters.

### Test Fixture Changes

Keep testdata as a minimal but representative Go project that:
- Contains real tests with coverage
- Passes linting and security checks
- Demonstrates typical Go project structure

## Pre-Submission Checklist

Before opening a pull request:

- Action YAML syntax is valid (`action.yml` passes validation)
- All tests pass:
  ```bash
  cd testdata && go test -v -race ./...
  ```
- Linting passes:
  ```bash
  golangci-lint run testdata/...
  ```
- Coverage analysis works:
  ```bash
  go tool cover -html=coverage.out
  ```
- Module is tidy:
  ```bash
  go mod tidy && git diff --exit-code
  ```
- Documentation updated if changing inputs, outputs, or behavior

## Pull Request Guidelines

### Before Opening

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make changes
4. Test locally
5. Commit with clear messages following the convention below

### Commit Message Convention

Prefix commits with one of:
- `fix:` - Bug fixes and corrections
- `feat:` - New features
- `docs:` - Documentation updates
- `test:` - Test additions or modifications
- `refactor:` - Code refactoring without behavior changes
- `perf:` - Performance improvements
- `ci:` - CI or workflow changes

Example:
```
fix: resolve lint timeout syntax error

Remove incorrect escaping from linter args and add validation
for go-version-file path.

Fixes #42
```

### PR Description Template

```markdown
## Changes

Description of what changed and why.

## Testing

- Tested locally with `act`
- Verified with testdata
- All checks pass

## Related Issues

Closes #123
```

## Common Issues

### act not found

Install act:
```bash
# macOS
brew install act

# Other platforms: https://github.com/nektos/act#installation
```

### Tests fail locally but pass in GitHub Actions

- Verify Go version: `go version` (must be 1.23+)
- Install dependencies: `go mod download`
- Clear test cache: `go clean -testcache`

### golangci-lint hangs

Increase timeout:
```bash
golangci-lint run --timeout=10m testdata/...
```

Or check for filesystem issues in testdata.

### Coverage report is empty

Verify tests run with coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Questions

Open a GitHub issue with the `question` label or start a discussion.

## Code of Conduct

Be respectful and constructive.
