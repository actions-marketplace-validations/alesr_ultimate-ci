# Ultimate Go CI

[![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat&logo=github-actions&logoColor=white)](https://github.com/features/actions)
[![Go Reference](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

GitHub Action for Go projects that runs quality checks in parallel with cached tooling.

> **Note**: Go 1.23+ is required (the documented minimum). All supported versions provide `go mod tidy -diff`, which the tidy check uses directly.

## Overview

The action executes the following checks concurrently:
- Linting via `golangci-lint`
- Code modernization via `modernize`
- Security scanning via `govulncheck`
- Module tidiness verification
- Test execution with coverage analysis

Tool versions are pinned per action release: `golangci-lint v2.15.0`, `modernize v0.52.0`, `govulncheck v1.9.0`. Tools are cached between runs for deterministic results without per-run compilation.

## Installation

There is no tagged release yet, so `alesr/ultimate-ci@v1` does not resolve today. Until v1.0.0 is published, pin a commit SHA (recommended for supply-chain safety) or track `main`:

```yaml
# pinned to a commit (recommended)
- uses: alesr/ultimate-ci@0c6c5969c5ab41d9eeb7d6486c48f433625795f9

# or track main
- uses: alesr/ultimate-ci@main
```

The `@v1` examples below take effect once v1.0.0 is released. The release process (see [Changelog](#changelog)) creates a `v1.0.0` tag and a movable `v1` tag, which is also what unlocks a GitHub Marketplace listing.

## Usage

### Basic

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/ultimate-ci@v1
        with:
          min-coverage: '80'
```

### Monorepo Support

For projects with Go modules in subdirectories, specify the working directory and version file:

```yaml
- uses: alesr/ultimate-ci@v1
  with:
    working-directory: 'packages/api'
    go-version-file: 'go.mod'
```

The action resolves paths intelligently. If `working-directory` is set, a relative `go-version-file` like `go.mod` resolves to `<working-directory>/go.mod`.

### Configuration Options

Enable or disable individual checks:

```yaml
- uses: alesr/ultimate-ci@v1
  with:
    run-lint: 'false'        # Disable linting
    run-modernize: 'false'   # Disable modernization
    run-vulncheck: 'false'   # Disable security scanning
    verify-tidy: 'true'      # Verify module tidiness
    min-coverage: '0'        # Disable coverage gate
```

## Inputs

| Input | Description | Default |
|-------|-------------|---------|
| `working-directory` | Relative path to the Go module directory | `.` |
| `go-version-file` | Path to go.mod or go.work file | `go.mod` |
| `go-version` | Explicit Go version (e.g., `1.26`). Overrides `go-version-file` when set. | `` |
| `cache-dependency-path` | Dependency file pattern for caching | `**/go.sum` |
| `run-modernize` | Run gomodernize check | `true` |
| `modernize-timeout` | Timeout for modernize check | `3m` |
| `run-lint` | Run golangci-lint | `true` |
| `lint-timeout` | Timeout for golangci-lint | `5m` |
| `lint-concurrency` | Parallel threads for golangci-lint (0 = auto) | `0` |
| `lint-config` | Path to golangci-lint config file (relative to your repository root, or absolute). If empty, golangci-lint auto-discovers from consumer repo. | `` |
| `run-tests` | Run go test | `true` |
| `race-detector` | Enable `-race` flag | `true` |
| `test-timeout` | Timeout for go test | `10m` |
| `test-parallel` | Parallel test executions (empty = auto) | `` |
| `test-verbose` | Enable verbose test output (`-v` flag) | `true` |
| `min-coverage` | Minimum coverage percentage. Must be a non-negative integer (e.g., `80` for 80%). Set to `0` to disable. | `0` |
| `run-vulncheck` | Run govulncheck security scan | `true` |
| `vulncheck-timeout` | Timeout for govulncheck | `3m` |
| `verify-tidy` | Verify go.mod and go.sum are tidy | `true` |
| `tidy-timeout` | Timeout for go mod tidy check | `1m` |

## Outputs

| Output | Description | Example |
|--------|-------------|---------|
| `test-coverage-file` | Path to coverage profile (prefixed with `working-directory` in monorepo mode) | `packages/api/coverage.out` |
| `test-coverage-percent` | Total code coverage percentage | `84.5%` |

Example usage of outputs:

```yaml
- name: Run Ultimate CI
  id: ci
  uses: alesr/ultimate-ci@v1
  with:
    min-coverage: '75'

- name: Use Coverage
  run: echo "Coverage: ${{ steps.ci.outputs.test-coverage-percent }}"
```

## Linting Configuration

By default, `golangci-lint` auto-discovers configuration from the consumer's repository. To point at a specific config committed in your repository, pass its path relative to your repository root:

```yaml
- uses: alesr/ultimate-ci@v1
  with:
    lint-config: '.golangci.yml'  # resolves to <your repo>/.golangci.yml
```

Relative paths resolve against your repository, not the action's directory. To use the action's own strict configuration instead, pass an absolute path: `lint-config: '${{ github.action_path }}/.golangci.yml'`.

The action's configuration enables linters such as `errcheck`, `errorlint`, `gosec`, `goimports`, `gofmt`, `revive`, `prealloc`, `unused`, `ineffassign`, `misspell`, `unconvert`, and `nolintlint`.

View the full configuration: [`.golangci.yml`](.golangci.yml)

## Examples

### Standard Project

```yaml
name: CI

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/ultimate-ci@v1
        with:
          min-coverage: '80'
          lint-timeout: '5m'
          test-timeout: '10m'
```

### Monorepo with Multiple Packages

```yaml
jobs:
  test-api:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/ultimate-ci@v1
        with:
          working-directory: 'packages/api'
          go-version-file: 'go.mod'
          min-coverage: '75'

  test-cli:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/ultimate-ci@v1
        with:
          working-directory: 'packages/cli'
          go-version-file: 'go.mod'
          min-coverage: '70'
```

### Pull Request with Coverage Comment

```yaml
name: PR Checks

on: [pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/ultimate-ci@v1
        id: ci
        with:
          min-coverage: '85'
      - uses: actions/github-script@v7
        if: always()
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: 'Coverage: ' + '${{ steps.ci.outputs.test-coverage-percent }}'
            });
```

## Coverage Artifacts

Coverage profiles are uploaded as GitHub Artifacts named `coverage-report-<run_id>`.

Download and analyze:

```bash
gh run download <RUN_ID> -n coverage-report-<RUN_ID>
go tool cover -html=coverage.out
```

## Performance

Checks run in parallel. Typical execution times:
- Simple project: 35-45 seconds (warm cache)
- Medium project: ~1 minute (warm cache)
- Large project: bounded by test suite duration

GitHub-hosted runners have 2-4 vCPUs. CPU-heavy checks contend, so parallelism provides 2-3x speedup over sequential execution.

Caching includes Go build cache, Go module cache, and pinned tool binaries. On warm cache, tool setup is effectively zero.

## Troubleshooting

### Coverage Gate Failing

**Error**: `Code coverage (45.2%) is below the required threshold of 80%.`

**Solution**: Increase test coverage or lower the threshold:

```yaml
min-coverage: '50'
```

### Linter Errors

**Error**: `./main.go:10:1: exported function Foo should have comment`

**Solution**: Add documentation:

```go
// Foo prints "foo" to stdout.
func Foo() {
    fmt.Println("foo")
}
```

### Path Resolution Issues

**Error**: `go.mod not found at: X`

**Solution**: Ensure `working-directory` and `go-version-file` are consistent.

Correct patterns:
```yaml
# Short filename with working-directory
working-directory: 'packages/api'
go-version-file: 'go.mod'  # Resolves to: packages/api/go.mod

# Full path
working-directory: 'packages/api'
go-version-file: 'packages/api/go.mod'
```

Incorrect pattern:
```yaml
# go.mod lives in packages/api, but no working-directory is set,
# so go-version-file resolves to ./go.mod at the repo root: not found
go-version-file: 'go.mod'
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and testing instructions.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Changelog

### v1.0.0 (unreleased)

No tag or release exists yet. To publish v1.0.0:

1. Merge the changes intended for the release.
2. Tag the release commit: `git tag v1.0.0 && git push origin v1.0.0`.
3. Create (or move) the major-version tag: `git tag -f v1 && git push -f origin v1`.
4. Publish a GitHub Release from `v1.0.0`; this also enables the Marketplace listing.

Changes:
- Parallel check execution (tidy, modernize, lint, vulncheck, tests)
- Pinned, cached tool binaries: golangci-lint v2.15.0, modernize v0.52.0, govulncheck v1.9.0
- Inline PR annotations via shipped problem matcher (monorepo-aware path normalization)
- Code modernization logic with go.mod version gating
- Coverage gates, GitHub Step Summary, and coverage artifacts
- Smart path resolution for monorepo subdirectories
- Dual Go version selection: explicit `go-version` or derived from go.mod
