# Contributing to Ultimate Go CI Suite

Thank you for your interest in contributing! This document explains how to set up your development environment and test changes locally.

## Development Setup

### Prerequisites

- **Go** 1.23 or later
- **golangci-lint** for linting
- **gomodernize** for code modernization checks
- **govulncheck** for security scanning
- **act** (optional) for local GitHub Actions testing

### Installation

```bash
# Clone the repository
git clone https://github.com/alesr/ultimate-ci.git
cd ultimate-ci

# Download dependencies
go mod download
cd testdata && go mod download && cd ..

# Install linters globally (optional but recommended)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/golangci/gomodernize/cmd/gomodernize@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## Local Testing

### Test the Action Locally

```bash
# Option 1: Using act (simulates GitHub Actions)
act -j test-composite-action

# Option 2: Run testdata manually
cd testdata

# Run tests
go test -v -race -coverprofile=coverage.out ./...

# Extract coverage
go tool cover -func=coverage.out

# Check linting
golangci-lint run ./...

# Check modernization
gomodernize check ./...

# Security scan
govulncheck ./...

# Verify module tidiness
go mod tidy
git diff --exit-code go.mod go.sum
```

### Run All Checks Manually

```bash
# From testdata directory
./scripts/ci.sh  # if available

# Or manually
go test -v -race -timeout=10m -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | grep total
golangci-lint run --timeout=5m ./...
gomodernize check ./...
timeout 3m govulncheck ./...
go mod tidy
```

## Code Changes

### File Structure

```
ultimate-ci/
├── action.yml           # Action definition (main logic)
├── .golangci.yml        # Linting configuration
├── README.md            # User documentation
├── CONTRIBUTING.md      # This file
├── testdata/            # Test fixture (minimal Go project)
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── main_test.go
└── .github/
    └── workflows/
        └── test.yml     # Self-testing workflow
```

### When Modifying `action.yml`

1. **Test locally first**: Use `act` to simulate the workflow
   ```bash
   act -j test-composite-action
   ```

2. **Validate YAML syntax**:
   ```bash
   # Using yamllint (if installed)
   yamllint action.yml
   
   # Or online at https://www.yamllint.com/
   ```

3. **Check for common issues**:
   - All variable substitutions use `${VAR}` syntax (not `\(VAR\)`)
   - Command substitutions use `$(...)` syntax
   - Proper quoting in bash scripts
   - Working-directory context applied correctly

4. **Test with different inputs**:
   ```bash
   # Modify testdata inputs in .github/workflows/test.yml
   # Then run: act -j test-composite-action
   ```

### When Modifying `.golangci.yml`

1. **Test against testdata**:
   ```bash
   cd testdata
   golangci-lint run --config=../.golangci.yml ./...
   ```

2. **Document changes** in the linting config
3. **Update README.md** if adding/removing linters

### When Modifying Test Fixtures

The `testdata/` directory should remain a minimal but representative Go project:
- ✅ Contains real tests with coverage
- ✅ Passes linting and security checks
- ✅ Demonstrates typical Go project structure

```go
// Example: add test coverage
package main

import "testing"

func TestFoo(t *testing.T) {
    // Test implementation
}
```

## Testing Checklist

Before submitting a pull request, ensure:

- ✅ Action syntax is valid (`action.yml` has no errors)
- ✅ All tests pass locally:
  ```bash
  cd testdata && go test -v -race ./...
  ```
- ✅ Linting passes:
  ```bash
  golangci-lint run testdata/...
  ```
- ✅ Coverage analysis works:
  ```bash
  go tool cover -html=coverage.out  # can open coverage report
  ```
- ✅ Module is tidy:
  ```bash
  go mod tidy && git diff --exit-code
  ```
- ✅ Documentation updated (if changing inputs/outputs)

## PR Guidelines

### Before Opening a PR

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Test locally with `act` or manual testing
5. Commit with clear messages:
   ```bash
   git commit -m "fix: resolve issue with linting args"
   git commit -m "feat: add support for go.work files"
   git commit -m "docs: clarify subdirectory usage"
   ```

### PR Description Template

```markdown
## Changes
Brief description of what changed and why.

## Testing
- [x] Tested locally with `act`
- [x] Verified with testdata/
- [x] All checks pass

## Related Issues
Closes #123
```

### Commit Message Convention

- **fix:** Bug fixes and corrections
- **feat:** New features
- **docs:** Documentation updates
- **test:** Test additions/modifications
- **refactor:** Code refactoring without behavior changes
- **perf:** Performance improvements
- **ci:** CI/workflow changes

Example:
```
fix: resolve golangci-lint timeout syntax error

- Remove \( and \) escaping from linter args
- Use proper YAML action args syntax
- Add validation step for go-version-file path

Fixes #42
```

## Common Issues and Solutions

### Issue: `act` not found

```bash
# Install act
brew install act  # macOS
# or: https://github.com/nektos/act#installation
```

### Issue: Tests fail locally but pass in GitHub Actions

- Check Go version: `go version` (should be 1.23+)
- Check if dependencies are installed: `go mod download`
- Clear cache: `go clean -testcache`

### Issue: golangci-lint hangs

- Increase timeout: `golangci-lint run --timeout=10m testdata/...`
- Check for file system issues in the testdata directory

### Issue: Coverage report is empty

```bash
# Make sure tests actually run and cover code:
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Questions?

Open a GitHub issue with the `question` label or start a discussion in the repository.

## Code of Conduct

Be respectful and constructive. We're building this together!
