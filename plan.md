# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

A stateless Go service that syncs Linear.app issues to PostgreSQL. All 9 source files are provided verbatim and must be written exactly as specified. The implementation follows a simple fetch-upsert loop with a 5-minute cadence, graceful shutdown via signal context, and PostgreSQL persistence with ON CONFLICT handling. No architectural decisions are needed — this is a pure file-writing and validation task.

## Delivery Target

github.com/bryanbarton525/linear-sync main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9
- github.com/stretchr/testify v1.9.0
- PostgreSQL

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| config | Reads LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL from environment; returns error if any missing. | environment variables | *Config struct |
| linear client | HTTP client that executes a GraphQL query against api.linear.app to fetch team issues, parsing the response into []Issue. | API key, team ID, context | []Issue |
| storage | PostgreSQL storage layer that upserts issues into linear_issues table using ON CONFLICT (id) DO UPDATE, serializing assignee as JSONB. | *sql.DB, []Issue, context | error |
| main service loop | Wires config, client, and storage; runs doSync() immediately then on a 5-minute ticker; exits on SIGINT/SIGTERM via signal.NotifyContext. | environment configuration | log output, persisted issues |

## Architectural Decisions

1. **Write all 9 files verbatim before executing any commands**
   - Rationale: The Matriarch and CRITICAL INSTRUCTIONS mandate file writing precedes all build/test commands.
   - Tradeoffs: No flexibility to adjust content; files must match specification exactly.
2. **No go mod operations permitted**
   - Rationale: go.mod and go.sum are pre-provided and must remain unchanged.
   - Tradeoffs: Dependencies are fixed; no version changes allowed.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 477c6289 | backend | Write all 9 source files verbatim | - | Write the following 9 files into the workspace at /var/lib/go-orca/workspaces/66cff7a8-3e73-49fd-a820-11751d4ad23c using the write_file tool. Each file must be written EXACTLY as provided — no additions, removals, or character changes.  Files to write: 1. go.mod — module linear-sync, go 1.21, requires github.com/lib/pq v1.10.9 and github.com/stretchr/testify v1.9.0 with indirect deps 2. go.sum — exact checksum file provided 3. main.go — package main; imports context, log, os, os/signal, syscall, time; implements main() with signal.NotifyContext, 5-minute ticker, doSync() closure calling client.fetchIssues and store.upsert 4. config.go — package main; Config struct with APIKey/TeamID/DatabaseURL; load() reads env vars LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL and errors if any missing 5. linear.go — package main; Issue and Assignee types; raw* internal decode types; linearClient with newClient() constructor; fetchIssues(ctx, teamID) sends GraphQL POST with Bearer auth and decodes response 6. storage.go — package main; Storage struct with *sql.DB; newDB(connStr) opens and pings postgres; newStorage(db) constructor; upsert(ctx, issues) runs a transaction with prepared INSERT...ON CONFLICT statement, marshaling Assignee to JSON 7. config_test.go — package main; TestLoad table-driven test with 4 subtests covering missing each env var and valid config; uses testify assert 8. linear_test.go — package main; TestFetchIssues uses httptest.NewServer to mock Linear API response and asserts parsed issues; TestFetchIssues_APIError asserts error propagation from API error response 9. storage_test.go — package main; TestStorage_Upsert and TestStorage_UpsertEmpty both skip when DATABASE_URL resolves to an unreachable PostgreSQL instance  Do NOT run any commands after writing files. Do NOT run go mod init, go mod tidy, or go work init. |
| 73de024a | backend | Build and test the service | 477c6289 | In the workspace at /var/lib/go-orca/workspaces/66cff7a8-3e73-49fd-a820-11751d4ad23c, after all 9 files have been written, execute the following commands in order:  1. Run: go build -buildvcs=false .    - Must exit with code 0 and produce a compiled executable.    - If the build fails, inspect errors carefully — do NOT modify any source files; report the issue.  2. Run: go test -buildvcs=false . -v    - Expected output: exactly 4 PASS and 2 SKIP.    - The 2 SKIP results come from TestStorage_Upsert and TestStorage_UpsertEmpty in storage_test.go when PostgreSQL is not available — this is correct and expected behavior.    - The 4 PASS results come from TestLoad (with 4 subtests) in config_test.go and TestFetchIssues and TestFetchIssues_APIError in linear_test.go.    - Do NOT treat SKIP as failure.  Do NOT run go mod init, go mod tidy, or go work init at any point. |
| 35229880 | ops | Commit and push to GitHub main branch | 73de024a | In the workspace at /var/lib/go-orca/workspaces/66cff7a8-3e73-49fd-a820-11751d4ad23c, after successful build and test validation, commit all 9 source files and push to the main branch of github.com/bryanbarton525/linear-sync.  Steps: 1. Ensure the git remote is set to https://github.com/bryanbarton525/linear-sync 2. Stage all files: git add go.mod go.sum main.go config.go linear.go storage.go config_test.go linear_test.go storage_test.go 3. Commit with a descriptive message: 'feat: implement Linear-to-PostgreSQL sync service' 4. Push to main branch: git push origin HEAD:main (or the configured primary branch)  The commit must include all 9 files. Do not amend or rebase existing history unless the push requires it for a clean merge. |

