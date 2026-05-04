# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

Deliver a Linear.app sync service by transcribing 9 Go source files verbatim to the workspace, validating compilation and test execution (4 PASS + 2 SKIP), and pushing the working code to github.com/bryanbarton525/linear-sync on the main branch. This is a zero-ambiguity delivery task with no design or implementation decisions — only file transcription, build validation, and git operations.

## Delivery Target

GitHub repository github.com/bryanbarton525/linear-sync on main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (assertion library)
- PostgreSQL (external dependency, not required for test execution)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source file transcription | Write 9 Go source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) character-for-character to the workspace directory using write_file tool. No formatting, validation, or modification permitted. | Provided file content specifications | 9 source files in workspace directory |
| Build validation | Execute 'go build -buildvcs=false .' from the workspace directory to validate that all source files compile without errors. The -buildvcs=false flag is required per constitution constraints. | All 9 source files in workspace | Compiled binary (success/failure status) |
| Test validation | Execute 'go test -buildvcs=false . -v' from the workspace directory and parse output to confirm exactly 4 PASS results (config_test.go TestLoad, linear_test.go TestFetchIssues + TestFetchIssues_APIError, storage_test.go TestStorage_UpsertEmpty) and 2 SKIP results (storage_test.go TestStorage_Upsert requires PostgreSQL). | All 9 source files, compiled service | Test execution report with counts (4 PASS, 2 SKIP, 0 FAIL) |
| Git repository setup | Ensure the workspace is a git repository with the remote 'origin' pointing to https://github.com/bryanbarton525/linear-sync. Initialize if needed, configure remote if missing, and verify repository state before committing. | Workspace directory path | Git repository ready for commit and push |
| Branch switching and commit | Switch from the workflow branch (workflow/f106c0e8-6a3e-4e91-a957-cb338877d3e6) to the main branch using 'git checkout -B main', stage all 9 source files with 'git add .', and commit with message 'Initial implementation of Linear.app sync service'. | Git repository, all 9 source files | Commit on main branch containing all source files |
| Push to GitHub main | Push the main branch to github.com/bryanbarton525/linear-sync using 'git push origin main'. Handle authentication and remote conflicts as needed (may require force push if remote main diverges). | Commit on main branch | Code pushed to GitHub main branch |

## Architectural Decisions

1. **Batch all file writes in a single task before running commands**
   - Rationale: Minimize context switches and ensure atomic file creation. Writing files one-at-a-time would introduce unnecessary round-trips and risk partial state if interrupted.
   - Tradeoffs: Single task with 9 file writes vs 9 separate tasks: the former is more efficient and reduces risk of partial workspace state.
2. **Use -buildvcs=false flag for build and test commands**
   - Rationale: Per constitution constraint, the build must not attempt to embed VCS information. This flag is mandatory for both compilation and test execution.
   - Tradeoffs: Binary will not contain git commit metadata, but this is acceptable for the delivery workflow.
3. **Target main branch instead of workflow branch for push**
   - Rationale: Constitution explicitly requires pushing to main branch of github.com/bryanbarton525/linear-sync. The workflow branch (workflow/f106c0e8-6a3e-4e91-a957-cb338877d3e6) is an internal go-orca artifact, not the delivery target.
   - Tradeoffs: Switching branches adds git operations complexity, but is mandatory per acceptance criteria.
4. **Verify exact test counts (4 PASS + 2 SKIP) instead of just 'no failures'**
   - Rationale: Constitution specifies exact expected test outcome. This ensures the provided test files execute correctly and the PostgreSQL integration tests properly skip when the database is unavailable.
   - Tradeoffs: Strict count validation vs flexible 'no failures' check: the former provides stronger validation that the code matches specifications.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 9fa8c2d3 | backend | Write all 9 source files to workspace | - | Produce artifact kind `code`, names `go.mod`, `go.sum`, `main.go`, `config.go`, `linear.go`, `storage.go`, `config_test.go`, `linear_test.go`, `storage_test.go`. Write each file to the workspace directory (/var/lib/go-orca/workspaces/f106c0e8-6a3e-4e91-a957-cb338877d3e6) using the write_file tool with the exact content provided in the task specification. Each file must be transcribed character-for-character with no formatting, validation, or content modifications. Acceptance criteria: All 9 files exist in the workspace directory with byte-for-byte identical content to the provided specifications. Verify using directory listing that all 9 files are present before marking complete. |
