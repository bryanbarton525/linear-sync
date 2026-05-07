# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

Verbatim file-writing workflow that creates a Linear.app-to-PostgreSQL sync service by writing 9 Go source files exactly as specified, building and testing the service, then deploying to GitHub main branch.

## Delivery Target

Git repository github.com/bryanbarton525/linear-sync main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (testing framework)
- Linear.app GraphQL API
- PostgreSQL (target persistence layer)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source Files | Nine Go source files comprising the service: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go | File specifications from original request | 9 Go source files in workspace root |
| Build Verification | Go build process that compiles source files without VCS metadata requirement | All 9 source files | Compiled binary, build success confirmation |
| Test Verification | Test suite execution validating config loading, Linear API client, and storage layer with expected skip behavior for PostgreSQL integration tests | Compiled binary, test source files | Test results: 4 PASS, 2 SKIP |
| Git Deployment | Version control operations committing all source files to workflow branch and pushing to main branch | All verified source files | Git commits, remote main branch updated |

## Architectural Decisions

1. **Write all 9 files in parallel using single task**
   - Rationale: Minimizes turn count and maximizes efficiency per Matriarch guidance
   - Tradeoffs: Single task failure blocks entire file-writing phase, but risk is low given verbatim specification
2. **Use -buildvcs=false flag for build and test**
   - Rationale: Workspace may not have full VCS metadata configured, flag ensures deterministic builds
   - Tradeoffs: Binary won't include VCS commit info, acceptable for this deployment pattern
3. **Accept 2 SKIP test results as success condition**
   - Rationale: Storage integration tests correctly skip when PostgreSQL unavailable in build environment
   - Tradeoffs: Integration tests not exercised in build phase, acceptable given they validate DB connection only
4. **Commit to workflow branch first, then push to main**
   - Rationale: Follows standard Git workflow pattern, preserves audit trail on workflow branch
   - Tradeoffs: Requires additional push operation, but provides rollback safety

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| c0558717 | backend | Write all 9 source files verbatim | - | Write exactly 9 Go source files to the workspace root (/var/lib/go-orca/workspaces/c7f92732-cc4d-4e98-ac43-fb5981fe9d96) with byte-for-byte accuracy matching the specifications in the original request.  Files to write: 1. go.mod - Module declaration with dependencies 2. go.sum - Cryptographic checksums for dependencies 3. main.go - Service entrypoint with signal handling and sync loop 4. config.go - Environment variable configuration loader 5. linear.go - Linear.app GraphQL API client 6. storage.go - PostgreSQL persistence layer with transactional upserts 7. config_test.go - Configuration loading unit tests 8. linear_test.go - API client unit tests with httptest 9. storage_test.go - Database integration tests (skip when PostgreSQL unavailable)  Acceptance criteria: - All 9 files written to workspace root with ZERO character modifications from specification - File content matches specification exactly including whitespace, newlines, and formatting - No additional files created - All files are valid Go source code (go.mod, go.sum, *.go, *_test.go)  Critical constraints: - Do NOT run go mod init, go mod tidy, or go work init - Do NOT modify go.mod or go.sum content - Write files using write_file tool with exact content from original request |
| 7f3cd388 | backend | Build the service | c0558717 | Execute the Go build command to compile the Linear sync service binary.  Command to execute from workspace root: go build -buildvcs=false .  Acceptance criteria: - Build completes with exit code 0 - No compilation errors reported - Binary file 'linear-sync' created in workspace root - Build output confirms successful compilation of all source files  The -buildvcs=false flag disables VCS metadata embedding, ensuring deterministic builds without requiring full git history. |
| 285c3d91 | backend | Execute test suite and verify output | 7f3cd388 | Run the Go test suite with verbose output and verify the expected test result pattern.  Command to execute from workspace root: go test -buildvcs=false . -v  Acceptance criteria: - Test execution completes with exit code 0 - Exactly 4 tests PASS: TestLoad (4 subtests in config_test.go), TestFetchIssues, TestFetchIssues_APIError (linear_test.go) - Exactly 2 tests SKIP: TestStorage_Upsert, TestStorage_UpsertEmpty (storage_test.go skip when DATABASE_URL PostgreSQL unavailable) - Output pattern matches: '4 passed, 2 skipped' or equivalent verbose output showing 4 PASS and 2 SKIP results - No test failures or panics  The 2 SKIP results are expected and correct behavior when PostgreSQL is not available in the test environment. |
| 320e1338 | ops | Commit and push to GitHub main branch | 285c3d91 | Commit all source files to the workflow branch and push to the main branch at github.com/bryanbarton525/linear-sync.  Git operations to perform: 1. Stage all 9 source files (go.mod, go.sum, *.go, *_test.go) 2. Commit with message: 'Add Linear.app to PostgreSQL sync service' 3. Push workflow branch to remote 4. Push to main branch at github.com/bryanbarton525/linear-sync  Acceptance criteria: - All 9 files committed to workflow branch workflow/c7f92732-cc4d-4e98-ac43-fb5981fe9d96 - Commit includes proper message and Co-authored-by trailer - Workflow branch pushed to remote successfully - Main branch at github.com/bryanbarton525/linear-sync updated with all source files - No merge conflicts or push errors - Remote repository contains working Linear sync service code  Note: Verify git credentials are configured before attempting push operations. |

