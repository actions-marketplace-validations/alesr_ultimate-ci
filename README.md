# Ultimate Go CI Suite

[![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat&logo=github-actions&logoColor=white)](https://github.com/features/actions)
[![Go Reference](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Enterprise-grade GitHub Action for Go projects combining testing, coverage gates, security scanning, linting, and code quality checks into a single composable workflow.

## Features

✨ **Comprehensive Quality Assurance**
- 🧪 Parallel testing with race detector enabled by default
- 📊 Code coverage analysis with configurable gate enforcement
- 🔐 Security scanning via `govulncheck`
- 🎯 Strict linting with `golangci-lint` and custom configuration
- 🔄 Code modernization checks via `gomodernize`
- 📋 Module hygiene verification (`go mod tidy`)
- ⚡ Dual-layer caching (Go build cache + dependency cache)
- 📦 Automatic coverage report artifacts

## Quick Start

### Basic Usage

```yaml
name: Go CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Ultimate Go CI
        uses: alesr/ultimate-ci@v1
        with:
          min-coverage: '80'
```

### Monorepo with Subdirectory

For projects with Go modules in subdirectories, the action intelligently resolves paths:

```yaml
- name: Test API Package
  uses: alesr/ultimate-ci@v1
  with:
    working-directory: 'packages/api'
    go-version-file: 'go.mod'           # Action resolves to packages/api/go.mod
    min-coverage: '75'
    lint-timeout: '10m'

- name: Test CLI Package
  uses: alesr/ultimate-ci@v1
  with:
    working-directory: 'packages/cli'
    go-version-file: 'go.mod'           # Action resolves to packages/cli/go.mod
    min-coverage: '70'
```

**Smart Path Resolution:**
- If `go-version-file: 'go.mod'` and `working-directory: 'packages/api'`, the action automatically resolves to `packages/api/go.mod`
- If `go-version-file: 'packages/api/go.mod'` and `working-directory: 'packages/api'`, it still works correctly
- Cache paths are also intelligently resolved for subdirectories

### Strict Coverage Gate

```yaml
- uses: alesr/ultimate-ci@v1
  with:
    min-coverage: '90'           # Require 90% coverage
    race-detector: 'true'        # Enable race condition detection
    run-modernize: 'true'        # Check code modernization
    run-vulncheck: 'true'        # Run security scans
    verify-tidy: 'true'          # Ensure go.mod is tidy
```

### Minimal Setup (Fast CI)

```yaml
- uses: alesr/ultimate-ci@v1
  with:
    run-lint: 'false'            # Skip linting
    run-modernize: 'false'       # Skip modernization checks
    run-vulncheck: 'false'       # Skip security scans
    min-coverage: '0'            # No coverage gate
```

## Inputs

| Input | Description | Default | Required |
|-------|-------------|---------|----------|
| `working-directory` | Relative path to Go module directory | `.` | No |
| `go-version-file` | Path to go.mod or go.work file | `go.mod` | No |
| `cache-dependency-path` | Dependency file pattern for caching | `**/go.sum` | No |
| `run-modernize` | Enable gomodernize checks | `true` | No |
| `run-lint` | Enable golangci-lint checks | `true` | No |
| `lint-timeout` | Timeout for linter (e.g., `5m`) | `5m` | No |
| `lint-concurrency` | Linter thread count (0 = auto) | `0` | No |
| `run-tests` | Enable go test execution | `true` | No |
| `race-detector` | Enable `-race` flag | `true` | No |
| `test-timeout` | Timeout for tests (e.g., `10m`) | `10m` | No |
| `test-parallel` | Parallel test executions (empty = auto) | `` | No |
| `min-coverage` | Minimum coverage % (0 = disabled) | `0` | No |
| `run-vulncheck` | Enable govulncheck security scan | `true` | No |
| `vulncheck-timeout` | Timeout for vulncheck (e.g., `3m`) | `3m` | No |
| `verify-tidy` | Verify go.mod and go.sum are tidy | `true` | No |

## Outputs

| Output | Description | Example |
|--------|-------------|---------|
| `test-coverage-file` | Path to coverage profile | `coverage.out` |
| `test-coverage-percent` | Total code coverage percentage | `84.5%` |

### Using Outputs in Downstream Steps

```yaml
- name: Run Ultimate Go CI
  id: quality
  uses: alesr/ultimate-ci@v1
  with:
    min-coverage: '75'

- name: Report Coverage
  run: |
    echo "Coverage: ${{ steps.quality.outputs.test-coverage-percent }}"
    echo "Profile: ${{ steps.quality.outputs.test-coverage-file }}"

- name: Download Coverage Artifact
  uses: actions/download-artifact@v4
  with:
    name: coverage-report-${{ github.run_id }}
```

## Linting Configuration

The action uses a strict but realistic linting configuration defined in `.golangci.yml`:

**Enabled Linters:**
- `errcheck` - Unchecked errors
- `errorlint` - Error wrapping best practices
- `gosec` - Security issues
- `goimports` - Import organization
- `gofmt` - Code formatting
- `revive` - Best practices (replacement for golint)
- `prealloc` - Slice pre-allocation
- And 10+ others for code quality

**Configuration highlights:**
- Strict exported function documentation requirements
- Error wrapping and handling best practices
- Security scanning for hardcoded credentials
- Unused code detection

View the full configuration: [`.golangci.yml`](.golangci.yml)

## Real-World Examples

### Example 1: Standard Go Project

```yaml
name: CI

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Quality Checks
        uses: alesr/ultimate-ci@v1
        with:
          min-coverage: '80'
          lint-timeout: '5m'
          test-timeout: '10m'
```

### Example 2: Monorepo Structure

```
project/
├── go.mod
├── go.sum
├── packages/
│   ├── api/
│   │   ├── go.mod
│   │   └── go.sum
│   └── cli/
│       ├── go.mod
│       └── go.sum
```

```yaml
jobs:
  test-api:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/ultimate-ci@v1
        with:
          working-directory: 'packages/api'
          go-version-file: 'packages/api/go.mod'
          min-coverage: '75'

  test-cli:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/ultimate-ci@v1
        with:
          working-directory: 'packages/cli'
          go-version-file: 'packages/cli/go.mod'
          min-coverage: '70'
```

### Example 3: Pull Request with Artifact Upload

```yaml
name: PR Quality Checks

on: [pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Quality Checks
        id: qa
        uses: alesr/ultimate-ci@v1
        with:
          min-coverage: '85'

      - name: Comment PR with Coverage
        if: always()
        uses: actions/github-script@v7
        with:
          script: |
            const coverage = '${{ steps.qa.outputs.test-coverage-percent }}';
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `📊 Coverage: ${coverage}`
            });
```

## Coverage Artifacts

Coverage profiles are automatically uploaded as GitHub Artifacts for:
- Historical tracking
- Local analysis via `go tool cover`
- Integration with coverage badges (Codecov, Coveralls)

### Download Coverage Profile

```bash
# From GitHub CLI
gh run download <RUN_ID> -n coverage-report-<RUN_ID>

# Then analyze
go tool cover -html=coverage.out
```

## Performance Optimization

The action is optimized for speed:

- **Dual-layer caching**: Go build cache + Go module cache
- **Parallel testing**: Race detector with configurable concurrency
- **Conditional steps**: Disable unnecessary checks
- **Timeout configuration**: Prevent hanging workflows

**Typical execution time:**
- Simple project: 30-45 seconds
- Medium project: 1-2 minutes
- Large project with many tests: 2-5 minutes

## Troubleshooting

### Coverage Gate Failing

**Error**: `Code coverage (45.2%) is below the required threshold of 80%.`

**Solution**: Either increase test coverage or reduce the threshold:
```yaml
min-coverage: '50'  # or add tests to increase coverage
```

### Linter Errors

**Error**: `./main.go:10:1: exported function Foo should have comment`

**Solution**: Add documentation comments:
```go
// Foo prints "foo" to stdout.
func Foo() {
    fmt.Println("foo")
}
```

### Path Resolution Issues

**Understanding Smart Path Resolution:**

The action intelligently resolves file paths based on `working-directory`:

```yaml
# Scenario 1: Single module at root
working-directory: '.'
go-version-file: 'go.mod'  # Resolves to: ./go.mod

# Scenario 2: Module in subdirectory (monorepo)
working-directory: 'packages/api'
go-version-file: 'go.mod'  # Resolves to: packages/api/go.mod

# Scenario 3: Explicit full path (also works)
working-directory: 'packages/api'
go-version-file: 'packages/api/go.mod'  # Uses as-is
```

**Error**: `go.mod not found at: X`

**Solution**: Verify your working-directory and go-version-file are consistent:
```yaml
# ✅ Correct - short filename with working-directory
working-directory: 'packages/api'
go-version-file: 'go.mod'

# ✅ Also Correct - full path
working-directory: 'packages/api'
go-version-file: 'packages/api/go.mod'

# ❌ Wrong - mismatched directory and filename
working-directory: 'packages/api'
go-version-file: 'go.mod'  # ONLY if there's go.mod at repo root
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for local development setup and testing instructions.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Changelog

### v1.1.0 (Latest)
- ✅ Fixed syntax errors in golangci-lint args
- ✅ Added coverage artifact upload
- ✅ Improved subdirectory path validation
- ✅ Added strict `.golangci.yml` configuration
- ✅ Comprehensive documentation

### v1.0.0
- Initial release with core features
