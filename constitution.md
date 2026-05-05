# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A Go service that periodically syncs Linear.app issues to PostgreSQL, maintaining an up-to-date issue database by fetching from the Linear GraphQL API at regular intervals.

## Goals

- Fetch issues from Linear.app GraphQL API
- Store or update issues in PostgreSQL database
- Run continuously with 5-minute sync intervals
- Handle configuration via environment variables
- Provide comprehensive test coverage with graceful degradation

## Constraints

- Must be written in Go 1.21
- Must use PostgreSQL for persistence
- All 9 source files must be written verbatim without modification
- Must not run go mod init, go mod tidy, or go work init (dependencies already provided)
- Service must sync every 5 minutes
- Build command: go build -buildvcs=false .
- Test command: go test -buildvcs=false . -v
- Tests must show exactly 4 PASS and 2 SKIP (storage tests skip without PostgreSQL)
- Final code must be committed and pushed to github.com/bryanbarton525/linear-sync on main branch

## Audience

Go developers integrating with Linear.app and PostgreSQL databases

## Output Medium

Compiled Go binary and source code in GitHub repository

## Acceptance Criteria

- All 9 source files written exactly as provided with no character modifications
- go build -buildvcs=false . completes successfully with no errors
- go test -buildvcs=false . -v produces exactly 4 PASS and 2 SKIP
- Code compiled and pushed to github.com/bryanbarton525/linear-sync main branch
- Service starts without errors when all required environment variables are set
- Service gracefully handles shutdown on interrupt/SIGTERM signals
- Service logs sync operations and errors with appropriate severity levels

## Out of Scope

- Modifying go.mod or go.sum files
- Adding new dependencies not in go.mod
- Changing the 5-minute sync interval
- Database schema migration tooling
- API key rotation or secret management beyond environment variables

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Load Configuration from Environment | Read LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables; fail fast with descriptive errors if any are missing | config.go provides the load() function |
| F2 | must | Database Connection Management | Open PostgreSQL connection, verify connectivity via Ping, and provide cleanup via Close() | storage.go newDB() function |
| F3 | must | Fetch Issues from Linear GraphQL API | Query Linear.app GraphQL endpoint for team issues; parse response and handle API errors | linear.go fetchIssues() method |
| F4 | must | Upsert Issues to PostgreSQL | Insert or update issues in linear_issues table using ON CONFLICT clause; serialize Assignee to JSONB | storage.go upsert() method |
| F5 | must | Periodic Sync Loop | Start service, run initial sync, then repeat every 5 minutes using time.Ticker | main.go ticker loop |
| F6 | must | Graceful Shutdown | Listen for SIGINT and SIGTERM; cancel context and exit cleanly without data loss | main.go signal.NotifyContext and ctx.Done() select |
| F7 | should | Structured Logging | Log sync start, completion (with issue count), errors, and shutdown signals using log package | main.go log.Printf calls |
| F8 | must | Context-Aware Database Operations | All database operations must accept context.Context and support cancellation | storage.go upsert() and newDB() with context parameters |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Go 1.21 Compatibility | All source code must compile with Go 1.21 using go build -buildvcs=false | go.mod declares go 1.21 |
| NF2 | must | Test Coverage and Isolation | Unit tests must be in separate *_test.go files with same package name; achieve 4 PASS and 2 SKIP with standard testing package | config_test.go, linear_test.go, storage_test.go |
| NF3 | must | Error Handling and Wrapping | All errors must be wrapped with %w verb or fmt.Errorf to preserve call chain; no silent failures | Across all files with error returns |
| NF4 | must | Database Isolation in Storage Tests | Storage tests must skip gracefully if PostgreSQL is unavailable (t.Skip); tests must not interfere with production data | storage_test.go TestStorage_Upsert and TestStorage_UpsertEmpty |
| NF5 | must | No Go Module Reinitialization | Do not run go mod init, go mod tidy, or go work init; use provided go.mod and go.sum | Task constraint |
| NF6 | must | Clean Compilation with -buildvcs=false | Build must complete without warnings or errors using -buildvcs=false flag | Task requirement |

## Dependencies

- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (assertion library for tests)
- Standard library: context, database/sql, encoding/json, net/http, time, log, os

