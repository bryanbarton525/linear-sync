# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

A Go microservice that periodically fetches issues from the Linear.app GraphQL API and upserts them into PostgreSQL. All 9 source files are fully specified and must be written verbatim to the workspace. The implementation follows idiomatic Go 1.21 patterns: context propagation, error wrapping, graceful shutdown via signal.NotifyContext, and transaction-safe upserts. The task graph writes files, builds, tests, then commits and pushes to the target repository.

## Delivery Target

github.com/bryanbarton525/linear-sync main branch

## Tech Stack

- Go 1.21
- PostgreSQL (lib/pq v1.10.9)
- github.com/stretchr/testify v1.9.0
- Linear GraphQL API

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Config (config.go) | Reads LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables; returns clear errors for missing values. | Environment variables | *Config struct |
| Linear API Client (linear.go) | Posts GraphQL queries to https://api.linear.app/graphql with Bearer auth; decodes issue nodes into []Issue. | API key, team ID, context.Context | []Issue |
| Storage (storage.go) | Opens and pings PostgreSQL via lib/pq; performs transactional INSERT ... ON CONFLICT upserts with JSONB assignee column. | *sql.DB, []Issue, context.Context | error |
| Main (main.go) | Wires config, client, and storage; runs initial sync then repeats every 5 minutes with 2-minute per-cycle timeout; handles SIGINT/SIGTERM gracefully. | Config, linearClient, Storage | Running service |
| Tests (config_test.go, linear_test.go, storage_test.go) | Table-driven unit tests for config loading, API fetch (using httptest), and storage upsert (skips without PostgreSQL). | Source packages | 4 PASS, 2 SKIP |

## Architectural Decisions

1. **Write all 9 files verbatim as specified**
   - Rationale: Files are fully designed and validated upstream; no modification needed.
   - Tradeoffs: Zero design latitude; correctness depends on exact file content.
2. **Use go build -buildvcs=false and go test -buildvcs=false flags**
   - Rationale: Workspace may not have full git history; -buildvcs=false prevents build failures from missing VCS metadata.
   - Tradeoffs: Build info embedded in binary will lack VCS stamp.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| bc130233 | backend | Write all 9 source files verbatim to workspace | - | Produce the following files in the workspace at /var/lib/go-orca/workspaces/7c0e7dfc-4307-4fde-946f-3f51d49f4d75 using the write_file tool. Write each file EXACTLY as specified — no added comments, no reformatting, no character changes.  Files to write: 1. go.mod — module linear-sync, go 1.21, requires github.com/lib/pq v1.10.9 and github.com/stretchr/testify v1.9.0 with indirect deps. 2. go.sum — exact checksum file provided; write verbatim. 3. main.go — package main; loads config, opens DB, creates storage and client, runs sync loop with 5-minute ticker and 2-minute per-cycle timeout, handles SIGINT/SIGTERM via signal.NotifyContext. 4. config.go — package main; Config struct with APIKey/TeamID/DatabaseURL; load() reads from env vars, returns error for any missing var. 5. linear.go — package main; Issue and Assignee types; linearClient with fetchIssues(ctx, teamID) posting GraphQL to Linear API with Bearer auth. 6. storage.go — package main; Storage struct wrapping *sql.DB; newDB() opens+pings postgres; upsert(ctx, issues) does transactional INSERT ... ON CONFLICT upsert with JSONB assignee. 7. config_test.go — package main; table-driven tests for load() covering missing API key, missing team ID, missing database URL, and valid config. 8. linear_test.go — package main; tests for fetchIssues using httptest.NewServer — happy path and API error path. 9. storage_test.go — package main; tests for upsert and empty-list upsert, each calling t.Skip if PostgreSQL is unavailable.  Acceptance criteria: All 9 files are present in the workspace with exact content as specified. |
| 25772718 | backend | Build the service | bc130233 | From the workspace directory /var/lib/go-orca/workspaces/7c0e7dfc-4307-4fde-946f-3f51d49f4d75, run:    go build -buildvcs=false .  Do NOT run go mod init, go mod tidy, or go work init — go.mod and go.sum are already provided and must not be modified.  Acceptance criteria: Command exits with status 0 and no error output. If it fails, report the exact compiler error and do not proceed. |
| a6536c7e | backend | Run tests and verify 4 PASS 2 SKIP | 25772718 | From the workspace directory /var/lib/go-orca/workspaces/7c0e7dfc-4307-4fde-946f-3f51d49f4d75, run:    go test -buildvcs=false . -v  Verify the output shows exactly: - 4 tests with status PASS: TestLoad (4 subtests count as part of TestLoad via t.Run), TestFetchIssues, TestFetchIssues_APIError — total 4 top-level or subtest PASSes as specified - 2 tests with status SKIP: TestStorage_Upsert and TestStorage_UpsertEmpty (these skip when PostgreSQL is unavailable — this is correct and expected behavior)  Acceptance criteria: go test exits with status 0. Output contains exactly 4 PASS results and 2 SKIP results. No FAIL results. If any test FAILs, report the full test output. |
| 0661d2cd | ops | Commit and push to github.com/bryanbarton525/linear-sync main branch | a6536c7e | From the workspace directory /var/lib/go-orca/workspaces/7c0e7dfc-4307-4fde-946f-3f51d49f4d75, perform the following git operations:  1. git add . 2. git commit -m 'Initial implementation of Linear.app sync service  Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>' 3. git push origin main  If the remote branch is named workflow/7c0e7dfc-4307-4fde-946f-3f51d49f4d75, push to main explicitly:   git push origin HEAD:main  If git push fails due to authentication or credential issues, report the exact error message and stop — do not attempt workarounds.  Acceptance criteria: All 9 files are committed and visible on the main branch of https://github.com/bryanbarton525/linear-sync. |

