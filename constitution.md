# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Create a production-ready Go service that synchronizes Linear.app issues to PostgreSQL by writing 9 provided source files verbatim, validating the build and tests pass, and committing the result to the main branch of github.com/bryanbarton525/linear-sync.

## Goals

- Write all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) exactly as provided without any modifications
- Successfully build the service using 'go build -buildvcs=false .'
- Execute tests with 'go test -buildvcs=false . -v' and verify 4 PASS and 2 SKIP results
- Commit all source files to the main branch of github.com/bryanbarton525/linear-sync
- Ensure the service is ready for deployment with environment variable configuration

## Constraints

- All 9 files must be written verbatim without adding, removing, or changing any characters
- Do NOT run 'go mod init', 'go mod tidy', or 'go work init' — go.mod and go.sum are already provided
- Build command must be exactly: 'go build -buildvcs=false .'
- Test command must be exactly: 'go test -buildvcs=false . -v'
- The 2 skipped storage tests are expected and correct behavior when PostgreSQL is not available
- All files must be committed to the main branch, not a feature branch

## Audience

DevOps engineers and developers deploying and operating the Linear.app to PostgreSQL sync service

## Output Medium

Go source code repository at github.com/bryanbarton525/linear-sync with all files committed to the main branch

## Acceptance Criteria

- All 9 source files exist in the workspace with content matching the provided specifications exactly
- The command 'go build -buildvcs=false .' completes successfully with exit code 0
- The command 'go test -buildvcs=false . -v' completes with exactly 4 PASS results and 2 SKIP results
- All source files are committed to the main branch of github.com/bryanbarton525/linear-sync
- The go toolchain validation profile passes (build succeeds, tests pass as expected)
- No modifications, additions, or deletions were made to the provided file contents

## Out of Scope

- Modifying any of the provided source files
- Adding additional source files beyond the 9 specified
- Running go mod init, go mod tidy, or go work init
- Setting up a PostgreSQL database for integration testing
- Deploying the service to a runtime environment
- Creating Dockerfile or deployment manifests
- Adding CI/CD pipeline configuration

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write go.mod file | Create go.mod in the workspace with exact content specifying module name 'linear-sync', Go version 1.21, and dependencies (lib/pq, stretchr/testify) with exact versions and indirect dependencies. | User request: File 1 of 9 |
| F2 | must | Write go.sum file | Create go.sum in the workspace with exact checksums for all direct and indirect dependencies as provided. | User request: File 2 of 9 |
| F3 | must | Write main.go | Create main.go containing the service entry point with signal handling, configuration loading, database initialization, Linear API client, and periodic sync loop (5-minute ticker). | User request: File 3 of 9 |
| F4 | must | Write config.go | Create config.go defining Config struct and load() function that reads LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables with error handling for missing values. | User request: File 4 of 9 |
| F5 | must | Write linear.go | Create linear.go implementing linearClient with fetchIssues() method that queries Linear's GraphQL API, parses response into Issue structs with proper type definitions (Issue, Assignee, raw response types). | User request: File 5 of 9 |
| F6 | must | Write storage.go | Create storage.go implementing Storage with newDB(), newStorage(), and upsert() methods that handle PostgreSQL connection, transaction-based bulk upserts with ON CONFLICT handling, and JSON marshaling for assignee data. | User request: File 6 of 9 |
| F7 | must | Write config_test.go | Create config_test.go with table-driven test TestLoad() covering missing LINEAR_API_KEY, missing LINEAR_TEAM_ID, missing DATABASE_URL, and valid config scenarios using stretchr/testify assertions. | User request: File 7 of 9 |
| F8 | must | Write linear_test.go | Create linear_test.go with TestFetchIssues() and TestFetchIssues_APIError() using httptest.NewServer to mock Linear API responses and verify parsing and error handling. | User request: File 8 of 9 |
| F9 | must | Write storage_test.go | Create storage_test.go with TestStorage_Upsert() and TestStorage_UpsertEmpty() that skip when DATABASE_URL is unavailable, create test table, verify insert and update behavior. | User request: File 9 of 9 |
| F10 | must | Build the service | Execute 'go build -buildvcs=false .' and verify it completes with exit code 0, producing a working binary. | User request: Build step |
| F11 | must | Run tests | Execute 'go test -buildvcs=false . -v' and verify it shows 4 PASS results (config_test.go, linear_test.go tests) and 2 SKIP results (storage_test.go tests skip without PostgreSQL). | User request: Test step |
| F12 | must | Commit and push to main branch | Commit all 9 source files to the main branch of github.com/bryanbarton525/linear-sync with a descriptive commit message. | User request: Git push step |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Verbatim file content | Every file must be written with exact character-for-character fidelity to the provided content. No whitespace changes, no formatting adjustments, no additions or deletions. | User constraint: CRITICAL INSTRUCTIONS |
| NF2 | must | Build reproducibility | The build must succeed on any system with Go 1.21+ installed without requiring go mod init or go mod tidy because dependencies are already locked in go.mod and go.sum. | User constraint: Do NOT run go mod commands |
| NF3 | must | Test determinism | Test results must be consistent: exactly 4 tests pass (config and linear client tests) and 2 tests skip (storage tests without PostgreSQL). No flaky or environment-dependent failures. | User requirement: Expected test results |
| NF4 | must | Version control integration | All files must be committed to version control on the main branch, not a temporary or feature branch, enabling immediate deployment from main. | User requirement: Push to main branch |

## Dependencies

- Go toolchain version 1.21 or higher
- Git client configured with access to github.com/bryanbarton525/linear-sync
- Write access to the main branch of github.com/bryanbarton525/linear-sync
- Workspace at /var/lib/go-orca/workspaces/8e569ed1-0092-4724-bb1f-5f056debc54b
- go-orca toolchain ID 'go' with 'default' validation profile

