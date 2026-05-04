# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Build a Go service that syncs Linear.app issues to PostgreSQL and deploys the complete codebase to github.com/bryanbarton525/linear-sync on the main branch.

## Goals

- Write all 9 source files exactly as specified without modification
- Successfully build the service using go build with -buildvcs=false flag
- Pass all executable tests with expected outcome of 4 PASS and 2 SKIP
- Commit all changes to the workflow branch
- Push the committed code to the main branch on GitHub

## Constraints

- All 9 files must be written VERBATIM - no character additions, removals, or changes permitted
- Do not run go mod init, go mod tidy, or go work init - dependency files are already provided
- Build and test commands must use -buildvcs=false flag
- Storage tests are expected to skip without PostgreSQL - this is correct behavior
- Must use the provided go.mod and go.sum without modification

## Audience

Development team and CI/CD pipeline consuming the linear-sync service

## Output Medium

Git repository at github.com/bryanbarton525/linear-sync

## Acceptance Criteria

- All 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are written to the workspace
- Command `go build -buildvcs=false .` completes with exit code 0
- Command `go test -buildvcs=false . -v` produces exactly 4 PASS and 2 SKIP test results
- All changes are committed to the workflow branch with a descriptive commit message
- Commits are pushed to the main branch at github.com/bryanbarton525/linear-sync
- No compilation errors, syntax errors, or import errors exist in any file
- The service implements Linear.app GraphQL API integration with proper authentication
- The service implements PostgreSQL storage with UPSERT logic for issue synchronization

## Out of Scope

- Setting up or configuring PostgreSQL database
- Running the service in a live environment
- Executing the 2 skipped storage tests (require PostgreSQL)
- Modifying file content to add features or fix issues
- Creating additional files beyond the 9 specified
- Running go mod tidy or other dependency management commands
- Adding linting, formatting, or other toolchain steps

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write all 9 source files to workspace | Use write_file tool to create go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, and storage_test.go with exact content as specified | User directive to write exactly 9 files VERBATIM |
| F2 | must | Build the Go service | Execute `go build -buildvcs=false .` and verify successful compilation with exit code 0 | Build step requirement after writing all files |
| F3 | must | Run test suite | Execute `go test -buildvcs=false . -v` and verify output shows 4 PASS and 2 SKIP results | Test validation requirement with specific expected outcomes |
| F4 | must | Commit changes to workflow branch | Commit all 9 files with descriptive message including Co-authored-by trailer | Git workflow requirement for version control |
| F5 | must | Push to main branch on GitHub | Push committed changes to main branch at github.com/bryanbarton525/linear-sync | Deployment requirement to target repository |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Byte-for-byte file accuracy | Each of the 9 files must match the specification exactly without any character changes | CRITICAL instruction to write files VERBATIM |
| NF2 | must | Build without VCS information | Build must succeed using -buildvcs=false flag to bypass version control checks | Build configuration constraint |
| NF3 | should | Test execution performance | Test suite must complete within 2 minutes without hanging or timing out | Standard test performance expectation |
| NF4 | must | Correct repository targeting | Git push operations must target the correct repository URL and main branch | Deployment target specification |

## Dependencies

- Git repository github.com/bryanbarton525/linear-sync must exist and be accessible
- Workspace must have write permissions for creating source files
- Go toolchain must be available for build and test operations
- Git credentials must be configured for push operations to GitHub

