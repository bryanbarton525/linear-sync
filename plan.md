# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This workflow implements a mechanical file-writing and validation task for a Linear.app-to-PostgreSQL sync service. All 9 source files are pre-written and provided verbatim—no design decisions are required. The implementation follows a strict validation pipeline: write files → build → test → commit/push. The service uses Go 1.21 with PostgreSQL driver (lib/pq) and structured error handling. Core architecture: config loader, Linear GraphQL client, PostgreSQL storage layer with upsert, and a main loop with 5-minute ticker and graceful shutdown.

## Delivery Target

Compiled Go binary and committed source code in github.com/bryanbarton525/linear-sync main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (test assertions)
- PostgreSQL 12+ (runtime dependency)
- Linear.app GraphQL API

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Configuration Module | Loads LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables with validation | Environment variables | Config struct with validated values |
| Linear API Client | GraphQL client for fetching issues from Linear.app. Handles request construction, response parsing, and error propagation | API key, Team ID, HTTP context | Slice of Issue structs with parsed fields |
| PostgreSQL Storage Layer | Manages database connection and upsert operations with transaction safety and prepared statements | Database connection string, Issue slice | Persisted issues in linear_issues table |
| Main Service Loop | Orchestrates periodic sync (5-minute ticker), signal handling for graceful shutdown, and operational logging | Config, Storage, Linear client | Running service with periodic sync and clean termination |
| Test Suite | Unit tests for config loading, Linear API mocking, and storage operations. Storage tests skip gracefully without PostgreSQL | Source files, Test fixtures | 4 PASS + 2 SKIP test results |

## Architectural Decisions

1. **Use verbatim file content without modification**
   - Rationale: All source files are pre-written and tested. Any deviation risks introducing compilation or test failures
   - Tradeoffs: No flexibility for improvements, but guarantees exact compliance with acceptance criteria
2. **Push to main branch (not workflow branch)**
   - Rationale: The request explicitly states 'push to main branch'. The workspace shows workflow/77e711cb-2688-4f21-9a5e-52b6f9bccb4e as current branch, but the user intent is clear: deliver to main
   - Tradeoffs: May require branch switch or force push depending on git state. Task description will handle both scenarios
3. **Storage tests must skip without PostgreSQL**
   - Rationale: Acceptance criteria explicitly states '2 SKIP' is expected and correct behavior when database is unavailable
   - Tradeoffs: Test coverage is reduced in CI environments without PostgreSQL, but this is intentional design
4. **Use -buildvcs=false flag for compilation**
   - Rationale: Disables VCS stamping to avoid build errors when git metadata is incomplete or unavailable
   - Tradeoffs: Binary loses version metadata, but ensures reliable build in any git state

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| bc4704d2 | backend | Write All Source Files to Workspace | - | Produce 9 code artifacts in the workspace root directory (/var/lib/go-orca/workspaces/77e711cb-2688-4f21-9a5e-52b6f9bccb4e). Write each file VERBATIM with exact content provided—do not add, remove, or change any characters.  Artifacts to write: 1. `go.mod` (artifact kind: code) — Module definition with dependencies github.com/lib/pq v1.10.9 and github.com/stretchr/testify v1.9.0 2. `go.sum` (artifact kind: code) — Checksums for all dependencies 3. `main.go` (artifact kind: code) — Entry point with signal handling, ticker loop, and sync orchestration 4. `config.go` (artifact kind: code) — Configuration loader reading LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL from environment 5. `linear.go` (artifact kind: code) — Linear.app GraphQL client with fetchIssues method 6. `storage.go` (artifact kind: code) — PostgreSQL storage layer with newDB and upsert functions 7. `config_test.go` (artifact kind: code) — Unit tests for config loading with 4 test cases 8. `linear_test.go` (artifact kind: code) — Unit tests for Linear client with httptest mock server 9. `storage_test.go` (artifact kind: code) — Integration tests for storage upsert with PostgreSQL (skip if unavailable)  Acceptance criteria: - All 9 files written to workspace root with exact content from provided specification - No additional files created - File content is byte-for-byte identical to provided text - Do NOT run go mod init, go mod tidy, or go work init - Use write_file tool for each file (can be parallelized) |
| edf451d7 | backend | Validate Go Build | bc4704d2 | Produce validation artifact kind `build_result` by compiling the Go source code. Execute `go build -buildvcs=false .` from the workspace directory (/var/lib/go-orca/workspaces/77e711cb-2688-4f21-9a5e-52b6f9bccb4e).  Acceptance criteria: - Command exits with status 0 (success) - No compilation errors in stderr - Binary file `linear-sync` is created in workspace root - Build uses only dependencies from go.mod (lib/pq, testify) - The -buildvcs=false flag must be present to disable VCS stamping  If build fails, report exact error message including file, line, and error description. Do not proceed to next task on failure. |
| 4f3eb5df | backend | Validate Test Suite | edf451d7 | Produce validation artifact kind `test_result` by executing the test suite. Run `go test -buildvcs=false . -v` from the workspace directory (/var/lib/go-orca/workspaces/77e711cb-2688-4f21-9a5e-52b6f9bccb4e).  Acceptance criteria: - Command completes successfully (exit code 0) - Test output shows exactly 4 PASS results:   * TestLoad (config_test.go) — validates all 4 test cases   * TestFetchIssues (linear_test.go) — validates successful API response parsing   * TestFetchIssues_APIError (linear_test.go) — validates error handling   * TestStorage_UpsertEmpty (storage_test.go) — validates empty slice handling (may skip if DB unavailable) - Test output shows exactly 2 SKIP results:   * TestStorage_Upsert (storage_test.go) — skips if PostgreSQL unavailable   * TestStorage_UpsertEmpty (storage_test.go) — skips if PostgreSQL unavailable - Total: 4 PASS + 2 SKIP (or 6 PASS + 0 SKIP if PostgreSQL is available) - No test failures (FAIL) allowed - The -buildvcs=false flag must be present  If test results do not match expected pattern, report exact test names and outcomes. Do not proceed to next task on failure. |
| 9c76e14e | backend | Commit and Push to GitHub | 4f3eb5df | Produce artifact kind `git_commit` by committing all source files and pushing to github.com/bryanbarton525/linear-sync main branch. Execute from workspace directory (/var/lib/go-orca/workspaces/77e711cb-2688-4f21-9a5e-52b6f9bccb4e).  Steps: 1. Check current branch with `git branch --show-current` 2. If not on main branch, switch: `git checkout main` or create: `git checkout -b main` 3. Stage all files: `git add .` 4. Commit with message: `git commit -m "Implement Linear.app sync service  Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"` 5. Push to remote: `git push origin main` (may need `-u` flag if branch is new)  Acceptance criteria: - Current branch is main after checkout step - All 9 source files are staged and committed - Commit message includes Co-authored-by trailer - Push completes successfully to github.com/bryanbarton525/linear-sync - Remote main branch contains all committed files  If git authentication fails, report the error and suggest credential configuration (SSH key or personal access token). If push is rejected (non-fast-forward), report conflict details. |

