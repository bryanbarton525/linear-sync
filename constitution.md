# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Create a Go service that synchronizes Linear.app issues to a PostgreSQL database by writing 9 pre-specified source files verbatim, validating build and test success, and deploying to the target GitHub repository.

## Goals

- Write exactly 9 source files using write_file tool with no modifications
- Validate the service builds successfully with specified flags
- Ensure tests pass with expected results (4 PASS, 2 SKIP)
- Commit and push all code to github.com/bryanbarton525/linear-sync main branch

## Constraints

- Files must be written VERBATIM with no additions, removals, or character changes
- Must NOT run go mod init, go mod tidy, or go work init
- Must use -buildvcs=false flag for build and test commands
- Must write ALL files before running any build or test commands
- Storage tests are expected to skip without PostgreSQL connection

## Audience

Go developer implementing the Linear.app synchronization service

## Output Medium

Go source code files in a Git repository with automated build validation

## Acceptance Criteria

- All 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are written exactly as specified
- Build command 'go build -buildvcs=false .' succeeds without errors
- Test command 'go test -buildvcs=false . -v' shows exactly 4 PASS and 2 SKIP results
- All code is committed and pushed to github.com/bryanbarton525/linear-sync main branch
- The service implements Linear.app API integration with PostgreSQL persistence
- Configuration loads from environment variables (LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL)
- Service runs with graceful shutdown handling and periodic sync intervals

## Out of Scope

- Modifying the provided source code in any way
- Adding additional files beyond the 9 specified
- Running go mod commands or workspace initialization
- Setting up PostgreSQL for storage tests to pass
- Custom build configurations or deployment scripts

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write source files verbatim | Use write_file tool to create exactly 9 source files with no character modifications | task instructions |
| F2 | must | Build service executable | Run 'go build -buildvcs=false .' to compile the Linear sync service | task instructions |
| F3 | must | Execute test suite | Run 'go test -buildvcs=false . -v' and verify 4 PASS + 2 SKIP results | task instructions |
| F4 | must | Git repository deployment | Commit all files and push to github.com/bryanbarton525/linear-sync main branch | workflow request |
| F5 | must | Linear API integration | Service fetches issues from Linear.app using GraphQL API with Bearer authentication | linear.go file content |
| F6 | must | PostgreSQL persistence | Upsert Linear issues to PostgreSQL with conflict resolution on issue ID | storage.go file content |
| F7 | must | Configuration management | Load API key, team ID, and database URL from environment variables | config.go file content |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Build reproducibility | Service must build with -buildvcs=false flag without dependency initialization | task constraints |
| NF2 | must | Test isolation | Storage tests skip gracefully when PostgreSQL is unavailable | expected test results |
| NF3 | must | File integrity | Source files must match provided specifications character-for-character | verbatim requirement |
| NF4 | must | Graceful shutdown | Service handles SIGTERM and SIGINT for clean termination | main.go signal handling |
| NF5 | must | Sync periodicity | Issues are synchronized every 5 minutes with 2-minute timeout per sync | main.go timing configuration |

## Dependencies

- Go 1.21 runtime environment
- github.com/lib/pq v1.10.9 for PostgreSQL connectivity
- github.com/stretchr/testify v1.9.0 for test assertions
- Git repository access for github.com/bryanbarton525/linear-sync

