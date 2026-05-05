# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Build a production-ready Go service that synchronizes Linear.app issues to a PostgreSQL database, enabling offline analytics, backup, and data export capabilities.

## Goals

- Fetch issues from Linear API on a configurable schedule
- Persist issues to PostgreSQL with upsert semantics
- Handle API errors gracefully and log them appropriately
- Support graceful shutdown via OS signals
- Maintain code quality through idiomatic Go patterns

## Constraints

- Use exactly 9 source files as specified in the task
- Do not modify go.mod or go.sum — they are provided
- Run go build and go test after writing all files
- Expected test results: 4 PASS and 2 SKIP (storage tests skip without PostgreSQL)
- Push final code to github.com/bryanbarton525/linear-sync on main branch
- Follow Go idioms from go-idioms.md references
- Never create new artifact versions during remediation cycles

## Audience

Go developers, DevOps engineers, backend platform teams responsible for Linear data integration

## Output Medium

Go source code files committed to GitHub repository workflow branch

## Acceptance Criteria

- go build -buildvcs=false . exits with code 0
- go test -buildvcs=false . -v shows exactly 4 PASS and 2 SKIP
- All storage_test.go tests skip gracefully when DATABASE_URL is not set
- No compilation errors in any file
- No test panics or leaked goroutines in concurrent tests
- All test servers properly closed with defer ts.Close()
- All HTTP tests use httptest.NewServer with fresh ServeMux

## Out of Scope

- User authentication flow beyond API key
- Issue filtering or querying logic
- Web dashboard or UI components
- Email notifications
- Multi-team support beyond single team ID

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Load configuration from environment variables | Read LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from OS environment. Return error if any required variable is missing. | config.go |
| F2 | must | Fetch issues from Linear API | Implement GraphQL query to Linear API, handle response decoding, and return array of Issue structs. Handle non-200 status codes and API errors appropriately. | linear.go |
| F3 | must | Persist issues to PostgreSQL | Implement upsert logic using ON CONFLICT clause. Serialize assignee data to JSONB. Handle empty issue lists gracefully. | storage.go |
| F4 | must | Run periodic sync loop | Execute initial sync, then repeat every 5 minutes via ticker. Exit gracefully on SIGINT or SIGTERM. | main.go |
| F5 | must | Create database connection with ping verification | Open PostgreSQL connection and verify connectivity via PingContext before accepting it as valid. | storage.go |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Context propagation | All blocking operations must accept context.Context as first parameter and check ctx.Done() in select statements. | code-generation/agent overlay |
| NF2 | must | Error wrapping | Wrap all errors using fmt.Errorf("%w", err) to preserve original context in the chain. | go-idioms.md |
| NF3 | must | Test isolation | All tests must be independent with no shared state. Setup and teardown must be explicit. Test servers must be closed. | qa-checklist.md |
| NF4 | must | Test file structure | Implementation and tests must be in separate files with matching package declaration. Never mix in a single file. | code-generation |
| NF5 | must | No DefaultServeMux pollution | Tests must create fresh http.ServeMux instances, never use http.DefaultServeMux. | qa-checklist.md |
| NF6 | must | Consistent timestamp usage | Use UnixMilli() consistently across all time operations; never mix with UnixNano(). | qa-checklist.md |

## Dependencies

- github.com/lib/pq v1.10.9
- github.com/stretchr/testify v1.9.0

