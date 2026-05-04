# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

Deliver a complete Linear.app-to-PostgreSQL sync service by writing 9 Go source files verbatim to the workspace, validating with build and test toolchain, and committing to the main branch. The service implements a 5-minute polling sync loop with graceful shutdown, GraphQL API client, and PostgreSQL persistence with upsert semantics.

## Delivery Target

Source code committed and pushed to github.com/bryanbarton525/linear-sync main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (test framework)
- PostgreSQL (runtime dependency)
- Linear.app GraphQL API (external dependency)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| File Materialization | Write all 9 source files exactly as specified to workspace root using write_file tool. Files: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. | Verbatim file content from task specification | 9 Go source files in /var/lib/go-orca/workspaces/04ec5122-7aac-4940-b33d-61cb7d652882 |
| Build Validation | Compile the service binary using go build with -buildvcs=false flag. Zero errors required. | 9 source files from File Materialization | Compiled binary linear-sync, zero build errors |
| Test Validation | Execute test suite with go test -buildvcs=false . -v. Validate exactly 4 PASS and 2 SKIP (storage tests skip without PostgreSQL). | 9 source files from File Materialization | Test results showing 4 PASS, 2 SKIP, zero failures |
| Git Operations | Stage all files, commit with descriptive message, push to main branch at github.com/bryanbarton525/linear-sync. Resolves branch strategy: workspace is on workflow branch but final delivery is to main branch via direct push. | Validated source files from Build and Test Validation | Committed source code on main branch with push confirmation |

## Architectural Decisions

1. **Write files directly to workspace root without module commands**
   - Rationale: go.mod and go.sum are provided complete with all dependencies pinned. Running go mod init, tidy, or work init would modify these files and violate the verbatim requirement.
   - Tradeoffs: No dependency updates possible, but ensures exact reproducibility and prevents unintended changes.
2. **Push directly to main branch from workflow branch**
   - Rationale: Task instructions explicitly state 'Commit and push to github.com/bryanbarton525/linear-sync main branch'. The workspace exists on a workflow branch but the delivery target is main. Matriarch's question about branch strategy is resolved by following the explicit instruction to push to main.
   - Tradeoffs: Bypasses pull request review workflow mentioned in Director summary, but satisfies explicit task requirement. If PR workflow is required, it must be clarified in constitution amendment.
3. **Accept 2 SKIP as passing condition**
   - Rationale: Storage integration tests skip gracefully when PostgreSQL is unavailable. This is expected and documented in acceptance criteria. 4 PASS + 2 SKIP = valid success state.
   - Tradeoffs: Reduced test coverage in CI environments without PostgreSQL, but allows development/testing without database dependency.
4. **Use -buildvcs=false flag for build and test**
   - Rationale: Prevents Git version control metadata from being embedded in binaries. Required by task specification to ensure consistent builds regardless of VCS state.
   - Tradeoffs: Binary does not contain commit hash or build metadata, but simplifies build process and prevents VCS-related failures.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| f81afda0 | backend | Write All 9 Source Files | - | Produce artifact kind `code` with 9 separate files in workspace root (/var/lib/go-orca/workspaces/04ec5122-7aac-4940-b33d-61cb7d652882). Write each file exactly as specified with zero character modifications. Files: go.mod (module definition with dependencies), go.sum (dependency checksums), main.go (service entry point with sync loop), config.go (environment variable configuration), linear.go (GraphQL API client), storage.go (PostgreSQL persistence layer), config_test.go (configuration validation tests), linear_test.go (API client tests with httptest), storage_test.go (storage integration tests with skip logic). Acceptance: All 9 files exist in workspace root with exact content matching specification. No additions, deletions, or modifications to any character. Use write_file tool for each file. Quality standards: Idiomatic Go style, proper package declarations, all imports present, no syntax errors. |
| b8c0b99e | backend | Build Service Binary | f81afda0 | Produce artifact kind `binary` named `linear-sync` by compiling the Go service. Execute command: go build -buildvcs=false . in workspace directory. Acceptance: Command exits with status 0, no error output, binary file linear-sync exists in workspace. Do NOT run go mod init, go mod tidy, or go work init. Input: 9 source files from Write All 9 Source Files task. Quality standards: Zero warnings, zero errors, clean compilation output. |
| ca3dac08 | backend | Execute Test Suite | f81afda0 | Produce artifact kind `test_results` by running the test suite. Execute command: go test -buildvcs=false . -v in workspace directory. Acceptance: Command exits with status 0, output shows exactly 4 PASS (config_test.go: 4 subtests for load validation, linear_test.go: 2 subtests for fetchIssues), exactly 2 SKIP (storage_test.go: 2 integration tests skip when PostgreSQL unavailable). No FAIL results. Input: 9 source files from Write All 9 Source Files task. Quality standards: All unit tests pass, integration tests skip gracefully with clear messages about PostgreSQL unavailability. |
| 71ecc731 | backend | Commit and Push to Main Branch | b8c0b99e, ca3dac08 | Produce artifact kind `git_commit` by staging all 9 files, creating a commit, and pushing to main branch. Stage all files with git add. Create commit with message 'Initial implementation of Linear.app to PostgreSQL sync service'. Push to remote main branch at github.com/bryanbarton525/linear-sync. Acceptance: All 9 files committed, commit exists on main branch, push succeeds with confirmation. Input: Validated source files from Build Service Binary and Execute Test Suite tasks. Quality standards: Commit message is descriptive, all files included in commit, push completes without errors. Note: Workspace is on workflow branch but delivery target is main branch per explicit task instruction. |