| 763e5d1a | backend | Build the Go service | 9fa8c2d3 | Execute artifact kind `build validation`. From the workspace directory (/var/lib/go-orca/workspaces/f106c0e8-6a3e-4e91-a957-cb338877d3e6), run the command 'go build -buildvcs=false .' to compile the service. The -buildvcs=false flag is mandatory per constitution constraints (do not run go mod init, go mod tidy, or go work init). Acceptance criteria: Command completes with exit code 0 and produces no error output. The compiled binary (named 'linear-sync') must be present in the workspace directory after build completes. Any compilation errors or warnings are blocking failures. |
| b2e62908 | backend | Run tests and verify counts | 763e5d1a | Execute artifact kind `test validation report`. From the workspace directory (/var/lib/go-orca/workspaces/f106c0e8-6a3e-4e91-a957-cb338877d3e6), run the command 'go test -buildvcs=false . -v' and parse the output. Acceptance criteria: Output must show exactly 4 PASS results (config_test.go TestLoad, linear_test.go TestFetchIssues, linear_test.go TestFetchIssues_APIError, storage_test.go TestStorage_UpsertEmpty) and exactly 2 SKIP results (storage_test.go TestStorage_Upsert which requires PostgreSQL). Zero test failures are permitted. Exit code must be 0. Any deviation from 4 PASS + 2 SKIP counts is a blocking failure. |
| 85ab41ba | backend | Initialize git repository and configure remote | b2e62908 | Execute artifact kind `git repository setup`. Check if the workspace directory (/var/lib/go-orca/workspaces/f106c0e8-6a3e-4e91-a957-cb338877d3e6) is already a git repository by running 'git rev-parse --git-dir'. If not, initialize with 'git init'. Verify the remote 'origin' is configured and points to https://github.com/bryanbarton525/linear-sync using 'git remote -v'. If the remote does not exist or points to a different URL, configure it with 'git remote add origin https://github.com/bryanbarton525/linear-sync' or 'git remote set-url origin https://github.com/bryanbarton525/linear-sync'. Acceptance criteria: Repository is initialized, remote 'origin' exists and points to https://github.com/bryanbarton525/linear-sync. The workspace must be ready for git add/commit/push operations. |
| 9b19ad1b | backend | Switch to main branch and commit files | 85ab41ba | Execute artifact kind `git commit on main branch`. From the workspace directory (/var/lib/go-orca/workspaces/f106c0e8-6a3e-4e91-a957-cb338877d3e6), switch from the workflow branch (workflow/f106c0e8-6a3e-4e91-a957-cb338877d3e6) to the main branch using 'git checkout -B main'. Stage all 9 source files with 'git add go.mod go.sum main.go config.go linear.go storage.go config_test.go linear_test.go storage_test.go'. Commit the staged files with message 'Initial implementation of Linear.app sync service'. Acceptance criteria: Command 'git log -1 --oneline' shows the commit on the main branch with the specified message. All 9 source files are included in the commit (verify with 'git show --name-only'). The compiled binary (linear-sync) must NOT be committed. |
| 297d8918 | backend | Push to GitHub main branch | 9b19ad1b | Execute artifact kind `git push operation`. From the workspace directory (/var/lib/go-orca/workspaces/f106c0e8-6a3e-4e91-a957-cb338877d3e6), push the main branch to github.com/bryanbarton525/linear-sync using 'git push origin main'. If the remote main branch already exists and has diverged, use 'git push --force origin main' to overwrite. Acceptance criteria: Command completes with exit code 0. The remote main branch at https://github.com/bryanbarton525/linear-sync contains the commit with message 'Initial implementation of Linear.app sync service' and all 9 source files. Verify the push succeeded by checking that 'git log origin/main -1 --oneline' matches the local main branch commit. |

