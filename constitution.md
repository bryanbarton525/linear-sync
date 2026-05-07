# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Refactor the linear-sync Go service from a flat root layout to follow golang-standards/project-layout conventions, establishing proper package boundaries and improving code organization for long-term maintainability.

## Goals

- Reorganize Go source files into cmd/, internal/config, internal/linear, and internal/storage directories following idiomatic Go structure
- Establish proper package boundaries and clear dependency flow between packages
- Maintain existing functionality, test coverage, and API behavior throughout the refactor
- Enable future contributors to navigate and extend the codebase confidently

## Constraints

- go.mod module path remains 'linear-sync'; Go version remains 1.21
- No new external dependencies may be introduced
- All existing tests must pass or skip appropriately; no tests may be removed or disabled
- Build output must be zero errors; compilation must succeed for all packages
- No changes to .gitignore, LICENSE, or top-level configuration files
- Import paths must follow 'linear-sync/internal/<package>' convention

## Audience

Go developers and project maintainers requiring standard project structure; contributors familiar with golang-standards/project-layout conventions

## Output Medium

Git repository (code files in cmd/ and internal/ directories) with git commit history tracking file moves and deletions

## Acceptance Criteria

- Directory structure exactly matches target: cmd/linear-sync/, internal/config/, internal/linear/, internal/storage/ with specified .go files
- go build -buildvcs=false ./... executes successfully with zero errors
- go test -buildvcs=false ./... produces: PASS for internal/config and internal/linear; SKIP for internal/storage (PostgreSQL unavailable in CI)
- All old root-level .go files (main.go, config.go, linear.go, storage.go, and all *_test.go) are deleted from repository root
- Artifact cruft files (git-operations-complete, linear-sync-implementation) are deleted
- No circular import dependencies exist between packages
- Package internal/linear.Issue, internal/linear.Client are properly exported and used by internal/storage and cmd/linear-sync
- main.go correctly imports and uses all internal packages; time.Ticker goroutine runs FetchIssues every 5 minutes with proper context cancellation
- Tests use httptest for HTTP mocking; database tests skip gracefully when DATABASE_URL unset
- All imports resolve correctly; no 'undefined' or 'no such package' errors

## Out of Scope

- Adding new features, endpoints, or API behavior
- Modifying the GraphQL schema or Linear API integration beyond what exists
- Changing database schema or migrations
- Adding new dependencies or updating existing versions in go.mod/go.sum
- Performance optimization or refactoring of algorithm logic
- Documentation beyond inline code comments required by idiom

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Create cmd/linear-sync/main.go | Create main package entrypoint that loads config, initializes Linear API client and database connection, runs a time.Ticker goroutine to fetch issues every 5 minutes, and handles OS signals for graceful shutdown via context cancellation. | Workflow specification; golang-standards/project-layout |
| F2 | must | Create internal/config package | Implement config.go with Config struct (APIKey, TeamID, DatabaseURL fields) and Load() function that reads LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL environment variables; return error if any variable is missing. | Workflow specification |
| F3 | must | Create internal/config tests | Implement config_test.go with table-driven tests covering: missing env variables return error for each, all variables set returns correct Config struct. | Workflow specification |
| F4 | must | Create internal/linear package | Implement linear.go with Assignee struct (Name field), Issue struct (ID, Title, State, Assignee fields), Client struct (apiKey, http.Client fields), NewClient() constructor, and FetchIssues(ctx, teamID) method that performs GraphQL POST to https://api.linear.app/graphql and returns []Issue. | Workflow specification |
| F5 | must | Create internal/linear tests | Implement client_test.go with httptest-based tests: TestFetchIssues (happy path with valid GraphQL response), TestFetchIssuesHTTPError (HTTP error handling), TestFetchIssuesGraphQLError (GraphQL error response handling). | Workflow specification |
| F6 | must | Create internal/storage package | Implement storage.go with Storage struct (db *sql.DB field), NewDB(connStr) function, New(db) constructor, and Upsert(ctx, issues) method that upserts []linear.Issue into linear_issues table using parameterized queries. | Workflow specification |
| F7 | must | Create internal/storage tests | Implement storage_test.go with tests TestStorageUpsert and TestStorageUpsertEmpty; both skip if DATABASE_URL environment variable is not set. | Workflow specification |
| F8 | must | Delete old root-level files | Remove all .go files from repository root: main.go, config.go, linear.go, storage.go, main_test.go, config_test.go, linear_test.go, storage_test.go and any other implementation files that have been moved to cmd/ or internal/. | Workflow specification |
| F9 | must | Delete artifact cruft | Remove artifact files: git-operations-complete and linear-sync-implementation from repository root. | Workflow specification |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Build validation | Execution of 'go build -buildvcs=false ./...' must complete successfully with zero compilation errors, warnings, or undefined symbols. | Workflow acceptance criteria |
| NF2 | must | Test execution validation | Execution of 'go test -buildvcs=false ./...' must produce: PASS status for internal/config and internal/linear packages; SKIP status for internal/storage package (due to missing PostgreSQL in CI environment). | Workflow acceptance criteria |
| NF3 | must | Package structure adherence | Directory layout and file organization must precisely match golang-standards/project-layout: cmd/<app>/main.go as entrypoint, internal/<package>/ for unexported code, proper separation of concerns by package. | golang-standards/project-layout convention |
| NF4 | must | Import path correctness | All import statements must use the canonical module path 'linear-sync/internal/<package>'; no relative imports or package aliasing outside idiomatic patterns. | Go module system and idiomatic practice |
| NF5 | must | Cyclic dependency prevention | No circular import dependencies between packages: cmd → internal/config, internal/linear, internal/storage; internal/storage → internal/linear; internal/linear and internal/config must not import each other. | Go idiom and package design |
| NF6 | must | Test isolation and cleanup | All tests must use isolated test fixtures (httptest.NewServer with defer close, fresh http.ServeMux); database tests must skip gracefully when DATABASE_URL is unset; no test pollution or shared state across test runs. | Go testing best practices |
| NF7 | must | Context propagation | All I/O functions (FetchIssues, Upsert) must accept context.Context as first parameter and respect cancellation signals; no background goroutines without explicit lifecycle management. | Go idiom and concurrency safety |
| NF8 | must | Error handling completeness | All error paths must be handled explicitly; no silent error suppression; errors must include sufficient context for diagnosis without requiring stack traces. | Go idiom and observability |

## Dependencies

- go.mod must remain unchanged (module: linear-sync, version: 1.21)
- go.sum must remain unchanged
- Standard library packages: context, database/sql, encoding/json, fmt, os, os/signal, sync, syscall, time, net/http, bytes, io, testing
- No external dependencies introduced or modified

