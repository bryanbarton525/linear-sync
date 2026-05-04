# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a mechanical implementation task requiring verbatim transcription of 9 provided Go source files into the workspace, followed by build validation, test execution, and Git commit to the main branch. No design decisions are required—the entire solution is pre-defined in the constitution. The primary technical challenge is ensuring character-for-character fidelity when writing files and correctly handling the workflow branch to main branch merge.

## Delivery Target

github.com/bryanbarton525/linear-sync main branch with all 9 source files committed

## Tech Stack

- Go 1.21
- lib/pq (PostgreSQL driver) v1.10.9
- stretchr/testify v1.9.0
- Git
- Linear.app GraphQL API
- PostgreSQL

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source Files | 9 Go source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) that implement a Linear.app to PostgreSQL sync service with environment-based configuration, GraphQL API client, transactional storage layer, and comprehensive test coverage. | Verbatim file content from constitution | 9 files written to workspace root |
| Build System | Go 1.21 toolchain executing 'go build -buildvcs=false .' to compile the service binary. The -buildvcs=false flag prevents version control metadata embedding, allowing builds in environments without full Git history. | All 9 source files in workspace | Compiled binary named 'linear-sync' |
| Test Suite | Go test runner executing 'go test -buildvcs=false . -v' to validate config parsing, Linear API client mocking, and storage layer behavior. 4 tests pass (config and Linear client), 2 tests skip (storage tests require PostgreSQL). | Compiled source and test files | Test results: 4 PASS, 2 SKIP |
| Version Control Integration | Git workflow to commit all 9 files to the main branch. The workspace starts on workflow/8e569ed1-0092-4724-bb1f-5f056debc54b but the constitution requires committing to main. Strategy: checkout main, apply files, commit, and push. | 9 source files validated by build and test | Commit on main branch at github.com/bryanbarton525/linear-sync |

## Architectural Decisions

1. **Use write_file tool to write all 9 files in a single response**
   - Rationale: The constitution requires verbatim file content with no modifications. The write_file tool preserves exact byte sequences without formatting or whitespace normalization. Writing all files in parallel maximizes efficiency.
   - Tradeoffs: Single-response bulk write vs sequential writes. Bulk approach is faster and ensures atomic completion of the file writing phase.
2. **Checkout main branch before committing**
   - Rationale: The constitution explicitly requires 'commit all source files to the main branch'. The workspace is on a workflow branch. Checking out main and committing directly satisfies the requirement without creating PRs or merge complexity.
   - Tradeoffs: Direct commit to main vs PR workflow. Direct commit matches the constitution and avoids the out-of-scope PR process.
3. **Use exact build and test commands from constitution**
   - Rationale: The constitution specifies 'go build -buildvcs=false .' and 'go test -buildvcs=false . -v' with exact flags. The -buildvcs=false flag is critical for builds in workflow environments without full Git metadata.
   - Tradeoffs: No tradeoffs—these are immutable requirements from the constitution.
4. **Skip go mod init, go mod tidy, and go work init**
   - Rationale: The constitution explicitly forbids running go mod commands because go.mod and go.sum are already provided with locked dependency versions. Running these commands would modify the files and violate the verbatim content requirement.
   - Tradeoffs: No tradeoffs—this is an explicit constraint from the constitution.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 9d711c7f | backend | Write 9 Go source files to workspace | - | Produce 9 code artifacts in the workspace root at /var/lib/go-orca/workspaces/8e569ed1-0092-4724-bb1f-5f056debc54b. Write each file verbatim with exact character-for-character fidelity to the content provided in the constitution. Use the write_file tool to create: (1) go.mod - module definition with dependencies, (2) go.sum - dependency checksums, (3) main.go - service entry point with signal handling and sync loop, (4) config.go - environment variable configuration loading, (5) linear.go - Linear API GraphQL client with Issue and Assignee types, (6) storage.go - PostgreSQL storage layer with transaction-based upsert, (7) config_test.go - table-driven config loading tests, (8) linear_test.go - httptest-based Linear client tests, (9) storage_test.go - PostgreSQL integration tests with skip behavior. Acceptance criteria: All 9 files exist in workspace root with byte-for-byte matching content to the constitution specifications. No whitespace normalization, no formatting changes, no additions or deletions. Verify each file path is correct (workspace root, not subdirectories). This task produces code artifacts only; it does not run builds or tests. |
| bcd2c36a | backend | Build Go service binary | 9d711c7f | Produce a compiled binary artifact by executing 'go build -buildvcs=false .' from the workspace root /var/lib/go-orca/workspaces/8e569ed1-0092-4724-bb1f-5f056debc54b. The build must complete with exit code 0. The -buildvcs=false flag is required to disable version control metadata embedding. Do NOT run 'go mod init', 'go mod tidy', or 'go work init'—the go.mod and go.sum are already provided and locked. Acceptance criteria: (1) Command 'go build -buildvcs=false .' exits with code 0, (2) A binary named 'linear-sync' (or platform equivalent) is created in the workspace root, (3) No build errors or warnings are emitted, (4) The build uses the exact dependency versions from go.mod without modification. If the build fails, capture the full error output and report it as a blocking issue. This task depends on all 9 source files being present in the workspace. |
| 07e40098 | backend | Execute test suite and verify results | bcd2c36a | Run the test suite by executing 'go test -buildvcs=false . -v' from the workspace root /var/lib/go-orca/workspaces/8e569ed1-0092-4724-bb1f-5f056debc54b. Capture the full test output. Verify the output contains exactly 4 PASS results and exactly 2 SKIP results. The 4 passing tests are: TestLoad (config_test.go with 4 subtests), TestFetchIssues (linear_test.go), and TestFetchIssues_APIError (linear_test.go). The 2 skipped tests are: TestStorage_Upsert and TestStorage_UpsertEmpty (storage_test.go) which skip when DATABASE_URL is unavailable. Acceptance criteria: (1) Command 'go test -buildvcs=false . -v' exits with code 0, (2) Output contains 'PASS: TestLoad' with 4 subtests, (3) Output contains 'PASS: TestFetchIssues' and 'PASS: TestFetchIssues_APIError', (4) Output contains exactly 2 lines with '--- SKIP:' for storage tests, (5) Final summary line shows 'PASS' and '4 passed, 2 skipped' or equivalent. If test results differ from expectations, report the full test output as a blocking issue. This task depends on successful build completion. |
| 7c3a83c1 | backend | Commit and push to main branch | 07e40098 | Commit all 9 source files to the main branch of github.com/bryanbarton525/linear-sync. The workspace is currently on branch 'workflow/8e569ed1-0092-4724-bb1f-5f056debc54b'. Strategy: (1) Run 'git checkout main' to switch to the main branch, (2) If main does not exist locally, run 'git checkout -b main origin/main' or 'git checkout -b main' if no remote tracking exists, (3) Copy or ensure all 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are staged, (4) Run 'git add go.mod go.sum main.go config.go linear.go storage.go config_test.go linear_test.go storage_test.go', (5) Commit with message 'Add Linear.app to PostgreSQL sync service', (6) Push to origin main with 'git push origin main'. Acceptance criteria: (1) All 9 files are committed to the main branch, (2) Commit message is descriptive, (3) Push succeeds without errors, (4) The main branch at github.com/bryanbarton525/linear-sync contains all 9 source files. If push fails due to permissions or branch protection, report the error as a blocking issue. This task depends on successful test validation. |

