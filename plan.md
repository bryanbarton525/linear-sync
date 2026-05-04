# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a verbatim transcription and validation workflow. The solution consists of writing 9 provided Go source files exactly as specified to the workspace, validating compilation with the Go toolchain, executing tests to verify expected outcomes (4 PASS, 2 SKIP), and pushing the complete codebase to the target GitHub repository. No design decisions are required - all source code is provided in final form.

## Delivery Target

Git repository at github.com/bryanbarton525/linear-sync with all source files committed to main branch, validated to compile and test successfully

## Tech Stack

- Go 1.21
- github.com/lib/pq (PostgreSQL driver)
- github.com/stretchr/testify (test assertions)
- PostgreSQL (runtime dependency, not required for build/test)
- Linear.app GraphQL API (runtime dependency)
- Git (for commit and push operations)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source File Set | 9 Go source files comprising the Linear.app sync service: go.mod (module definition), go.sum (dependency checksums), main.go (service entry point with sync loop), config.go (environment variable loader), linear.go (Linear API client), storage.go (PostgreSQL persistence layer), config_test.go (configuration tests), linear_test.go (API client tests), storage_test.go (database tests) | Provided file content from task specification | 9 files in workspace root directory |
| Build Validation | Go compilation check using 'go build -buildvcs=false .' to verify all imports resolve, syntax is correct, and the codebase compiles to a valid binary | All 9 source files in workspace | Compiled binary (exit code 0 indicates success) |
| Test Validation | Test execution using 'go test -buildvcs=false . -v' to verify 4 tests pass (config tests and linear client tests) and 2 tests skip (storage tests requiring PostgreSQL) | All source files and successful build | Test results: 4 PASS, 2 SKIP |
| Git Commit | Commit all 9 files to the workflow branch with descriptive message and Co-authored-by trailer | All validated source files | Git commit on workflow branch |
| Git Push | Push workflow branch to github.com/bryanbarton525/linear-sync main branch | Committed changes on workflow branch | Changes visible on main branch in GitHub repository |

## Architectural Decisions

1. **Write all files before any build/test commands**
   - Rationale: Task specification explicitly requires all 9 files be written BEFORE running go build or go test to ensure complete codebase is available for validation
   - Tradeoffs: Sequential execution ensures file availability but increases task latency compared to interleaved write/validate approach
2. **Use -buildvcs=false flag for build and test**
   - Rationale: Task specification requires this flag to disable VCS stamping during build, likely due to workflow branch not being the main branch at time of build
   - Tradeoffs: Binary will not contain VCS metadata (commit hash, timestamp) but ensures build succeeds in workflow context
3. **Accept 2 skipped tests as success condition**
   - Rationale: Storage tests require PostgreSQL connection which is not available in build environment - skipped tests are expected and correct behavior
   - Tradeoffs: Partial test coverage in CI environment, but full coverage available when PostgreSQL is configured in production
