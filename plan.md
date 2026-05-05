# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a direct transcription workflow with verification. The task requires writing 9 pre-specified Go source files verbatim into the workspace, then executing build and test commands to verify correctness, and finally committing and pushing to the GitHub repository. No design decisions are needed—the implementation is fully specified.

## Delivery Target

GitHub repository github.com/bryanbarton525/linear-sync main branch with all source files committed

## Tech Stack

- Go 1.21
- PostgreSQL (lib/pq driver)
- testify assertion library
- Linear.app GraphQL API
- Git for version control

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source Files | Nine Go source files (main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go, go.mod, go.sum) that implement a Linear.app to PostgreSQL sync service. Each file is provided verbatim and must be written exactly as specified. | Provided file contents from task description | 9 files in workspace root |
| Build Verification | Compilation step using go build -buildvcs=false to verify all source files are syntactically correct and dependencies are properly resolved. | All 9 source files in workspace | linear-sync compiled binary |
| Test Verification | Test execution using go test -buildvcs=false . -v to verify expected test outcomes (4 PASS, 2 SKIP). | All 9 source files in workspace | Test results showing 4 PASS and 2 SKIP |
| Git Deployment | Git operations to commit all files and push to github.com/bryanbarton525/linear-sync on main branch. | All verified source files, Build success, Test success | Code pushed to GitHub repository |

## Architectural Decisions

1. **Write all 9 files before any build or test commands**
   - Rationale: Task explicitly requires writing ALL files using write_file tool before running any commands. This ensures complete file set is available for compilation.
   - Tradeoffs: Serial file writing vs parallel—chosen serial to ensure atomic completion before proceeding to build step.
2. **Use -buildvcs=false flag for both build and test**
   - Rationale: Required by task specification; prevents go toolchain from embedding VCS metadata during build.
   - Tradeoffs: None—this is a hard constraint from the task.
3. **Verify exactly 4 PASS and 2 SKIP test results**
   - Rationale: Storage tests are expected to skip without PostgreSQL database; this verifies correct test isolation behavior.
   - Tradeoffs: Strict verification ensures correct implementation but requires exact test output matching.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| f1d71a4c | backend | Write All 9 Source Files | - | Produce artifact kind `code`, names `main.go`, `config.go`, `linear.go`, `storage.go`, `config_test.go`, `linear_test.go`, `storage_test.go`, `go.mod`, `go.sum`. Write each file exactly as provided in the task description without adding, removing, or changing any characters. Use the write_file tool for each file. All files must be written to the workspace root directory (/var/lib/go-orca/workspaces/3fdc6a07-50c6-4603-96a3-e0cddab5c1b9). Do NOT run go mod init, go mod tidy, or go work init. Acceptance: All 9 files exist in workspace with exact byte-for-byte match to provided content. |
| 985b0fcd | backend | Build the Service | f1d71a4c | Run the command `go build -buildvcs=false .` in the workspace directory to compile the linear-sync service. Verify the build completes successfully with no errors or warnings. The -buildvcs=false flag is mandatory and must not be omitted. Acceptance: go build command exits with status 0 and produces a linear-sync binary in the workspace root. |
| b15c1657 | backend | Run Test Suite | f1d71a4c | Run the command `go test -buildvcs=false . -v` in the workspace directory to execute all tests. Verify the output shows exactly 4 PASS and 2 SKIP test results. The 2 SKIP results are expected from storage_test.go tests when PostgreSQL is unavailable. The -buildvcs=false flag is mandatory. Acceptance: go test command completes with exactly 4 tests passing and 2 tests skipped (TestStorage_Upsert and TestStorage_UpsertEmpty). |
| a3cd3568 | backend | Commit and Push to GitHub | 985b0fcd, b15c1657 | Commit all 9 source files to the git repository with message 'Initial commit: Linear.app sync service'. Push the commit to github.com/bryanbarton525/linear-sync on the main branch. If the workspace is not yet a git repository, run git init first. If the remote origin is not configured, add it with git remote add origin https://github.com/bryanbarton525/linear-sync. If the main branch does not exist on remote, create it with git push -u origin main. Acceptance: All files successfully pushed to github.com/bryanbarton525/linear-sync main branch and visible in the remote repository. |

