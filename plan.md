# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a mechanical file transcription task that requires writing 9 Go source files exactly as specified, validating the build and test outputs, and deploying to the target GitHub repository. The service implements a Linear.app to PostgreSQL sync mechanism with scheduled polling every 5 minutes.

## Delivery Target

Git repository at github.com/bryanbarton525/linear-sync, main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (testing assertions)
- Linear.app GraphQL API
- PostgreSQL

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source Files | Nine Go source files that comprise the linear-sync service: go.mod (module definition), go.sum (dependency checksums), main.go (entry point with sync loop), config.go (environment variable loader), linear.go (Linear.app GraphQL client), storage.go (PostgreSQL persistence layer), and three test files (config_test.go, linear_test.go, storage_test.go) | Verbatim file specifications from user | 9 source files written to workspace root |
| Build Validation | Go build step using -buildvcs=false flag to compile the service without version control metadata. Must complete with exit code 0. | All 9 source files in workspace | Compiled binary (linear-sync) or build error output |
| Test Validation | Go test execution using -buildvcs=false flag with verbose output. Expected outcome: 4 PASS (TestLoad, TestFetchIssues, TestFetchIssues_APIError, TestStorage_UpsertEmpty) and 2 SKIP (TestStorage_Upsert subtests requiring PostgreSQL). | Compiled source files from build step | Test results showing 4 PASS and 2 SKIP |
| Git Commit | Commit all 9 source files to the workflow branch with descriptive message and Co-authored-by trailer | Validated source files from test step | Git commit on workflow branch |
| Git Push | Push workflow branch to main at github.com/bryanbarton525/linear-sync | Committed changes on workflow branch | Code deployed to main branch on GitHub |

## Architectural Decisions

1. **Write all 9 files in a single parallel write_file operation**
   - Rationale: Matriarch directive explicitly requires parallelizing the write_file calls for maximum efficiency, and all 9 files are independent with no sequential dependencies
   - Tradeoffs: Single-response approach means any failure affects all files, but the verbatim nature of the task makes this acceptable
2. **Use -buildvcs=false flag for build and test commands**
   - Rationale: Specified in constitution and toolchain validation profile; required when building in environments where VCS metadata is unavailable or unreliable
   - Tradeoffs: Binary lacks VCS commit information but this is acceptable for this deployment workflow
3. **Accept 2 SKIP test results as correct behavior**
   - Rationale: Storage tests require PostgreSQL database connection which is not available in the build environment; skipping is the expected and correct behavior per constitution
   - Tradeoffs: Limited test coverage for storage layer but acceptable given the deployment-focused nature of this workflow

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 5cfc1d0e | backend | Write All 9 Source Files to Workspace | - | Produce artifact kind `code`, names: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Write all 9 files to the workspace root (/var/lib/go-orca/workspaces/718c0d44-a954-4a9e-8918-5a51ad0d18c6) in a single parallel write_file operation. Each file must be written VERBATIM - byte-for-byte identical to the specification with zero character additions, removals, or modifications. Use write_file tool for each file with the exact content provided in the task specification. Acceptance: All 9 files exist in workspace with exact content matching specification. |
| 4c602d1a | backend | Build Linear-Sync Service | 5cfc1d0e | Produce artifact kind `validation`, name `build_output.txt`. Execute `go build -buildvcs=false .` from the workspace root directory. The -buildvcs=false flag is required per toolchain validation profile. Verify the command completes with exit code 0 indicating successful compilation. Acceptance: Command exits with code 0, no compilation errors in output, binary artifact (linear-sync or linear-sync.exe) is created in the workspace. |
| 6e2562b1 | backend | Run Test Suite with Validation | 4c602d1a | Produce artifact kind `validation`, name `test_output.txt`. Execute `go test -buildvcs=false . -v` from the workspace root directory. Verify the output contains exactly 4 PASS results for: TestLoad, TestFetchIssues, TestFetchIssues_APIError, and TestStorage_UpsertEmpty. Verify the output contains exactly 2 SKIP results for TestStorage_Upsert subtests (these require PostgreSQL database connection and skipping is correct behavior per constitution). Acceptance: Test command exits with code 0, output shows 4 PASS and 2 SKIP, no test failures. |
| 401ec6e7 | backend | Commit Changes to Workflow Branch | 6e2562b1 | Produce artifact kind `git_commit`, name `workflow_branch_commit`. Stage all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) using `git add` and commit to the workflow branch (workflow/718c0d44-a954-4a9e-8918-5a51ad0d18c6) with message: 'Add Linear.app to PostgreSQL sync service' followed by a blank line and the Co-authored-by trailer: 'Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>'. Acceptance: Git commit exists on workflow branch containing all 9 files with specified commit message and trailer. |
| 842a16e9 | backend | Push to Main Branch on GitHub | 401ec6e7 | Produce artifact kind `git_push`, name `main_branch_deployment`. Push the workflow branch (workflow/718c0d44-a954-4a9e-8918-5a51ad0d18c6) to the main branch at github.com/bryanbarton525/linear-sync using the command: `git push origin workflow/718c0d44-a954-4a9e-8918-5a51ad0d18c6:main`. This pushes the local workflow branch to the remote main branch. Acceptance: Push command exits with code 0, remote main branch contains all 9 committed files, GitHub repository shows the commit. |

