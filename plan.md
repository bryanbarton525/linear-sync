# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a file transcription and validation task. Write 9 provided Go source files character-for-character into the workspace using parallel write_file calls, validate the implementation through go build and go test, then commit and push to the main branch of the target GitHub repository.

## Delivery Target

Git repository at github.com/bryanbarton525/linear-sync, branch main

## Tech Stack

- Go 1.21
- github.com/lib/pq (PostgreSQL driver)
- github.com/stretchr/testify (test assertions)
- Git for version control
- GitHub as repository host

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| File Transcription | Write all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) exactly as provided using write_file tool in a single batched operation | Provided file content specifications from workflow request | 9 Go source files in workspace root |
| Branch Management | Verify current git branch state and switch to main branch if currently on workflow-specific branch | Git repository state | Active main branch in workspace |
| Build Validation | Execute go build with -buildvcs=false flag and verify successful compilation | All 9 written source files | Compiled binary and exit code 0 confirmation |
| Test Validation | Execute go test with -buildvcs=false and -v flags, verify 4 PASS and 2 SKIP results | All source files including test files | Test execution log showing 4 PASS, 2 SKIP |
| Git Commit and Push | Stage all files, commit with standard message and Co-authored-by trailer, push to origin main | Validated source files on main branch | Committed changeset pushed to github.com/bryanbarton525/linear-sync main branch |

## Architectural Decisions

1. **Batch all 9 write_file calls in a single pod response**
   - Rationale: Parallel execution reduces latency and satisfies the critical requirement to write all files before running any commands
   - Tradeoffs: Single-response batching means no incremental progress visibility, but the task is deterministic transcription with no decision points
2. **Switch to main branch before commit/push rather than using PR workflow**
   - Rationale: The workflow request explicitly requires pushing to main branch; no mention of PR workflow or code review gates
   - Tradeoffs: Direct push to main assumes no branch protection rules; if push fails due to protection, the task will escalate with the exact error
3. **Accept 2 SKIP test results as valid completion**
   - Rationale: Storage tests require PostgreSQL which is not available in the test environment; skipping is expected and documented in the constitution
   - Tradeoffs: No runtime validation of storage layer, but the constitution explicitly defines 4 PASS + 2 SKIP as the acceptance baseline

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| ee18cd00 | backend | Write All 9 Source Files | - | Produce artifact kind `code` for all 9 files: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Write each file into the workspace root directory using the write_file tool with exact character-for-character content as provided in the workflow request. Use parallel write_file calls in a single response to write all files simultaneously. Acceptance criteria: All 9 files exist in workspace root with byte-for-byte match to provided content. No modifications permitted — write content verbatim including all whitespace, line endings, and formatting. |
| a73d1079 | backend | Verify and Switch to Main Branch | ee18cd00 | Execute `git status` and `git branch` to verify current branch state. If the workspace is on branch workflow/fc201a4a-... (or any branch other than main), execute `git checkout main` or `git switch main` to switch to the main branch. Acceptance criteria: The workspace is on branch main as confirmed by `git branch` output showing `* main`. If the branch switch fails with authentication or permission errors, capture the exact error message for escalation. |
| 727464e2 | backend | Execute Go Build Validation | ee18cd00 | Run `go build -buildvcs=false .` in the workspace root directory. Verify the command completes with exit code 0. Acceptance criteria: Build completes successfully with no compilation errors, producing a binary artifact (or validating that the module is buildable). The -buildvcs=false flag is required as specified in the workflow request. If build fails, capture the full error output including line numbers and error messages. |
| 9c3537a7 | backend | Execute Go Test Validation | 727464e2 | Run `go test -buildvcs=false . -v` in the workspace root directory. Verify the output shows exactly 4 PASS results and 2 SKIP results. Acceptance criteria: Test output contains 4 lines matching 'PASS: Test...' and 2 lines matching '--- SKIP: Test...' (storage tests skip without PostgreSQL, which is expected and correct). The -buildvcs=false and -v flags are required. If test results differ from 4 PASS + 2 SKIP, capture the full test output for diagnosis. |
| 040c2729 | ops | Commit and Push to GitHub Main | a73d1079, 9c3537a7 | Stage all 9 source files using `git add .`, then commit with message 'Add Linear.app sync service implementation' followed by the Co-authored-by trailer: 'Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>'. Execute `git push origin main` to push the commit to github.com/bryanbarton525/linear-sync on the main branch. Acceptance criteria: Push completes successfully with no authentication errors, no branch protection violations, and the commit appears on the main branch at the remote repository. If push fails, capture the exact error message including any hints about authentication or branch protection — do not attempt workarounds like force-push. |

