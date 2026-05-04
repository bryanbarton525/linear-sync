# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This workflow implements a Linear.app to PostgreSQL synchronization service by writing 9 pre-specified Go source files verbatim to the workspace, validating build and test success, and deploying to the target GitHub repository. The service features GraphQL API integration with Linear.app, database persistence with upsert conflict resolution, environment-based configuration, and graceful shutdown handling. The implementation is deterministic: all source files are provided in full, requiring no code generation or design decisions.

## Delivery Target

GitHub repository github.com/bryanbarton525/linear-sync on main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (test assertions)
- Linear.app GraphQL API
- PostgreSQL (runtime dependency)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| File Writing Phase | Write exactly 9 source files to workspace using write_file tool with no character modifications | Provided source file content in task specification | go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go |
| Build Validation Phase | Compile the service executable with -buildvcs=false flag to verify all dependencies resolve correctly | All 9 source files written in phase 1 | Compiled linear-sync executable |
| Test Validation Phase | Execute test suite and verify exactly 4 PASS and 2 SKIP results (storage tests skip without PostgreSQL) | Compiled service and test files | Test execution report showing 4 PASS + 2 SKIP |
| Repository Deployment Phase | Commit all source files and push to github.com/bryanbarton525/linear-sync main branch | All validated source files from prior phases | Git commit on main branch with all 9 files |

## Architectural Decisions

1. **Write all 9 files before running any build or test commands**
   - Rationale: Build and test commands depend on complete source tree; partial writes will fail compilation
   - Tradeoffs: Sequential file writing increases initial latency but ensures atomic workspace state before validation
2. **Use -buildvcs=false flag for build and test commands**
   - Rationale: Task specification explicitly requires this flag; prevents Git metadata embedding during build
   - Tradeoffs: Binary will not contain VCS information, but this is acceptable per requirements
3. **Accept 2 SKIP test results as success condition**
   - Rationale: Storage tests require PostgreSQL connection which is not available in build environment; skipping is expected behavior
   - Tradeoffs: Storage layer is not fully tested in CI, but this is acceptable per requirements
4. **Push directly to main branch without pull request**
   - Rationale: Workflow request specifies pushing to main branch; no PR workflow mentioned
   - Tradeoffs: Bypasses code review but satisfies deployment requirement

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 5def52be | backend | Write All Source Files to Workspace | - | Produce artifact kind `code`, name `linear-sync-sources`. Use write_file tool to create exactly 9 Go source files in the workspace directory /var/lib/go-orca/workspaces/716d63b2-35a0-4d3e-9711-13e2242d50ec: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Each file must be written VERBATIM with no additions, removals, or character changes from the provided content. Write ALL 9 files before proceeding to any build or test commands. Acceptance criteria: (1) All 9 files exist in workspace with exact content matching task specification, (2) go.mod declares module linear-sync with Go 1.21, (3) go.sum contains 7 dependency entries, (4) main.go contains signal-based shutdown and 5-minute sync ticker, (5) config.go loads three environment variables, (6) linear.go implements GraphQL client with Bearer auth, (7) storage.go implements PostgreSQL upsert with conflict resolution, (8) Three test files contain table-driven tests and httptest server usage. |
| 3b045949 | backend | Build Service Executable | 5def52be | Produce artifact kind `executable`, name `linear-sync`. From workspace directory /var/lib/go-orca/workspaces/716d63b2-35a0-4d3e-9711-13e2242d50ec, run command: go build -buildvcs=false . This compiles the Linear sync service without VCS metadata embedding. Do NOT run go mod init, go mod tidy, or go work init - the go.mod and go.sum files are already provided. Acceptance criteria: (1) Build command exits with status 0, (2) No compilation errors in stdout/stderr, (3) Executable file linear-sync exists in workspace directory, (4) All dependencies from go.mod resolve correctly without additional downloads. |
| 3b6b293b | backend | Execute Test Suite with Validation | 3b045949 | Produce artifact kind `test_report`, name `linear-sync-test-results`. From workspace directory /var/lib/go-orca/workspaces/716d63b2-35a0-4d3e-9711-13e2242d50ec, run command: go test -buildvcs=false . -v This executes all test files with verbose output. Acceptance criteria: (1) Test command exits with status 0, (2) Stdout shows exactly 4 tests with PASS status: TestLoad (4 subtests), TestFetchIssues, TestFetchIssues_APIError, (3) Stdout shows exactly 2 tests with SKIP status: TestStorage_Upsert and TestStorage_UpsertEmpty (skipping is expected when PostgreSQL is unavailable), (4) No test failures or panics reported. The 2 SKIP results are correct and expected behavior per task requirements. |
| 3fb2bd8f | ops | Commit and Push to GitHub Main Branch | 3b6b293b | Produce artifact kind `git_commit`, name `linear-sync-deployment`. From workspace directory /var/lib/go-orca/workspaces/716d63b2-35a0-4d3e-9711-13e2242d50ec, commit all 9 source files to Git and push to github.com/bryanbarton525/linear-sync main branch. Use commit message: 'Implement Linear.app to PostgreSQL sync service' followed by the required Co-authored-by trailer: 'Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>'. Ensure git is configured for the target repository and branch. Acceptance criteria: (1) git add command includes all 9 files, (2) git commit succeeds with specified message and trailer, (3) git push to origin main succeeds without authentication or branch protection errors, (4) All files are visible in github.com/bryanbarton525/linear-sync main branch after push, (5) Commit history shows the implementation commit with Co-authored-by trailer. |

