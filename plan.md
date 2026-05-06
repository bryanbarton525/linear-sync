# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a mechanical file-writing and validation workflow with zero architectural decisions. All 9 Go source files are provided verbatim and must be written to the workspace in exact byte-for-byte form, then built, tested, and pushed to GitHub. The service is a Linear.app-to-PostgreSQL sync daemon with config, API client, storage layer, and test coverage.

## Delivery Target

Git repository at github.com/bryanbarton525/linear-sync, main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (testing utilities)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| File Writing Phase | Write all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to workspace verbatim without modification | Provided file contents from user request | 9 Go source files in workspace root |
| Build Validation Phase | Execute go build -buildvcs=false to compile the service binary and verify successful compilation | All 9 source files in workspace | linear-sync binary, exit code 0 |
| Test Validation Phase | Execute go test -buildvcs=false with verbose output to run test suite and verify expected results (4 PASS, 2 SKIP for storage tests) | Compiled service and test files | Test results showing 4 PASS and 2 SKIP |
| Repository Push Phase | Commit all 9 source files to workflow branch and push to github.com/bryanbarton525/linear-sync main branch | All validated source files | Committed and pushed code to GitHub main branch |

## Architectural Decisions

1. **Write all 9 files in a single task before any build/test commands**
   - Rationale: User explicitly requires ALL files to exist before running go commands. This prevents partial state and ensures go.mod/go.sum are present before Go toolchain operations.
   - Tradeoffs: Larger initial task but enforces critical execution order constraint
2. **Use -buildvcs=false flag for build and test commands**
   - Rationale: User explicitly specified this flag. It disables VCS stamping which may not be needed for this validation workflow.
   - Tradeoffs: Binary won't contain VCS metadata but follows explicit instructions
3. **Accept 2 SKIP test results as success criteria**
   - Rationale: Storage tests skip when PostgreSQL is unavailable, which is expected and correct behavior in environments without a running database.
   - Tradeoffs: Storage layer not fully validated but user explicitly states this is correct

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 6ad27c55 | backend | Write All 9 Source Files Verbatim to Workspace | - | Produce artifact kind `code`, names: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Write all 9 files to workspace root (/var/lib/go-orca/workspaces/cb5e88ce-5a8f-4951-9893-9e07c607d87d) using the write_file tool with exact content provided in user request. Each file must be byte-for-byte identical to the provided content - preserve all whitespace, newlines, and characters exactly. Do NOT modify, reformat, or add any content. Acceptance: All 9 files exist in workspace root with exact matching content, verified by reading each file back or checking file existence. |
| 08c45a4f | backend | Build the Linear-Sync Service Binary | 6ad27c55 | Execute shell command `go build -buildvcs=false .` in workspace root. This builds the linear-sync service binary from the 9 source files written in the previous task. Acceptance: Command exits with code 0, no error output, and produces a `linear-sync` binary in the workspace root. The -buildvcs=false flag is mandatory per user instructions. |
| dad27cd3 | backend | Run Test Suite and Verify Results | 08c45a4f | Execute shell command `go test -buildvcs=false . -v` in workspace root. This runs all tests with verbose output. Acceptance: Command exits with code 0 or non-zero with expected skip count. Test output MUST show exactly 4 tests passing (PASS) and exactly 2 tests skipping (SKIP). The 2 SKIP results are from storage_test.go tests that skip when PostgreSQL is unavailable - this is expected and correct behavior. Verify output contains 'PASS: TestLoad', 'PASS: TestFetchIssues', 'PASS: TestFetchIssues_APIError', 'PASS: TestStorage_UpsertEmpty', and 'SKIP' messages for TestStorage_Upsert and one other storage test. The -buildvcs=false flag is mandatory per user instructions. |
| 659ef81c | backend | Commit and Push to GitHub Main Branch | dad27cd3 | Commit all 9 source files to the workflow branch (workflow/cb5e88ce-5a8f-4951-9893-9e07c607d87d) using git add, git commit with message 'feat: add linear-sync service with Linear.app to PostgreSQL sync', then push to github.com/bryanbarton525/linear-sync targeting the main branch. Acceptance: All 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are committed, git push succeeds with no errors, and files are visible in the GitHub repository on the main branch at root level. |

