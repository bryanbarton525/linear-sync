# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A stateless Go service that syncs Linear.app issues to PostgreSQL at regular intervals, providing a persistent store for issue tracking and enabling custom reporting and analysis workflows.

## Goals

- Implement a Linear API client that fetches team issues via GraphQL
- Build a PostgreSQL storage layer that upserts issues with conflict handling
- Create a main service loop that syncs issues on a 5-minute cadence with graceful shutdown
- Validate all code builds without errors and passes the test suite
- Commit and push implementation to the primary branch at github.com/bryanbarton525/linear-sync

## Constraints

- All 9 source files must be written exactly as provided — no additions, removals, or modifications to character content
- Go version: 1.21 with existing go.mod and go.sum; no go mod init, go mod tidy, or go work init allowed
- Build command: `go build -buildvcs=false .`
- Test command: `go test -buildvcs=false . -v`; expected result is 4 PASS and 2 SKIP (storage tests skip without PostgreSQL — this is expected and correct)
- Database connectivity is optional at build time; storage tests gracefully skip if PostgreSQL is unavailable
- Service reads configuration from three environment variables: LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL

## Audience

DevOps engineers and developers who need to mirror Linear.app issues into PostgreSQL for offline analysis, reporting, or integration with other tools.

## Output Medium

Go source code (9 files), compiled executable, and a repository commit on the main branch.

## Acceptance Criteria

- All 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are written exactly as provided with no character modifications
- Compilation succeeds: `go build -buildvcs=false .` produces an executable with zero errors
- Test execution passes validation: `go test -buildvcs=false . -v` shows 4 PASS and 2 SKIP
- The 2 SKIP results occur in storage_test.go when PostgreSQL is not available in the test environment
- Code is committed with a descriptive message and pushed to github.com/bryanbarton525/linear-sync on the main branch
- No intermediate go mod commands are executed; go.mod and go.sum remain unchanged from their provided state

## Out of Scope

- Schema migrations, database initialization scripts, or DDL management
- Linear API authentication beyond Bearer token (OAuth flows, token refresh)
- HTTP server or REST API endpoints
- Monitoring, metrics, or observability beyond log output
- Configuration management beyond environment variables
- Database connection pooling or advanced tuning
- Backoff/retry logic for API or database failures beyond error logging
- Support for multiple teams or batch team synchronization
- GraphQL query optimization or pagination beyond current implementation

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write all 9 source files | Write go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, and storage_test.go using the write_file tool with exact verbatim content provided. | Task specification, CRITICAL INSTRUCTIONS section 1 |
| F2 | must | Build executable without errors | Execute `go build -buildvcs=false .` after all files are written. The build must complete with exit code 0 and produce a working executable. | Task specification, step 3 |
| F3 | must | Execute and validate test suite | Run `go test -buildvcs=false . -v` and verify output shows exactly 4 PASS and 2 SKIP. The 2 SKIP results are expected for storage tests when PostgreSQL is unavailable. | Task specification, step 4-5 |
| F4 | must | Commit and push to repository | After successful build and test validation, commit all files with a descriptive message and push to the main branch of github.com/bryanbarton525/linear-sync. | Task specification, step 6 |
| F5 | must | Load configuration from environment | Implement config.go to read LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables and return error if any are missing. | config.go specification |
| F6 | must | Fetch issues from Linear API | Implement linear.go with a linearClient that authenticates with Bearer token and executes a GraphQL query to fetch issues by team ID, parsing the response into Issue structs. | linear.go specification |
| F7 | must | Persist issues to PostgreSQL | Implement storage.go with an upsert operation that inserts or updates issues in the linear_issues table, handling assignee as JSON and using ON CONFLICT for idempotency. | storage.go specification |
| F8 | must | Main service loop with periodic sync | Implement main.go with a 5-minute ticker that calls doSync(), handles context cancellation gracefully, and logs sync results and errors. | main.go specification |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Code must compile | All Go source code must compile without errors or warnings using Go 1.21 with the specified dependencies. | Toolchain validation profile |
| NF2 | must | Tests must be isolated and idempotent | Test cases must not share state, must clean up after themselves, and must produce consistent results across multiple runs. Storage tests must gracefully skip if PostgreSQL is unavailable. | qa-validation skill: Test Isolation Rule |
| NF3 | must | Error handling must propagate context | All error returns must be wrapped with fmt.Errorf('%w', ...) to preserve error chain; context.Context must be passed to all blocking operations. | code-generation skill: Error Handling |
| NF4 | must | No pre-existing modifications to go.mod or go.sum | The provided go.mod and go.sum must not be modified during implementation; no go mod init, go mod tidy, or go work init commands are permitted. | CRITICAL INSTRUCTIONS section 2 |
| NF5 | must | Service must respect context cancellation | The main loop must honor signal.NotifyContext and context.WithTimeout to ensure graceful shutdown and operation timeouts. | code-generation skill: Concurrency |

## Dependencies

- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (testing assertions)

---

## Constitution Amendment — Cycle 2

Root cause: implementation defect. Prior pods claimed to write 9 Go source files but workspace validation confirms zero files exist in /var/lib/go-orca/workspaces/66cff7a8-3e73-49fd-a820-11751d4ad23c. This is an implementation defect in file-writing process — not a requirement gap, not a design gap, but a failure to execute the write_file tool correctly. Remediation: (1) Write all 9 files using write_file tool with absolute paths; (2) Verify file presence via directory listing; (3) Build with `go build -buildvcs=false .`; (4) Test with `go test -buildvcs=false . -v` (expect 4 PASS + 2 SKIP); (5) Commit and push to main branch. No architectural changes needed — all acceptance criteria remain valid. The constitution is complete and unmodified.
