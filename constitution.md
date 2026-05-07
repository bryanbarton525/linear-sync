# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Create a production-ready Go service that synchronizes Linear.app issues to a PostgreSQL database with continuous sync capabilities, deployed to GitHub at github.com/bryanbarton525/linear-sync.

## Goals

- Write exactly 9 Go source files verbatim to the workspace without any modifications
- Build the service successfully using go build -buildvcs=false
- Verify test suite execution with 4 passing tests and 2 expected skips (storage tests skip without PostgreSQL)
- Commit all source files and push to github.com/bryanbarton525/linear-sync main branch

## Constraints

- Files must be written VERBATIM - no character-level changes, additions, or removals permitted
- Do NOT run go mod init, go mod tidy, or go work init - go.mod and go.sum are already provided
- Must use go build -buildvcs=false for building to avoid VCS requirement errors
- Must use go test -buildvcs=false . -v for testing
- Storage tests are expected to skip in environments without PostgreSQL - this is correct behavior

## Audience

Development and operations teams managing Linear issue synchronization workflows and requiring reliable database persistence of issue tracking data.

## Output Medium

Git repository at github.com/bryanbarton525/linear-sync containing Go source code, configuration files, and test suites.

## Acceptance Criteria

- All 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are written to the workspace exactly as specified
- go build -buildvcs=false . completes successfully with exit code 0
- go test -buildvcs=false . -v executes and reports exactly 4 PASS results and 2 SKIP results
- All source files are committed to the workflow branch
- Code is pushed to github.com/bryanbarton525/linear-sync main branch
- Repository contains a working Linear-to-PostgreSQL sync service with proper error handling, context propagation, and transaction management

## Out of Scope

- Modifying or improving the provided source code
- Adding linters, formatters, or additional tooling beyond build and test
- Creating database migration scripts or schema definitions
- Implementing additional features or refactoring existing code
- Setting up CI/CD pipelines or deployment automation

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write go.mod file | Write the go.mod file with module declaration and exact dependency specifications including github.com/lib/pq v1.10.9 and github.com/stretchr/testify v1.9.0 | Task specification - File: go.mod |
| F2 | must | Write go.sum file | Write the go.sum file with cryptographic checksums for all direct and transitive dependencies | Task specification - File: go.sum |
| F3 | must | Write main.go service entrypoint | Write main.go implementing the service lifecycle with signal handling, periodic sync ticker, and graceful shutdown | Task specification - File: main.go |
| F4 | must | Write config.go configuration loader | Write config.go implementing environment variable-based configuration loading for LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL | Task specification - File: config.go |
| F5 | must | Write linear.go API client | Write linear.go implementing GraphQL client for Linear.app API with issue fetching and JSON unmarshaling | Task specification - File: linear.go |
| F6 | must | Write storage.go database layer | Write storage.go implementing PostgreSQL connection management and transactional upsert logic for Linear issues | Task specification - File: storage.go |
| F7 | must | Write config_test.go unit tests | Write config_test.go with table-driven tests validating environment variable loading and error conditions | Task specification - File: config_test.go |
| F8 | must | Write linear_test.go API tests | Write linear_test.go with httptest-based tests for successful issue fetching and API error handling | Task specification - File: linear_test.go |
| F9 | must | Write storage_test.go integration tests | Write storage_test.go with PostgreSQL integration tests that skip when database is unavailable (expected behavior) | Task specification - File: storage_test.go |
| F10 | must | Build the service | Execute go build -buildvcs=false . to compile all source files into a working binary | Task specification - CRITICAL INSTRUCTIONS step 3 |
| F11 | must | Run test suite | Execute go test -buildvcs=false . -v and verify output shows 4 PASS and 2 SKIP results | Task specification - CRITICAL INSTRUCTIONS step 4 |
| F12 | must | Commit and push to GitHub | Commit all source files and push to github.com/bryanbarton525/linear-sync main branch | Task specification - CRITICAL INSTRUCTIONS step 6 |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Byte-for-byte file accuracy | All 9 files must be written VERBATIM with no character additions, removals, or modifications | Task specification - CRITICAL INSTRUCTIONS |
| NF2 | must | Build reproducibility | Build must succeed deterministically without requiring VCS metadata (buildvcs=false flag) | Task specification - CRITICAL INSTRUCTIONS step 3 |
| NF3 | must | Test execution predictability | Test suite must produce exactly 4 PASS and 2 SKIP results consistently when PostgreSQL is unavailable | Task specification - CRITICAL INSTRUCTIONS step 5 |
| NF4 | must | Dependency integrity | Dependencies specified in go.mod and checksums in go.sum must remain unmodified and not be regenerated | Task specification - CRITICAL INSTRUCTIONS step 2 |

## Dependencies

- GitHub repository github.com/bryanbarton525/linear-sync must exist and be accessible
- Git credentials must be configured for pushing to the repository
- Go 1.21+ toolchain must be available in the workspace environment
- Network access to fetch go.mod dependencies (github.com/lib/pq, github.com/stretchr/testify) if not cached

