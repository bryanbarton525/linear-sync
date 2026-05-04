# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

A Go microservice that periodically fetches Linear.app issues via GraphQL API and persists them to PostgreSQL. The service runs a 5-minute sync loop with graceful shutdown, structured logging, and upsert-based storage to handle issue updates. All configuration is environment-driven with comprehensive error handling and context propagation.

## Delivery Target

GitHub repository with compilable Go source code and passing test suite

## Tech Stack

- Go 1.21
- PostgreSQL 11+
- Linear GraphQL API
- github.com/lib/pq
- github.com/stretchr/testify

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Configuration Loader | Reads LINEAR_API_KEY, LINEAR_TEAM_ID, and DATABASE_URL from environment variables with validation | Environment variables | Config struct |
| Linear GraphQL Client | Fetches issues from Linear.app API using team-scoped GraphQL queries with JSON response parsing | API key, Team ID, HTTP context | Issue slice, API errors |
| PostgreSQL Storage Layer | Handles database connections and upsert operations with transaction safety and JSON serialization for assignee data | Database URL, Issue slice | Persistence confirmation, Database errors |
| Sync Orchestrator | Main loop with 5-minute ticker, timeout enforcement, signal handling, and structured logging | Configuration, Client, Storage | Sync status logs, Graceful shutdown |

## Architectural Decisions

1. **Environment-based configuration**
   - Rationale: Follows 12-factor app principles for containerized deployment and prevents secrets in source code
   - Tradeoffs: Less flexible than file-based config but more secure and cloud-native
2. **GraphQL for Linear API integration**
   - Rationale: Linear.app's native API format allows precise field selection and efficient data fetching
   - Tradeoffs: More complex than REST but reduces over-fetching and supports Linear's schema evolution
3. **Upsert-based PostgreSQL storage**
   - Rationale: Handles both new issues and updates in a single operation, preventing duplicate key errors
   - Tradeoffs: Requires unique constraints but eliminates complex conflict resolution logic
4. **Fixed 5-minute sync interval**
   - Rationale: Balances API rate limits with reasonable data freshness for integration use cases
   - Tradeoffs: Not configurable but prevents misconfiguration and API abuse

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| bb712ae4 | backend | Transcribe Go Source Files | - | Write all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to the workspace exactly as provided. Each file must be written verbatim without any character changes, additions, or deletions. The files implement a complete Linear-to-PostgreSQL sync service with configuration loading, GraphQL client, storage layer, and comprehensive unit tests. Use write_file tool for each source file to ensure exact transcription. |
| 2ddae52f | backend | Build and Test Linear Sync Service | bb712ae4 | Execute build and test commands to validate the transcribed Go source code. Run 'go build -buildvcs=false .' to compile the service and verify no build errors. Then run 'go test -buildvcs=false . -v' to execute the test suite. Verify the output shows exactly 4 PASS results and 2 SKIP results (storage tests skip without PostgreSQL environment, which is expected behavior). The build must produce a working executable and all non-skipped tests must pass without errors or panics. |
| 17e9e48e | backend | Commit and Push to GitHub Repository | 2ddae52f | Commit all 9 source files to the git repository with the message 'Add Linear.app to PostgreSQL sync service (workflow 13ebcc13-3f88-4186-bc3c-2ecb581e4ceb)' and include the Co-authored-by trailer 'Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>'. Push the commit to the main branch of github.com/bryanbarton525/linear-sync. The repository workspace is already configured, so use standard git add, git commit, and git push commands. Verify the push succeeds and the remote repository reflects the exact state of the 9 provided files. |

---

## Remediation Cycle 1 — PM Triage

Classification: **implementation defect**. The 9 required Go source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) were not written to the workspace during the implementation phase, causing the test validation to fail with 'no packages to test'. The constitution and requirements are correct and complete; no amendments needed. Remediation: Execute write_file operations to transcribe all 9 files exactly as provided in the task specification. After files are written, rerun build (`go build -buildvcs=false .`) and tests (`go test -buildvcs=false . -v`) to verify 4 PASS and 2 SKIP results, then commit and push to the GitHub repository.

**QA blocking issues being triaged:**

- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stderr":"go: warning: \"./...\" matched no packages\nno packages to test\n","output":"go: warning: \"./...\" matched no packages\nno packages to test","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":6,"exit_code":1,"truncated":false}}
- [workspace files] Source files were not written to the workspace. Validation shows 'go test ./...' failed with 'no packages to test', indicating the 9 required Go source files are missing from the workspace directory.: Execute write_file operations for all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) exactly as provided in the task specification before attempting build and test commands.

---

## Remediation Cycle 1 — Architect

**Current overview:** A Go microservice that periodically fetches Linear.app issues via GraphQL API and persists them to PostgreSQL. The service runs a 5-minute sync loop with graceful shutdown, structured logging, and upsert-based storage to handle issue updates. All configuration is environment-driven with comprehensive error handling and context propagation.

### Remediation Tasks

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| f3ae9c09 | backend | Write All 9 Go Source Files to Workspace | - | Produce artifact kind `code`, name `go-source-files`. Use write_file tool to transcribe all 9 source files exactly as provided in the task specification: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Each file must be written verbatim to the workspace root directory without any character modifications, additions, or deletions. Verify all files are written successfully before proceeding. This resolves the QA blocking issue where source files were missing from the workspace, causing 'no packages to test' error. |
| 553c8810 | backend | Build and Test Linear Sync Service | f3ae9c09 | Produce artifact kind `validation`, name `build-test-results`. Execute build command 'go build -buildvcs=false .' to compile the service and verify no build errors. Then execute test command 'go test -buildvcs=false . -v' to run the test suite. Verify the output shows exactly 4 PASS results and 2 SKIP results (storage tests skip without PostgreSQL environment, which is expected behavior). The build must produce a working executable and all non-skipped tests must pass without errors or panics. Report any deviations from expected test results. |
| 88284673 | backend | Commit and Push to GitHub Repository | 553c8810 | Produce artifact kind `git-commit`, name `linear-sync-commit`. Commit all 9 source files to the git repository with commit message 'Add Linear.app to PostgreSQL sync service (workflow 13ebcc13-3f88-4186-bc3c-2ecb581e4ceb)' and include Co-authored-by trailer 'Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>'. Push the commit to the main branch of github.com/bryanbarton525/linear-sync. Use standard git add, git commit, and git push commands. Verify the push succeeds and the remote repository reflects the exact state of the 9 provided files. |

---

## Remediation Cycle 2 — PM Triage

The 9 Go source files were not written to the workspace directory. QA validation reports 'no packages to test', indicating the workspace is empty despite prior artifact claims of successful file transcription. Remediation: Use write_file tool to transcribe all 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to the workspace root (/var/lib/go-orca/workspaces/13ebcc13-3f88-4186-bc3c-2ecb581e4ceb/) exactly as provided in the task specification. Verify files exist in workspace before proceeding to build and test commands.

**QA blocking issues being triaged:**

- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stderr":"go: warning: \"./...\" matched no packages\nno packages to test\n","output":"go: warning: \"./...\" matched no packages\nno packages to test","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":7,"exit_code":1,"truncated":false}}

