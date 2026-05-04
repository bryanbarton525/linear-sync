# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A production-ready Go service that continuously syncs Linear.app issues to PostgreSQL with robust error handling, graceful shutdown, and reliable persistence

## Goals

- Implement a Go service that fetches Linear.app issues via GraphQL API
- Persist issues to PostgreSQL using upsert logic to handle updates
- Run sync loop every 5 minutes with graceful shutdown support
- Provide comprehensive test coverage with unit tests for all major components
- Deploy code to github.com/bryanbarton525/linear-sync main branch

## Constraints

- Must use Go 1.21 as specified in go.mod
- Must use exact file structure and content provided in task specification
- Must not modify or regenerate go.mod and go.sum files
- Must use github.com/lib/pq for PostgreSQL and github.com/stretchr/testify for testing
- Must build with -buildvcs=false flag
- Must not run go mod init, go mod tidy, or go work init

## Audience

Backend engineers responsible for maintaining Linear.app integration and data synchronization infrastructure

## Output Medium

Go source code repository hosted on GitHub

## Acceptance Criteria

- All 9 source files written verbatim with no modifications: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go
- Build succeeds with command: go build -buildvcs=false .
- Tests pass with command: go test -buildvcs=false . -v
- Test results show exactly 4 PASS and 2 SKIP (storage tests skip without PostgreSQL, which is expected and correct)
- Code committed and pushed to github.com/bryanbarton525/linear-sync main branch
- Repository contains all required files in workspace root directory

## Out of Scope

- Database schema migrations or setup automation
- API rate limiting beyond context timeouts
- Monitoring, alerting, or observability tooling
- Deployment manifests or CI/CD pipelines
- Documentation beyond code comments

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write 9 source files verbatim | Write all 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) exactly as specified with no character changes | Task specification |
| F2 | must | Implement Linear.app GraphQL API client | Fetch issues from Linear.app using GraphQL API with authentication, including issue metadata (title, description, state, priority, assignee, timestamps) | linear.go specification |
| F3 | must | Implement PostgreSQL storage layer | Store and update issues in PostgreSQL using upsert logic to handle both inserts and updates, with proper transaction handling and JSONB for assignee field | storage.go specification |
| F4 | must | Implement sync loop with 5-minute intervals | Run continuous sync loop that fetches issues every 5 minutes, with proper logging and error handling for failures | main.go specification |
| F5 | must | Handle graceful shutdown | Listen for SIGTERM and SIGINT signals, cancel context, and shut down cleanly without data loss | main.go specification |
| F6 | must | Load configuration from environment | Read LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables with validation | config.go specification |
| F7 | must | Build successfully | Execute go build -buildvcs=false . without errors after writing all files | Task instructions |
| F8 | must | Pass all tests | Execute go test -buildvcs=false . -v and verify 4 PASS and 2 SKIP (storage tests skip without PostgreSQL) | Task instructions |
| F9 | must | Push to GitHub main branch | Commit all files and push to github.com/bryanbarton525/linear-sync main branch | Workflow request |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Build must succeed | go build -buildvcs=false . must exit with status 0 and produce a valid binary | Go toolchain validation |
| NF2 | must | Tests must pass with expected skip count | go test -buildvcs=false . -v must show 4 PASS and 2 SKIP; storage tests skipping without PostgreSQL is expected and correct behavior | Task acceptance criteria |
| NF3 | must | Follow Go idioms | Code must use idiomatic Go patterns: proper error wrapping, context propagation, defer for cleanup, table-driven tests | code-generation skill |
| NF4 | must | Preserve dependency files | go.mod and go.sum must remain exactly as provided; no go mod tidy or go mod init commands allowed | Task constraints |
| NF5 | must | Use specified dependencies | Must use github.com/lib/pq v1.10.9 for PostgreSQL and github.com/stretchr/testify v1.9.0 for testing as declared in go.mod | go.mod specification |

## Dependencies

- github.com/lib/pq v1.10.9
- github.com/stretchr/testify v1.9.0
- PostgreSQL database (for production use, not required for test pass criteria)

