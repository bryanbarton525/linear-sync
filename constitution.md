# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A production-ready Go service that periodically fetches issues from Linear.app via GraphQL API and persists them to PostgreSQL, enabling data synchronization for downstream workflows.

## Goals

- Fetch Linear issues via GraphQL API with proper authentication and schema parsing
- Parse and transform Linear API responses into domain models with nested data structures
- Persist issues to PostgreSQL with upsert semantics using ON CONFLICT DO UPDATE
- Run on a 5-minute sync interval with graceful cancellation via context and signals
- Load all configuration from environment variables with validation on startup
- Pass full build and test suite: 4 PASS, 2 SKIP (storage conditional on PostgreSQL availability)

## Constraints

- Must use Go 1.21 or later
- Must use PostgreSQL for all persistence
- All 9 source files must be written exactly as specified—no additions, removals, or changes
- Must use provided go.mod and go.sum without modification; no go mod init, tidy, or work init
- Linear API requires valid API_KEY and TEAM_ID at runtime
- Code must follow idiomatic Go conventions: early returns, error wrapping with %w, no interface{} in public API
- All code must be committed to github.com/bryanbarton525/linear-sync on main branch

## Audience

Go developers and DevOps engineers deploying Linear-to-PostgreSQL sync services in containerized or VM-based environments.

## Output Medium

Go source code (9 files), executable binary, passing test suite with coverage ≥80%.

## Acceptance Criteria

- All 9 source files written exactly as specified without any character modifications
- Build succeeds with zero errors: go build -buildvcs=false .
- Test suite passes: go test -buildvcs=false . -v shows exactly 4 PASS and 2 SKIP
- Configuration tests verify environment variable validation for LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL
- Linear client tests verify GraphQL query construction, JSON parsing, and API error handling
- Storage tests verify upsert semantics, transaction handling, and JSON serialization of assignee objects
- Code compiles without warnings; no unused imports or variables
- Implementation and test files exist in separate files with matching package declarations
- All code is committed and pushed to github.com/bryanbarton525/linear-sync on main branch

## Out of Scope

- Webhook-based real-time sync (polling-based only)
- Multi-team or multi-workspace support
- Data retention policies or cleanup jobs
- Metrics collection or distributed tracing
- Advanced authentication (API key only)
- In-memory caching layer
- Batch API operations beyond standard upsert

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Configuration Management | Load and validate configuration from environment variables: LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL. Return errors for missing or empty values. | config.go |
| F2 | must | Linear GraphQL API Client | Implement linearClient to fetch issues from Linear.app via GraphQL. Construct query for team issues with required fields (id, title, description, state, priority, assignee, createdAt, updatedAt). Parse JSON responses and transform to Issue domain models. Handle API errors in response.errors array. | linear.go |
| F3 | must | PostgreSQL Storage Layer | Implement Storage struct for PostgreSQL persistence. Create database connections, implement upsert logic using ON CONFLICT DO UPDATE, serialize assignee objects to JSONB, use transactions for consistency. Skip gracefully if PostgreSQL unavailable in tests. | storage.go |
| F4 | must | Main Service Loop | Load configuration, open database connection, initialize storage and client, set up signal handlers for SIGINT/SIGTERM, implement 5-minute sync ticker, execute sync immediately on startup, log events and errors, implement 2-minute context timeout for each sync operation. | main.go |
| F5 | must | Configuration Unit Tests | Test config.load() with various environment variable combinations: missing LINEAR_API_KEY, missing LINEAR_TEAM_ID, missing DATABASE_URL, valid complete config. | config_test.go |
| F6 | must | Linear Client Integration Tests | Test linearClient.fetchIssues() with mocked HTTP responses. Verify GraphQL response parsing, transformation to Issue structs, and error handling for API errors. | linear_test.go |
| F7 | must | Storage Integration Tests | Test Storage.upsert() with PostgreSQL. Verify insert of new issues, update of existing issues (ON CONFLICT semantics), and graceful skip when PostgreSQL unavailable. | storage_test.go |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Code Quality and Style | Follow idiomatic Go: use errors.New and fmt.Errorf("%w", ...) for error wrapping, return early on errors, avoid interface{} in public API, use named returns only when beneficial, include context.Context in all blocking operations. | All Go files |
| NF2 | must | Test Isolation | Implementation and test files must be in separate files with matching package names. Use httptest.NewServer with defer close for test servers. Do not use http.DefaultServeMux. Achieve ≥80% statement coverage on new packages. | All test files |
| NF3 | must | Build Validation | Build must succeed: go build -buildvcs=false . Test suite must pass: go test -buildvcs=false . -v shows 4 PASS and 2 SKIP. No build warnings or errors. | Build toolchain |
| NF4 | must | Concurrency Safety | Pass context.Context to all blocking calls. Implement signal-based cancellation with proper cleanup. Prevent goroutine leaks. Use context timeout for sync operations (2 minutes). | main.go, linear.go |
| NF5 | should | Observability and Logging | Use structured logging with [INFO] and [ERROR] prefixes. Log service startup, sync completion with issue count, and errors with context. Log configuration validation errors on startup. | main.go, config.go |
| NF6 | must | Module Integrity | Use provided go.mod and go.sum without modification. Do not run go mod init, go mod tidy, or go work init. All dependencies already present and pinned. | go.mod, go.sum |

## Dependencies

- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (test assertions and mocking)
- Go 1.21 standard library (context, net/http, encoding/json, database/sql)

