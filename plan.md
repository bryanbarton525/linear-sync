# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a mechanical implementation task with no design decisions. The service is a Linear.app to PostgreSQL synchronization daemon that polls Linear issues every 5 minutes and upserts them into a PostgreSQL table. All 9 source files are provided verbatim and must be written exactly as specified, then built and tested to verify correctness before committing to the repository.

## Delivery Target

github.com/bryanbarton525/linear-sync main branch with all 9 files, successful build, and passing tests

## Tech Stack

- Go 1.21
- PostgreSQL (lib/pq driver)
- testify for assertions
- httptest for HTTP mocking

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Main Service Loop | Entry point in main.go that loads configuration, establishes database connection, and runs a ticker-based sync loop with graceful shutdown handling | Environment variables (LINEAR_API_KEY, LINEAR_TEAM_ID, DATABASE_URL) | Synchronized issues in PostgreSQL, structured logs |
| Configuration Loader | config.go reads required environment variables and returns a Config struct or error if any are missing | Environment variables | Config struct with APIKey, TeamID, and DatabaseURL |
| Linear Client | linear.go implements GraphQL client for Linear.app API to fetch team issues with their state, priority, assignee, and timestamps | API key, team ID, HTTP context | Slice of Issue structs |
| Storage Layer | storage.go handles PostgreSQL connection and upsert operations using transactions and prepared statements with JSON marshaling for assignee data | Database connection string, slice of issues | Persisted issues in linear_issues table |
| Test Suite | config_test.go, linear_test.go, and storage_test.go provide unit tests using testify assertions and httptest for the Linear client | Source code, test environment | 4 passing tests, 2 skipped storage tests (no PostgreSQL) |

## Architectural Decisions

1. **Write all files verbatim before any commands**
   - Rationale: The user explicitly requires all 9 files to be written using write_file before running build or test commands. This is a mechanical task with no interpretation needed.
   - Tradeoffs: No design flexibility, but ensures exact reproduction of the provided implementation
2. **Use -buildvcs=false flag for build and test**
   - Rationale: Workspace may not have git metadata available during build/test phase, so VCS stamping must be disabled
   - Tradeoffs: Binary will not contain VCS information, but this is acceptable for the workflow validation
3. **Accept 2 skipped storage tests**
   - Rationale: Storage tests require a live PostgreSQL instance which is not available in the test environment. Skipping is the expected and correct behavior.
   - Tradeoffs: Storage layer is not integration tested in this workflow, but unit test structure is validated

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| c7844d65 | backend | Write all 9 source files to workspace | - | Produce artifact kind `code` for all 9 files: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Each file must be written VERBATIM to the workspace at /var/lib/go-orca/workspaces/3cc4f0e3-7a14-4098-a4e6-d107e82cb8fc using the write_file tool. Do NOT add, remove, or change any characters from the provided file contents. Write all files in a single batch before proceeding to any build or test commands. Acceptance criteria: All 9 files exist in the workspace with exact byte-for-byte content matching the provided specifications. |
| 8edb890a | backend | Build the Go service | c7844d65 | Produce artifact kind `binary` by executing `go build -buildvcs=false .` in the workspace directory. This compiles the linear-sync service with VCS stamping disabled. The build must complete successfully with exit code 0 and no error output. Acceptance criteria: go build command exits successfully, binary artifact is created in the workspace. |
| 66302a9e | backend | Run all unit tests | 8edb890a | Execute `go test -buildvcs=false . -v` in the workspace directory to run the complete test suite. Verify the output shows exactly 4 PASS results (TestLoad with 4 subtests from config_test.go, TestFetchIssues and TestFetchIssues_APIError from linear_test.go) and exactly 2 SKIP results (TestStorage_Upsert and TestStorage_UpsertEmpty from storage_test.go which skip when PostgreSQL is unavailable). The 2 storage test skips are expected and correct behavior. Acceptance criteria: Test command exits with code 0, output contains 4 PASS and 2 SKIP, no FAIL results. |
| 914ba125 | backend | Commit and push to repository | 66302a9e | Commit all 9 written files to the workflow branch using `git add .` and `git commit -m "Add linear-sync service implementation"`. Then push the workflow branch to the main branch of github.com/bryanbarton525/linear-sync using `git push origin HEAD:main`. Ensure the push completes successfully. Acceptance criteria: All files are committed to the workflow branch, push to main branch completes without errors, repository at github.com/bryanbarton525/linear-sync contains all 9 files on main branch. |

