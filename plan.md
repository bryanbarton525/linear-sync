# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a mechanical transcription task with zero design latitude. The architecture is fully specified: a Go service that polls Linear.app's GraphQL API every 5 minutes, fetches team issues, and upserts them into PostgreSQL. The service uses standard library patterns for HTTP, context propagation, graceful shutdown via signal handling, and transactional database writes. All 9 source files are provided verbatim and must be written character-for-character without modification.

## Delivery Target

Repository: github.com/bryanbarton525/linear-sync, Branch: main

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (test assertions)
- PostgreSQL (external dependency)
- Linear.app GraphQL API (external dependency)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Configuration Loader | Reads LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables with validation errors for missing values. | Environment variables | Config struct |
| Linear API Client | GraphQL HTTP client that constructs team-scoped issue queries, authenticates with Bearer token, and parses nested JSON responses into Issue structs. | API key, Team ID | []Issue |
| PostgreSQL Storage | Transactional upsert layer that uses ON CONFLICT DO UPDATE to merge issues; serializes assignee to JSONB; handles empty lists gracefully. | []Issue, Database connection | Upsert confirmation |
| Main Loop | 5-minute ticker loop with signal-based graceful shutdown; propagates context cancellation to all I/O operations; logs sync events. | Config, Database, Linear client | Running service process |

## Architectural Decisions

1. **Use http.DefaultClient for Linear API calls**
   - Rationale: The provided code uses the standard HTTP client without custom transport; no requirements for connection pooling or TLS configuration beyond defaults.
   - Tradeoffs: No custom timeouts or retry logic, but meets the functional requirements as specified.
2. **5-minute polling interval**
   - Rationale: Fixed in main.go as time.NewTicker(5 * time.Minute); balances API rate limits with reasonable sync freshness.
   - Tradeoffs: Not configurable without code change, but matches the provided specification exactly.
3. **Transaction-per-batch upsert**
   - Rationale: All issues in a sync cycle are written in a single transaction with rollback on error, ensuring consistency.
   - Tradeoffs: Large batches could hold locks longer, but provides atomicity for the sync operation.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| c842d927 | backend | Write All Source Files Verbatim | - | Produce artifact kind `code`, writing exactly 9 Go source files into the workspace. Write each file character-for-character as provided—no additions, deletions, or modifications. Files: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Use the write_file tool for each file in a single batch. Acceptance: All 9 files exist in the workspace with exact byte-for-byte match to the provided content. Do NOT run any go commands during this task—file writing only. |
| a927bf47 | backend | Build the Project | c842d927 | Produce artifact kind `executable`, running go build -buildvcs=false . in the workspace directory. Acceptance: Command exits with code 0, producing a binary named 'linear-sync' (or the module name). No compilation errors or warnings. Do NOT run go mod init, go mod tidy, or go work init. The go.mod and go.sum files are already present and must not be modified. If build fails, report the error immediately without attempting to fix source files. |
| 093a9694 | backend | Run Test Suite | a927bf47 | Produce artifact kind `test_results`, running go test -buildvcs=false . -v in the workspace directory. Acceptance: Command completes with exactly 4 PASS results and 2 SKIP results. The 2 storage tests will skip gracefully when PostgreSQL is unavailable—this is expected and correct behavior, NOT a failure. No FAIL results allowed. If test output does not match 4 PASS + 2 SKIP, report the discrepancy immediately without modifying source files. |
| ee546c45 | ops | Commit and Push to GitHub | 093a9694 | Produce artifact kind `git_commit`, initializing git repository if needed, staging all 9 source files, committing with message 'Initial commit: Linear.app sync service' and Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com> trailer, setting remote origin to https://github.com/bryanbarton525/linear-sync, and pushing to main branch. Acceptance: Workspace is a git repository, remote is set to the correct URL, all 9 files are committed, and main branch exists on GitHub with the commit visible. If workspace is already a git repo with existing commits, add files, commit, and push. If push fails due to authentication or network issues, report the error with guidance on manual resolution. |