4. **Do not run go mod init/tidy/work init**
   - Rationale: Task specification explicitly prohibits these commands because go.mod and go.sum are provided in final form and must not be regenerated
   - Tradeoffs: Cannot fix dependency issues if they exist, but ensures exact dependency graph specified in provided files

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| b2f997bf | backend | Write all 9 source files verbatim | - | Produce artifact kind `code` comprising 9 Go source files in workspace root directory /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f. Write each file with EXACT content as provided in task specification - no character-level modifications permitted. Files to write: (1) go.mod with module definition and dependencies, (2) go.sum with dependency checksums, (3) main.go with service entry point and sync loop, (4) config.go with environment variable loader, (5) linear.go with Linear API client, (6) storage.go with PostgreSQL persistence layer, (7) config_test.go with configuration tests, (8) linear_test.go with API client tests, (9) storage_test.go with database tests. Use write_file tool for each file. Acceptance criteria: All 9 files exist in workspace root with byte-for-byte identical content to provided specifications, file paths are relative to workspace root (e.g. 'go.mod', 'main.go'), no additional files created, no existing files modified. Quality standard: Verbatim transcription - any deviation from provided content is a failure. |
| ae0d0423 | backend | Validate Go build | b2f997bf | Produce artifact kind `build_validation` by executing 'go build -buildvcs=false .' in workspace directory /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f. This task validates that all 9 source files (written by previous task) compile successfully with Go 1.21 toolchain. Use bash tool to run build command. Acceptance criteria: Command exits with code 0, no compilation errors in stderr, no import resolution failures, binary artifact 'linear-sync' created in workspace root. If build fails with any errors, stop immediately and report exact error message - do not attempt to modify source files. Quality standard: Zero compilation errors - any syntax error, import error, or type error is a blocking failure. |
| d0f5cc43 | backend | Validate test execution | ae0d0423 | Produce artifact kind `test_validation` by executing 'go test -buildvcs=false . -v' in workspace directory /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f. This task validates that test suite executes with expected results. Use bash tool to run test command and capture stdout. Acceptance criteria: Command exits with code 0, stdout contains exactly 4 lines with 'PASS:' prefix (TestLoad with 4 subtests, TestFetchIssues, TestFetchIssues_APIError, one of TestStorage_UpsertEmpty), stdout contains exactly 2 lines with '--- SKIP:' prefix (TestStorage_Upsert and TestStorage_UpsertEmpty skipping due to missing PostgreSQL connection), no lines with 'FAIL:' prefix. If test results differ from 4 PASS and 2 SKIP, stop immediately and report discrepancy - do not modify test files. Quality standard: Exact match on test result counts - any FAIL result or different PASS/SKIP counts is a blocking failure. |
| 3eb7a7c6 | ops | Commit source files to workflow branch | d0f5cc43 | Produce artifact kind `git_commit` by committing all 9 source files to workflow branch 'workflow/ac924edc-f127-4e01-88aa-cf5ae16b892f' in workspace /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f. Use bash tool to execute git commands. Stage all files with 'git add go.mod go.sum main.go config.go linear.go storage.go config_test.go linear_test.go storage_test.go'. Commit with message 'Add Linear.app to PostgreSQL sync service' followed by blank line and Co-authored-by trailer 'Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>'. Acceptance criteria: 'git status' shows working tree clean after commit, 'git log -1' shows commit with correct message and trailer, all 9 files included in commit. Quality standard: Commit must include all files and correct metadata per git commit trailer requirement. |
| 1e117b3f | ops | Push workflow branch to main | 3eb7a7c6 | Produce artifact kind `git_push` by pushing workflow branch 'workflow/ac924edc-f127-4e01-88aa-cf5ae16b892f' to 'main' branch on remote repository github.com/bryanbarton525/linear-sync. Use bash tool to execute 'git push origin workflow/ac924edc-f127-4e01-88aa-cf5ae16b892f:main'. Verify git credentials are configured before attempting push. Acceptance criteria: Push completes successfully with 'Branch workflow/ac924edc-f127-4e01-88aa-cf5ae16b892f set up to track remote branch main' or similar success message, no authentication errors, no rejected push errors, changes visible on main branch at github.com/bryanbarton525/linear-sync when queried with 'git ls-remote origin main'. If push fails with authentication error, report error immediately - credentials must be configured outside this workflow. Quality standard: Successful push to main branch with all commits from workflow branch - any network error, auth error, or rejected push is a blocking failure. |

---

## Remediation Cycle 1 — PM Triage

Implementation defect: The first task (b2f997bf) created artifact metadata claiming success but never executed the write_file tool to physically write the 9 source files to workspace /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f. QA validation correctly detected empty workspace (zero packages found). The Architect must create a single remediation task that executes 9 sequential write_file tool calls with the exact file content provided in the original task specification. No requirement gaps exist - the acceptance criteria and functional requirements clearly specify all 9 files must exist in the workspace before any build or test validation can proceed.

**QA blocking issues being triaged:**

- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stderr":"go: warning: \"./...\" matched no packages\nno packages to test\n","output":"go: warning: \"./...\" matched no packages\nno packages to test","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":7,"exit_code":1,"truncated":false}}
- [Source Files] Workspace validation shows zero packages exist - 'go test ./...' returned 'matched no packages' and 'no packages to test'. The 9 required source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) were never written to workspace /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f.: Use write_file tool to write all 9 files to workspace directory with exact content from task specification. Files must exist before any validation can succeed.

