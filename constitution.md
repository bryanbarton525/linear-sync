# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Build a production-ready Go service that syncs Linear.app issues to PostgreSQL with periodic polling, maintaining a persistent, up-to-date record of team issues with graceful error handling and proper observability.

## Goals

- Implement a working Go service that fetches issues from Linear.app GraphQL API
- Persist issues to PostgreSQL with upsert semantics to handle conflicts
- Implement automatic periodic syncing on a 5-minute interval
- Ensure all code compiles and passes the configured test suite
- Successfully push working code to github.com/bryanbarton525/linear-sync main branch

## Constraints

- Write all 9 source files verbatim as provided—no character additions, removals, or modifications
- Do not run go mod init, go mod tidy, or go work init—go.mod and go.sum are provided
- Go 1.21 or compatible runtime
- Build command: go build -buildvcs=false .
- Test command: go test -buildvcs=false . -v
- Expected test results: 4 PASS and 2 SKIP (storage tests skip without PostgreSQL—this is expected and correct)
- Code must be pushed to the main branch of github.com/bryanbarton525/linear-sync

## Audience

Go developers integrating Linear.app issue tracking with PostgreSQL databases; DevOps teams deploying issue-sync services.

## Output Medium

Go source code repository with passing tests, compiled binary, and commits on the main branch.

## Acceptance Criteria

- All 9 source files are written exactly as provided with no alterations
- go build -buildvcs=false . succeeds with exit code 0 and produces a valid binary
- go test -buildvcs=false . -v shows exactly 4 PASS and 2 SKIP with no FAIL
- Code imports all required packages without errors (context, http, sql, json, time, etc.)
- Storage and config tests verify core functionality; storage tests correctly skip when PostgreSQL is unavailable
- Linear API client tests mock HTTP responses and validate GraphQL parsing and error handling
- Code is committed and pushed to github.com/bryanbarton525/linear-sync on the main branch
- Commits include a Co-authored-by trailer for Copilot as per git conventions

## Out of Scope

- Modifying go.mod or go.sum files
- Adding new dependencies or packages beyond those provided
- Creating database schema migrations or setup scripts
- Deploying the service to any environment
- Modifying any provided source code after initial write

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Configuration Loading | Read and validate three required environment variables: LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL. Return descriptive errors if any are missing. | config.go |
| F2 | must | Database Connection | Open a PostgreSQL connection with the provided connection string, ping to verify connectivity, and provide a reusable *sql.DB handle with proper error handling. | storage.go |
| F3 | must | Linear API Client | Implement a GraphQL client that authenticates with Linear API using Bearer token, constructs a team-scoped issues query, and handles HTTP and JSON decode errors. | linear.go |
| F4 | must | Issue Data Mapping | Parse Linear API JSON responses (nodes, state objects, assignee details) into strongly-typed Issue and Assignee structs with proper timestamp handling. | linear.go |
| F5 | must | Upsert to PostgreSQL | Insert or update issues in linear_issues table using ON CONFLICT DO UPDATE; serialize assignee to JSONB; wrap all operations in a transaction; handle empty issue lists gracefully. | storage.go |
| F6 | must | Periodic Polling | Initialize a 5-minute ticker that triggers fetchIssues and upsert operations repeatedly until the service receives a shutdown signal. | main.go |
| F7 | must | Graceful Shutdown | Register signal handlers for SIGINT and SIGTERM; propagate shutdown via context cancellation; exit cleanly without dangling goroutines. | main.go |
| F8 | should | Error Logging | Log sync successes, failures, API errors, and shutdown events with [INFO], [ERROR] prefixes; wrap errors with context to aid debugging. | main.go |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Code Compilation | All 9 source files must compile without errors or warnings using go build -buildvcs=false . with Go 1.21. | Build validation |
| NF2 | must | Test Coverage | Minimum 4 passing tests covering config validation, API client behavior, and storage operations; 2 storage tests must skip gracefully when PostgreSQL is unavailable. | Test validation |
| NF3 | must | Transaction Safety | All database writes must use transactions with proper rollback on error; no silent failures or partial updates. | storage.go |
| NF4 | must | Context Propagation | All I/O operations (HTTP requests, database queries) must accept and respect context.Context; no blocking operations without timeout or cancellation checks. | Concurrency requirements |
| NF5 | must | JSON Serialization | Assignee structs must correctly marshal/unmarshal to/from JSONB; nil assignees must be handled without panic. | storage.go, linear.go |
| NF6 | should | Idiomatic Go | Code must follow Go idioms: error wrapping with fmt.Errorf %w, early returns on error, no interface{}, named returns only when necessary. | Code generation guidelines |

## Dependencies

- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (Test assertions)

