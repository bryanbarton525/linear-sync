# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Provide a small, reliable Go service (linear-sync) that continuously syncs Linear.app issues for a team into PostgreSQL and exposes simple HTTP endpoints for inspection and manual sync control.

## Goals

- Poll the Linear GraphQL API every 5 minutes and upsert issues into Postgres.
- Expose HTTP endpoints: GET /issues, POST /sync, GET /healthz.
- Be configuration-driven via environment variables (LINEAR_API_KEY, DATABASE_URL) and safe to run in containers.
- Ship with a clear, testable acceptance criteria set so the toolchain can validate correctness.

## Constraints

- Polling interval default: 5 minutes (300s) and must be configurable via POLL_INTERVAL (seconds).
- Use LINEAR_API_KEY env var for Linear API auth and DATABASE_URL for Postgres connection.
- Database table schema is fixed: issues (id PRIMARY KEY, title, status, priority, updated_at).
- Module path: github.com/bryanbarton525/linear-sync.
- All long-running/blocking operations must accept context.Context for cancellation.

## Audience

Operators and backend engineers who need a lightweight service to mirror Linear issues into Postgres for reporting, analytics, or downstream processing.

## Output Medium

Go module (command) with source code, tests, and instructions; the service exposes an HTTP API and persists data to PostgreSQL.

## Acceptance Criteria

- Build: go build ./... succeeds with module path github.com/bryanbarton525/linear-sync.
- Tests: go test ./... passes.
- Database schema: a Postgres table named issues exists with columns: id (text PRIMARY KEY), title (text), status (text), priority (int), updated_at (timestamptz).
- Upsert behavior: When sync runs (manual via POST /sync or automatic via poll), issues returned by Linear are upserted into issues table by id; updated_at reflects the issue's updatedAt timestamp.
- HTTP endpoints: GET /issues returns application/json array of rows from issues table; POST /sync triggers an immediate sync and returns 200 on success; GET /healthz returns 200 when the app can connect to the DB and (optionally) the Linear API.
- Polling: When LINEAR_API_KEY and DATABASE_URL are set, the service polls the Linear GraphQL API at the configured interval (default 300s) and performs upserts; this behavior is observable in logs or by verifying DB state after two polling intervals.
- Configuration: LINEAR_API_KEY and DATABASE_URL are required in runtime; missing required vars cause startup failure with clear error logs.
- Graceful shutdown: Service stops polling and returns non-200 healthz during shutdown; running go vet/format passes (gofmt/gofmt).

## Out of Scope

- Syncing Linear attachments, comments, or user mapping beyond storing basic issue fields.
- Providing a web UI or frontend beyond the three HTTP endpoints specified.
- Performing advanced conflict resolution beyond upsert-by-id using updated_at.
- Supporting multiple teams in one service instance (single-team per instance is assumed).

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Linear API Poller | Poll the Linear GraphQL API on a configurable interval (default 300 seconds) using the LINEAR_API_KEY env var, fetch issues for the configured team, and produce a deterministic list of issues with fields: id, title, status, priority, updatedAt. | Request and director |
| F2 | must | Postgres Upsert | Connect to Postgres using DATABASE_URL and upsert issues into the issues table. The table columns must be id (text PRIMARY KEY), title (text), status (text), priority (int), updated_at (timestamptz). Upsert must replace or update rows when Linear's updatedAt is newer. | Request |
| F3 | must | HTTP API | Expose GET /issues (returns JSON array of issues from DB), POST /sync (triggers immediate sync and returns operation status), and GET /healthz (returns service health based on DB connectivity and optionally last successful sync). | Request |
| F4 | must | Environment Configuration | Read runtime configuration from environment variables: LINEAR_API_KEY (required), DATABASE_URL (required), POLL_INTERVAL (optional, seconds), and LOG_LEVEL (optional). Fail fast on missing required variables. | Request |
| F5 | should | Manual Sync Trigger | POST /sync must trigger a foreground sync using the same logic as the poller and return 200 on success; concurrent manual triggers are serialized or deduplicated to prevent concurrent syncs. | Request |
| F6 | should | Error Handling and Retries | Transient errors when calling Linear or Postgres should be retried with exponential backoff for the duration of an operation; persistent failures must be logged and exposed via health/metrics. | Common operational practice |
| F7 | could | Logging and Observability | Emit structured logs for sync start/finish, errors, and number of records upserted. Optionally expose basic metrics (sync duration, last sync time). | Operational needs |
| F8 | must | Context Propagation & Graceful Shutdown | All blocking operations must accept context.Context to allow cancellation. The service should stop accepting requests, cancel in-flight syncs, and close DB connections on SIGINT/SIGTERM. | Architect overlay and best practices |
| F9 | must | Data Model Mapping | Map Linear issue fields to DB columns: Linear.id -> id, Linear.title -> title, Linear.state.name or comparable -> status, Linear.priority (or mapped numeric) -> priority, Linear.updatedAt -> updated_at (timestamptz). | Request |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Reliability | Service must be resilient to transient network failures; polls should continue after transient errors and not crash the process. | Operational requirements |
| NF2 | must | Testability | Unit tests and integration tests must exist covering: the GraphQL request builder, DB upsert logic (using a test Postgres or test double), and HTTP handlers. go test ./... must pass. | Toolchain validation |
| NF3 | should | Performance | Single-instance memory and CPU usage must be small (suitable for a single container); bulk upserts must be batched where appropriate to avoid per-row transactions. | Operational constraints |
| NF4 | must | Security | Do not log sensitive values (LINEAR_API_KEY or DATABASE_URL). Use parameterized queries to avoid SQL injection. Secrets must come from env vars only; do not embed secrets in code. | Security best practices |
| NF5 | should | Observability | Provide enough logs and status via /healthz to diagnose failures; include last sync time and status in logs/health output. | Operational needs |

## Dependencies

- Linear GraphQL API (requires LINEAR_API_KEY)
- PostgreSQL instance (accessible via DATABASE_URL)
- Go toolchain (go 1.20+ recommended)
- pgx or database/sql Postgres driver (implementation detail)
- Network connectivity to Linear and Postgres

