# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This is a mechanical file-writing and build-test-deploy workflow with zero design interpretation. All 9 source files are provided verbatim and must be written exactly as specified, followed by build verification, test execution, and GitHub deployment.

## Delivery Target

Repository at github.com/bryanbarton525/linear-sync on main branch

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9
- github.com/stretchr/testify v1.9.0
- PostgreSQL (production only)

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source File Initialization | Write all 9 source files verbatim to workspace root: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go | Verbatim file content from task specification | 9 Go source files and dependency manifests in workspace root |
| Build Verification | Execute go build -buildvcs=false . to verify all files compile correctly | All 9 source files | linear-sync binary, build success confirmation |
| Test Execution | Execute go test -buildvcs=false . -v to verify tests pass with expected results (4 PASS, 2 SKIP) | All source files, compiled binary | Test results showing 4 PASS and 2 SKIP |
| GitHub Deployment | Commit all files and push to github.com/bryanbarton525/linear-sync main branch | All source files, successful build and test results | Code deployed to GitHub main branch |

## Architectural Decisions

1. **Write all files in single batch before any commands**
   - Rationale: Ensures atomic file creation and prevents partial state that could cause build failures
   - Tradeoffs: None - this is the only valid approach per task specification
2. **Do not run go mod init, go mod tidy, or go work init**
   - Rationale: go.mod and go.sum are provided verbatim and must remain unmodified per constitution
   - Tradeoffs: None - modification would violate acceptance criteria
3. **Use -buildvcs=false flag for build and test commands**
   - Rationale: Required by task specification to avoid VCS metadata issues
   - Tradeoffs: Binary won't include VCS information, but this is acceptable for the workflow

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 7ca0c686 | backend | Write All 9 Source Files | - | Produce artifact kind `code`, names: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Write each file VERBATIM using the write_file tool with exact content from task specification. All files must be written to workspace root directory (/var/lib/go-orca/workspaces/2503c0ad-b42c-483b-8b92-e45c768d3ad1). Do NOT modify, add, or remove any characters from provided content. Write all 9 files in a single response before proceeding to any other task. Acceptance criteria: All 9 files exist in workspace root with byte-for-byte match to specification. |
| 229c8eed | backend | Build Go Service | 7ca0c686 | Execute build command in workspace root: go build -buildvcs=false . The build must succeed without errors and produce a linear-sync binary. Do NOT run go mod init, go mod tidy, or go work init. Acceptance criteria: Command exits with status 0, binary file 'linear-sync' exists in workspace root, no compilation errors in output. |
| 0e4e79e4 | backend | Execute Test Suite | 229c8eed | Execute test command in workspace root: go test -buildvcs=false . -v. Verify test output shows exactly 4 PASS and 2 SKIP. The 2 SKIP tests are storage tests that skip without PostgreSQL - this is expected and correct behavior per task specification. Acceptance criteria: Command exits with status 0, output contains 4 tests marked PASS, output contains 2 tests marked SKIP (TestStorage_Upsert and TestStorage_UpsertEmpty), no test failures. |
| 3e6a1d85 | backend | Commit and Push to GitHub | 0e4e79e4 | Commit all 9 source files to git repository with message 'Initial commit: Linear.app to PostgreSQL sync service' and push to github.com/bryanbarton525/linear-sync main branch. Repository URL is configured as https://github.com/bryanbarton525/linear-sync. Acceptance criteria: All files committed to git, commit pushed to main branch successfully, remote repository contains all 9 files. |

