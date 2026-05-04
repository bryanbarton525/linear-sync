# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A Go microservice that periodically fetches issues from Linear.app via GraphQL and persists them to PostgreSQL, enabling offline access and integration with external systems.

## Goals

- Establish a reliable sync mechanism that fetches Linear issues on a 5-minute interval
- Store Linear issues in PostgreSQL with upsert semantics to handle updates and new issues
- Provide a clean, testable API for Linear client and storage interactions
- Support graceful shutdown on SIGTERM/SIGINT signals
- Enable local testing without requiring a live PostgreSQL instance (tests skip gracefully)

## Constraints

- All 9 source files must be written exactly as provided; no modifications to content
- Must pass compilation with `go build -buildvcs=false .` without errors
- Must pass all non-skipped tests; storage tests may skip when PostgreSQL is unavailable (expected behavior)
- Dependencies are fixed in go.mod and go.sum; no `go mod tidy` or `go mod init` operations allowed
- Service operates in a single process; no clustering or horizontal scaling in this phase
- Configuration loaded exclusively from environment variables (LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL)
- Sync interval is hardcoded to 5 minutes; no configuration of polling frequency

## Audience

Backend engineers integrating Linear.app with PostgreSQL-backed systems; DevOps teams deploying the sync service.

## Output Medium

GitHub repository (github.com/bryanbarton525/linear-sync) with committed source code and executable binaries built via `go build`.

## Acceptance Criteria

- All 9 source files are written to the workspace without modification
- Source code compiles successfully: `go build -buildvcs=false .` produces no errors
- Test suite executes and reports 4 PASS results and 2 SKIP results (storage tests skip without PostgreSQL environment)
- No test failures or runtime panics during `go test -buildvcs=false . -v`
- Code is committed with commit message referencing the workflow task and includes Co-authored-by trailer for Copilot
- Final commit is pushed to github.com/bryanbarton525/linear-sync on the main branch
- Repository reflects the exact state of the 9 provided files with no drift

## Out of Scope

- Database schema migration tooling (tests create tables on-the-fly)
- API exposure (HTTP endpoints); service is pull-only via scheduled sync
- Authentication to GitHub beyond SSH key/token for git operations
- Linear API error recovery strategies (failed syncs log and continue on next interval)
- Metrics, tracing, or observability beyond structured logging
- Deployment manifests (Kubernetes, Docker Compose, systemd units)

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Load configuration from environment | Read LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables; fail startup with clear error if any are missing. | config.go, main.go |
| F2 | must | Establish PostgreSQL connection | Open and verify PostgreSQL database connection at startup via newDB(connStr); defer close on exit. | storage.go, main.go |
| F3 | must | Fetch issues from Linear GraphQL API | Query Linear GraphQL endpoint with team ID to retrieve issue list; decode response and extract issue metadata (id, title, description, state, priority, assignee, timestamps). | linear.go |
| F4 | must | Upsert issues to PostgreSQL | Insert or update issue records in linear_issues table using ON CONFLICT DO UPDATE; serialize assignee as JSON; handle empty issue lists gracefully. | storage.go |
| F5 | must | Schedule periodic sync on 5-minute interval | Use time.Ticker(5 * time.Minute) to trigger sync operations; execute doSync() on each tick and at startup. | main.go |
| F6 | must | Graceful shutdown on signals | Listen for SIGTERM and SIGINT; cancel context and exit cleanly without data loss. | main.go |
| F7 | should | Structured error logging | Log sync failures, API errors, and database errors with context-aware messages at ERROR level; log successful syncs at INFO level. | main.go, linear.go, storage.go |
| F8 | should | Sync timeout enforcement | Each sync operation has a 2-minute timeout to prevent hangs; context cancellation is propagated. | main.go |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Compilation without errors | All source code must compile with `go build -buildvcs=false .` producing a single executable; no build warnings or errors. | Go toolchain |
| NF2 | must | Test coverage and execution | Minimum 4 passing unit tests (config, linear_test x2, storage_upsert_empty); 2 storage tests skip without PostgreSQL. No test failures. | Test suite |
| NF3 | must | Context propagation in async operations | All blocking calls (HTTP, database) must accept and respect context.Context; context cancellation is checked in select statements. | Code review, qa-validation skill |
| NF4 | must | Error handling completeness | All error paths wrapped with context using fmt.Errorf("%w", err); no silently swallowed errors. | Code review |
| NF5 | must | Transaction safety in storage operations | Upsert operations use database transactions with rollback on error; prepared statements prevent SQL injection. | storage.go |
| NF6 | must | HTTP client resource cleanup | HTTP response bodies closed in defer blocks; test servers created with httptest.NewServer() are deferred closed. | linear.go, linear_test.go |
| NF7 | should | Idiomatic Go patterns | Code follows Go conventions: error as last return value, early return on error, interfaces at package boundaries, table-driven tests. | Code review |

## Dependencies

- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (testing assertions)
- Go 1.21+ runtime
- PostgreSQL 11+ (optional for testing; tests skip without it)
- Linear.app GraphQL API endpoint (requires valid API key)
- Git for commit and push operations

