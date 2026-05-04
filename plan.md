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

---

## Remediation Cycle 1 — PM Triage

Implementation defect: Pod created document artifacts instead of writing actual Go source files to workspace. The Pod misinterpreted task c7844d65 and produced four document artifacts describing the work rather than using write_file to write the 9 Go source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to /var/lib/go-orca/workspaces/3cc4f0e3-7a14-4098-a4e6-d107e82cb8fc. The requirements, design, and validation are correct — this is purely an execution error where the wrong tool or artifact type was used. Architect must revise task c7844d65 to explicitly require write_file tool calls with artifact kind 'code' and workspace path, ensuring the Pod writes actual files rather than creating document summaries.

**QA blocking issues being triaged:**

- [Source Files] Workspace validation found no Go source files in /var/lib/go-orca/workspaces/3cc4f0e3-7a14-4098-a4e6-d107e82cb8fc. Pod must write the source files before Go validation can run.
- [Workspace Source Files] Toolchain validation failed: no Go source files found in workspace at /var/lib/go-orca/workspaces/3cc4f0e3-7a14-4098-a4e6-d107e82cb8fc. The Pod created 4 document artifacts describing the implementation but did not write the actual 9 Go source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to the workspace using the write_file tool. All artifacts are of kind 'document' containing summaries, not the actual source code.: Use the write_file tool to write all 9 Go source files verbatim to the workspace path /var/lib/go-orca/workspaces/3cc4f0e3-7a14-4098-a4e6-d107e82cb8fc before running any build or test commands. Each file must be written with the exact content provided in the task specification. After writing all files, the toolchain validation should pass and enable subsequent build/test operations.

---

## Remediation Cycle 1 — Architect

**Current overview:** This is a mechanical implementation task with no design decisions. The service is a Linear.app to PostgreSQL synchronization daemon that polls Linear issues every 5 minutes and upserts them into a PostgreSQL table. All 9 source files are provided verbatim and must be written exactly as specified, then built and tested to verify correctness before committing to the repository.

### Remediation Tasks

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 076e2228 | backend | Write all 9 Go source files to workspace using write_file tool | - | Use the write_file tool to write all 9 Go source files verbatim to the workspace at /var/lib/go-orca/workspaces/3cc4f0e3-7a14-4098-a4e6-d107e82cb8fc. Make all 9 write_file calls in parallel (single response). Each file must be written with the exact content provided below, with no additions, removals, or character changes. Artifact kind for each file is 'code'. Files to write: (1) go.mod with module declaration and dependencies, (2) go.sum with checksums, (3) main.go with service entry point and sync loop, (4) config.go with environment variable loading, (5) linear.go with Linear API client, (6) storage.go with PostgreSQL storage layer, (7) config_test.go with configuration tests, (8) linear_test.go with API client tests, (9) storage_test.go with storage tests. After writing all files, the workspace must contain exactly these 9 files with byte-for-byte matching content. Acceptance criteria: All 9 write_file tool calls succeed, workspace validation passes showing 9 Go source files present, no document artifacts created, files are written to /var/lib/go-orca/workspaces/3cc4f0e3-7a14-4098-a4e6-d107e82cb8fc (not to artifacts storage). |

