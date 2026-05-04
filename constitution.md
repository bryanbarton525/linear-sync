# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Create a working Go service that synchronizes Linear.app issues to PostgreSQL with proper testing and deployment to GitHub

## Goals

- Write 9 source files verbatim to the workspace without modifications
- Successfully build the Go service using go build -buildvcs=false
- Pass all tests with go test -buildvcs=false . -v showing 4 PASS and 2 SKIP
- Commit and push all code to github.com/bryanbarton525/linear-sync main branch

## Constraints

- Do NOT run go mod init, go mod tidy, or go work init - go.mod and go.sum are already provided
- All 9 files must be written using write_file tool BEFORE running any build or test commands
- Files must be written VERBATIM - no additions, removals, or character changes
- Use -buildvcs=false flag for both build and test commands
- Storage tests will skip without PostgreSQL - this is expected and correct

## Audience

Developers and DevOps teams requiring Linear.app issue synchronization to PostgreSQL

## Output Medium

Go source code repository with build artifacts and test results

## Acceptance Criteria

- All 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are written verbatim to the workspace
- go build -buildvcs=false . completes successfully with no errors
- go test -buildvcs=false . -v produces exactly 4 PASS and 2 SKIP results
- All code is committed to the workflow branch
- All code is pushed to github.com/bryanbarton525/linear-sync main branch
- The final result passes the configured toolchain validation profile (build and test)

## Out of Scope

- Running go mod init, go mod tidy, or go work init
- Modifying the provided source files in any way
- Setting up PostgreSQL database for local testing
- Adding additional files beyond the 9 specified files
- Implementing actual Linear.app API integration testing

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write go.mod file verbatim | Write the exact go.mod file content to the workspace without any modifications | User request specifying 9 files to write verbatim |
| F2 | must | Write go.sum file verbatim | Write the exact go.sum file content to the workspace without any modifications | User request specifying 9 files to write verbatim |
| F3 | must | Write main.go file verbatim | Write the exact main.go file content containing the service entry point with signal handling and sync loop | User request specifying 9 files to write verbatim |
| F4 | must | Write config.go file verbatim | Write the exact config.go file content containing environment variable configuration loading | User request specifying 9 files to write verbatim |
| F5 | must | Write linear.go file verbatim | Write the exact linear.go file content containing Linear.app API client implementation | User request specifying 9 files to write verbatim |
| F6 | must | Write storage.go file verbatim | Write the exact storage.go file content containing PostgreSQL storage implementation | User request specifying 9 files to write verbatim |
| F7 | must | Write config_test.go file verbatim | Write the exact config_test.go file content containing configuration loading tests | User request specifying 9 files to write verbatim |
| F8 | must | Write linear_test.go file verbatim | Write the exact linear_test.go file content containing Linear API client tests | User request specifying 9 files to write verbatim |
| F9 | must | Write storage_test.go file verbatim | Write the exact storage_test.go file content containing PostgreSQL storage tests | User request specifying 9 files to write verbatim |
| F10 | must | Build the Go service | Execute go build -buildvcs=false . command after all files are written to compile the service | User request step 3 |
| F11 | must | Run all tests | Execute go test -buildvcs=false . -v command to run all unit tests | User request step 4 |
| F12 | must | Commit to workflow branch | Commit all written files to the workflow branch using git commit | User request step 6 and workflow configuration |
| F13 | must | Push to main branch | Push all committed code to github.com/bryanbarton525/linear-sync main branch | User request step 6 |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | File writing sequence | All 9 files must be written using write_file tool before any build or test commands are executed | User CRITICAL INSTRUCTIONS step 1 |
| NF2 | must | No module initialization | Must NOT run go mod init, go mod tidy, or go work init as go.mod and go.sum are already provided | User CRITICAL INSTRUCTIONS step 2 |
| NF3 | must | Verbatim file content | Files must be written exactly as specified without adding, removing, or changing any characters | User CRITICAL INSTRUCTIONS and task description |
| NF4 | must | Build success | The go build -buildvcs=false . command must complete successfully without errors | Acceptance criteria and toolchain validation profile |
| NF5 | must | Test results validation | Tests must show exactly 4 PASS and 2 SKIP results when executed with go test -buildvcs=false . -v | User CRITICAL INSTRUCTIONS step 5 |
| NF6 | must | Expected skip behavior | Storage tests skipping without PostgreSQL is expected and correct behavior | User CRITICAL INSTRUCTIONS step 5 |

## Dependencies

- write_file tool for writing all 9 source files
- go build toolchain for compilation
- go test toolchain for running tests
- git toolchain for commit and push operations
- Access to github.com/bryanbarton525/linear-sync repository

