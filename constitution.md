# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A Go microservice that periodically synchronizes Linear.app issues to PostgreSQL, enabling offline access, data warehousing, and integration with downstream systems.

## Goals

- Fetch issues from Linear.app GraphQL API every 5 minutes
- Upsert issues into PostgreSQL with conflict-on-duplicate-key semantics
- Load configuration from environment variables (API key, team ID, database URL)
- Provide graceful shutdown on SIGTERM and SIGINT signals
- Log operational events and errors with structured messages

## Constraints

- Must compile with Go 1.21 or later
- Must use PostgreSQL as the sole data store (via lib/pq driver)
- Must pass context.Context to all blocking I/O operations for proper cancellation
- All configuration must come from environment variables; no config files
- Sync operations must complete within 2 minutes (timeout enforced per cycle)
- Linear API authentication must use Bearer token in Authorization header

## Audience

Backend infrastructure teams deploying self-hosted issue tracking and data warehousing systems; platform engineers integrating Linear.app with their data lake

## Output Medium

Go executable and library code pushed to github.com/bryanbarton525/linear-sync on the main branch; deployable as a containerized service or local daemon

## Acceptance Criteria

- go build -buildvcs=false . completes without errors
- go test -buildvcs=false . -v produces exactly 4 PASS results
- go test -buildvcs=false . -v produces exactly 2 SKIP results (storage tests skip gracefully when PostgreSQL unavailable)
- All 9 source files are committed and pushed to github.com/bryanbarton525/linear-sync main branch
- Configuration loader rejects missing LINEAR_API_KEY, LINEAR_TEAM_ID, or DATABASE_URL with clear error messages
- Issues table in PostgreSQL uses id as PRIMARY KEY and correctly performs UPSERT on conflict
- Service fetches and syncs issues every 5 minutes with individual sync operations bounded to 2 minutes
- Service catches SIGTERM and SIGINT, logs shutdown message, and exits cleanly
- All error paths include wrapped error context (no silent failures or log-only errors)
- Linear API client errors (non-200 status, malformed JSON, API errors) are captured and returned to caller

## Out of Scope

- User authentication or authorization mechanisms
- Web UI, REST API, or GraphQL endpoint
- Issue creation or modification through the service (read-only from Linear)
- Issue deletion (only create/update via UPSERT)
- Data normalization or transformation beyond direct storage
- Metrics collection, distributed tracing, or observability backends
- Batch API operations or bulk data export
- Encryption at rest or in transit (TLS handled by operator infrastructure)

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Fetch issues from Linear API | Client must POST a GraphQL query to https://api.linear.app/graphql with Authorization Bearer token, including team ID in query, and return list of Issue structs with fields: ID, Title, Description, State, Priority, Assignee, CreatedAt, UpdatedAt. | linear.go implementation and linear_test.go test coverage |
| F2 | must | Upsert issues into PostgreSQL | Storage layer must execute INSERT ... ON CONFLICT (id) DO UPDATE to persist issues. Assignee must be stored as JSONB. Transaction must commit atomically or rollback entirely on any error. | storage.go implementation and storage_test.go test coverage |
| F3 | must | Load configuration from environment | Application must read LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables at startup. Return clear error if any required variable is missing. | config.go implementation and config_test.go test coverage |
| F4 | must | Orchestrate sync cycle | Main service must run initial sync on startup, then repeat every 5 minutes via time.Ticker. Each sync calls fetchIssues() then upsert(). Sync timeout is 2 minutes per cycle. | main.go implementation |
| F5 | must | Handle graceful shutdown | Application must catch os.Interrupt (SIGINT) and syscall.SIGTERM signals and exit cleanly after logging shutdown message. Context cancellation must propagate to in-flight operations. | main.go signal.NotifyContext usage |
| F6 | must | Error handling and logging | All error paths must log at appropriate level ([ERROR], [INFO]). Errors must be wrapped with context via fmt.Errorf(%w). No silent failures. | All files implement error wrapping; main.go logs errors |
| F7 | must | Context propagation for cancellation | All blocking operations (HTTP, database, time.Sleep) must respect context.Context. fetchIssues and upsert accept ctx parameter and use it in all I/O calls. | linear.go fetchIssues(ctx) and storage.go upsert(ctx) implementations |
| F8 | must | Unit test coverage | config_test.go tests missing environment variables and valid config; linear_test.go tests successful fetch and API error handling; storage_test.go tests upsert and empty list handling. | config_test.go, linear_test.go, storage_test.go with 4 PASS, 2 SKIP results |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Performance / Sync cycle timeout | Each individual sync cycle must complete within 2 minutes. If fetchIssues or upsert takes longer, context timeout cancels the operation and logs error. Main sync loop continues to next ticker. | main.go uses context.WithTimeout(ctx, 2*time.Minute) |
| NF2 | must | Database reliability | Connection to PostgreSQL must be verified on startup via db.PingContext(). Connection pool managed by database/sql with sensible defaults. | storage.go newDB() calls db.PingContext() |
| NF3 | must | Code quality and idioms | All code must follow Go 1.21 idioms: named errors (errors.New), error wrapping (%w), early returns on error, no interface{}, table-driven tests with t.Run. | All source files adhere to Go idiom guidelines |
| NF4 | must | Security: API key handling | Linear API key must never appear in logs, error messages, or HTTP query strings. Key passed only via Authorization Bearer header. No key leakage in panics or debug output. | linear.go sets Authorization header; config.go does not log key |
| NF5 | must | Observability: structured logging | All log lines must include [INFO], [ERROR], or equivalent prefix. Messages must include operation and outcome (e.g., 'Synced 42 issues', 'Failed to fetch issues: <error>'). | main.go uses log.Printf with prefixes |
| NF6 | must | Test coverage and isolation | Minimum 80% statement coverage on new packages. Tests must be isolated (no shared state, proper cleanup). Table-driven tests for multiple scenarios. No time.Sleep in tests. | All *_test.go files follow table-driven pattern and httptest isolation |

## Dependencies

- github.com/lib/pq v1.10.9 (PostgreSQL driver for Go)
- github.com/stretchr/testify v1.9.0 (assertion library for tests)
- Go 1.21 or later (language and stdlib)

