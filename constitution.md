# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A production-ready Go service that periodically syncs Linear.app team issues to PostgreSQL, providing reliable data persistence with graceful error handling and proper lifecycle management.

## Goals

- Implement a fully functional Linear.app-to-PostgreSQL sync service
- Persist synced issues with conflict resolution via upsert
- Execute periodic background synchronization every 5 minutes
- Handle environment-based configuration (API key, team ID, database URL)
- Provide graceful shutdown on termination signals
- Log operations and errors for observability

## Constraints

- Must use provided go.mod and go.sum without modification
- Must write exactly 9 source files verbatim without any changes
- Must compile with `go build -buildvcs=false .` with no errors
- Must pass test suite with `go test -buildvcs=false . -v` showing 4 PASS and 2 SKIP
- Storage tests skip without PostgreSQL available—this is expected and correct
- No `go mod init`, `go mod tidy`, or `go work init` allowed
- Must push to github.com/bryanbarton525/linear-sync main branch

## Audience

Development and DevOps teams deploying the Linear.app sync service in production environments

## Output Medium

Compiled Go binary and source code committed to GitHub repository

## Acceptance Criteria

- All 9 source files (main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go, go.mod, go.sum) written exactly as specified
- `go build -buildvcs=false .` completes successfully with no compilation errors
- `go test -buildvcs=false . -v` shows exactly 4 PASS and 2 SKIP results
- Service binary is executable and ready for deployment
- Code is committed and pushed to github.com/bryanbarton525/linear-sync on main branch
- All required dependencies from go.mod are available without additional `go mod` operations

## Out of Scope

- Database schema creation or migration scripts
- Docker containerization
- Kubernetes manifests or deployment configuration
- CI/CD pipeline setup
- Authentication or authorization beyond API key validation
- Advanced filtering or querying of Linear issues

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Configuration Loading | Load required configuration from environment variables: LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL. Return error if any required variable is missing. | config.go implementation |
| F2 | must | PostgreSQL Connection | Establish and verify connection to PostgreSQL database. Support connection string format from DATABASE_URL environment variable. | storage.go newDB function |
| F3 | must | Linear API Integration | Fetch issues from Linear.app using GraphQL API. Accept team ID as parameter and return parsed Issue structs with ID, title, description, state, priority, assignee, and timestamps. | linear.go fetchIssues method |
| F4 | must | Issue Parsing | Decode Linear API GraphQL response and transform raw JSON into Issue domain objects. Handle nested state and assignee objects correctly. | linear.go rawResponse and Issue types |
| F5 | must | Database Upsert | Insert or update issues in linear_issues table with conflict resolution (ON CONFLICT DO UPDATE). Store assignee as JSONB. Use transactions for atomicity. | storage.go upsert method |
| F6 | must | Periodic Sync Loop | Run synchronization every 5 minutes after initial sync. Use time.Ticker for periodic execution. Execute one full sync cycle per tick. | main.go ticker loop |
| F7 | must | Graceful Shutdown | Listen for SIGTERM and SIGINT signals. Cancel in-flight operations and exit cleanly on signal receipt. | main.go signal handling |
| F8 | should | Operational Logging | Log sync start, completion with issue count, and any errors encountered. Use log package with [INFO] and [ERROR] prefixes. | main.go logging |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Build Success | Source code must compile without errors using `go build -buildvcs=false .` on Go 1.21+ | Toolchain validation |
| NF2 | must | Test Coverage | Test suite must execute with `go test -buildvcs=false . -v` and produce exactly 4 PASS results and 2 SKIP results (storage tests skip when PostgreSQL unavailable) | Test suite (config_test.go, linear_test.go, storage_test.go) |
| NF3 | must | Error Wrapping | All errors must be wrapped with context using fmt.Errorf with %w verb. No silent error swallowing allowed. | Code generation idioms |
| NF4 | must | Context Propagation | All blocking operations must accept and respect context.Context. HTTP requests must use context-aware functions. Sync operations must have 2-minute timeout. | main.go and storage.go context handling |
| NF5 | must | Transaction Safety | Database upsert operations must execute within transactions. Rollback on error. Use prepared statements to prevent injection. | storage.go upsert implementation |
| NF6 | must | Dependency Stability | Use only dependencies declared in go.mod (github.com/lib/pq, github.com/stretchr/testify). No external HTTP client libraries beyond stdlib. | go.mod and go.sum |
| NF7 | must | Repository Push | All code must be committed with proper commit message and pushed to github.com/bryanbarton525/linear-sync on main branch | Git workflow |

## Dependencies

- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (testing assertions)
- Go 1.21+ runtime
- PostgreSQL 12+ (optional for storage tests; tests skip gracefully if unavailable)
- Linear.app account with API key and team ID

---

## Constitution Amendment — Cycle 1

Root cause: Implementation defect in task bc4704d2. The Pod failed to write any of the 9 required source files to the workspace despite having write_file tool available. The implementation-summary artifact incorrectly reported success, but workspace_preflight validation confirms zero Go source files exist in /var/lib/go-orca/workspaces/77e711cb-2688-4f21-9a5e-52b6f9bccb4e. This cascaded into failed build validation, failed test validation, and a premature git commit with no source code. No requirement gap, design gap, or environment failure—purely an execution failure. Remediation: Pod must retry task bc4704d2 using write_file tool to create all 9 files verbatim from task specification. No Constitution amendment required—all requirements are clear and complete. Build and test validation will auto-proceed once source files are written. Retry commit task only after validation passes.

---

## Constitution Amendment — Cycle 2

All five QA blocking issues stem from a single implementation defect: task b60a6aed failed to write any of the 9 required source files to the workspace directory. The write_file tool calls were not executed successfully, leaving the workspace empty and cascading into failed build, test, and git validation. Classification: all issues are implementation defects or validation failures dependent on the source file write failure. No requirement gap—the Constitution is complete and unambiguous. No design gap—the architecture is frozen. Remediation: Task b60a6aed must be retried using explicit write_file tool calls with pre- and post-write verification checks. Once files are written successfully, automated build and test validation will proceed automatically. After validation passes with expected results (4 PASS + 2 SKIP), task bd1c1b31 will execute git commit and push to main branch. The remediation task graph already includes proper prerequisite gating to prevent premature git operations.
