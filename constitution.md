# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Deliver a production-ready Go service that syncs Linear.app issues to PostgreSQL, demonstrating idiomatic Go patterns, proper error handling, context propagation, and comprehensive test coverage.

## Goals

- Write all 9 source files exactly as specified (main.go, config.go, linear.go, storage.go, go.mod, go.sum, config_test.go, linear_test.go, storage_test.go)
- Validate build succeeds with `go build -buildvcs=false .`
- Validate test suite executes with `go test -buildvcs=false . -v` showing 4 PASS and 2 SKIP
- Commit and push all code to github.com/bryanbarton525/linear-sync on the main branch

## Constraints

- Write each file VERBATIM—no additions, removals, or character changes
- Do not run go mod init, go mod tidy, or go work init; use provided go.mod and go.sum
- Build must use -buildvcs=false flag
- Tests must use -buildvcs=false flag and -v for verbose output
- Storage tests (2) will skip without PostgreSQL—this is expected and correct
- Config and linear tests (4) must pass; only database connectivity tests skip

## Audience

Go developers, infrastructure engineers, and DevOps teams deploying Linear-to-PostgreSQL sync services.

## Output Medium

GitHub repository containing Go source code, compiled binary, and test results.

## Acceptance Criteria

- All 9 source files written exactly as provided with no modifications
- `go build -buildvcs=false .` exits with status 0 and produces no errors
- `go test -buildvcs=false . -v` output shows exactly 4 PASS and 2 SKIP
- Code follows idiomatic Go: error wrapping with %w, context propagation, no silent error swallowing
- Test isolation verified: config_test.go and linear_test.go pass without external service dependencies; storage_test.go skips gracefully when PostgreSQL unavailable
- All commits pushed to github.com/bryanbarton525/linear-sync on the main branch

## Out of Scope

- Database schema migrations or initialization beyond test table creation
- Docker or container configuration
- CI/CD pipeline setup beyond git push
- Documentation beyond code comments
- Performance tuning or load testing
- Dependency upgrades or modifications to go.mod/go.sum

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write all 9 source files verbatim | Write main.go, config.go, linear.go, storage.go, go.mod, go.sum, config_test.go, linear_test.go, and storage_test.go exactly as provided with zero character changes. | Task specification: Write exactly the following 9 source files into the workspace using the write_file tool. Write each file VERBATIM. |
| F2 | must | Build succeeds without errors | Execute `go build -buildvcs=false .` and verify exit code 0 with no compilation errors. | Task specification: After writing all files, run: go build -buildvcs=false . |
| F3 | must | Test suite validates with expected pass/skip counts | Execute `go test -buildvcs=false . -v` and verify 4 PASS (config tests and linear client tests) and 2 SKIP (storage tests skip without PostgreSQL). | Task specification: Tests should show 4 PASS and 2 SKIP (storage tests skip without PostgreSQL — this is expected and correct). |
| F4 | must | Commit and push to GitHub | Commit all 9 files and push to github.com/bryanbarton525/linear-sync on the main branch. | Workflow request: pushing all code to github.com/bryanbarton525/linear-sync on the main branch. |
| F5 | must | Configuration loading from environment variables | Service loads LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment; fails gracefully if any is missing. | config.go and config_test.go implementation |
| F6 | must | Linear API issue fetching | linearClient.fetchIssues() makes authenticated GraphQL request to Linear API, parses response, and returns Issue structs with ID, title, description, state, priority, assignee, and timestamps. | linear.go implementation |
| F7 | must | PostgreSQL persistence | Storage.upsert() accepts context and []Issue, performs INSERT...ON CONFLICT DO UPDATE to persist issues with assignee serialized as JSONB. | storage.go implementation |
| F8 | must | Graceful shutdown and context propagation | main() creates context with signal notification (SIGINT, SIGTERM); passes context to all blocking operations (API calls, database operations); cancels context on shutdown. | main.go implementation |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Idiomatic Go error handling | All errors wrapped with fmt.Errorf("%w", err) or errors.New(); no silent error swallowing; early return pattern for error paths. | code-generation skill: Go idioms |
| NF2 | must | Context propagation | Every blocking function (API calls, database operations, timeouts) receives context.Context parameter; respects context cancellation. | code-generation skill: Concurrency patterns |
| NF3 | must | Test isolation and no external dependencies | config_test.go and linear_test.go pass without PostgreSQL or external services; storage_test.go gracefully skips when database unavailable; tests use httptest.NewServer for mocking. | qa-validation skill: Test isolation |
| NF4 | should | Table-driven test structure | Tests with multiple cases (config_test.go) use table-driven pattern with t.Run; HTTP mocking tests use httptest.NewServer with proper defer cleanup. | code-generation skill: Testing best practices |
| NF5 | must | No dependency modifications | go.mod and go.sum remain exactly as provided; no go mod init, go mod tidy, or dependency upgrades. | Task specification: Do NOT run go mod init, go mod tidy, or go work init |
| NF6 | should | Package organization | All source files in root package (main); no cyclic imports; clear separation between API client (linear.go), storage layer (storage.go), and configuration (config.go). | code-generation skill: Package layout |

## Dependencies

- github.com/lib/pq (PostgreSQL driver)
- github.com/stretchr/testify (assertion library)
- encoding/json (JSON parsing)
- net/http (HTTP client for Linear API)
- database/sql (SQL interface)
- context (cancellation and timeouts)

