# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Deliver a working Go service that syncs Linear.app issues to PostgreSQL by writing 9 provided source files verbatim, validating compilation and test execution, and pushing to the target GitHub repository.

## Goals

- Write all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) exactly as provided without modification
- Validate the codebase compiles successfully with go build -buildvcs=false
- Validate tests execute correctly with expected results: 4 PASS and 2 SKIP
- Commit the complete codebase to the workflow branch
- Push all changes to github.com/bryanbarton525/linear-sync main branch

## Constraints

- Files must be written verbatim - no character-level modifications permitted
- Do not run go mod init, go mod tidy, or go work init - dependency files are provided
- Must use -buildvcs=false flag for both build and test commands
- Storage tests will skip without PostgreSQL - 2 skipped tests are expected and correct
- All files must be written BEFORE running any build or test commands

## Audience

Automated workflow execution in go-orca orchestration system

## Output Medium

Git repository at github.com/bryanbarton525/linear-sync with working Go service code

## Acceptance Criteria

- All 9 files exist in workspace with exact content as specified
- go build -buildvcs=false . completes successfully with exit code 0
- go test -buildvcs=false . -v completes with exactly 4 PASS and 2 SKIP test results
- All changes committed to workflow branch with proper git metadata
- Workflow branch successfully pushed to github.com/bryanbarton525/linear-sync main branch
- No compilation errors, import errors, or syntax errors present in any file

## Out of Scope

- Modifying or improving the provided source code
- Adding additional files beyond the 9 specified
- Running PostgreSQL to enable skipped storage tests
- Refactoring code structure or applying style improvements
- Creating documentation files or README

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write all source files verbatim | Use write_file tool to create all 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) with exact content as provided in task specification | Task instruction section listing 9 files with complete content |
| F2 | must | Build the Go service | Execute 'go build -buildvcs=false .' after all files are written to validate compilation | Step 3 of CRITICAL INSTRUCTIONS |
| F3 | must | Run test suite | Execute 'go test -buildvcs=false . -v' and verify output shows 4 PASS and 2 SKIP results | Step 4-5 of CRITICAL INSTRUCTIONS |
| F4 | must | Commit changes | Commit all written files to the workflow branch with appropriate commit message | Step 6 of CRITICAL INSTRUCTIONS |
| F5 | must | Push to main branch | Push committed changes to github.com/bryanbarton525/linear-sync main branch | Step 6 of CRITICAL INSTRUCTIONS and workflow request |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | File content fidelity | Written files must match provided content exactly - no additions, deletions, or character-level changes | CRITICAL INSTRUCTIONS statement 'Write each file VERBATIM' |
| NF2 | must | Execution order | All files must be written before any go build or go test commands are executed | CRITICAL INSTRUCTIONS statement 1 |
| NF3 | must | Build validation | Build must complete without errors, warnings, or import resolution failures | Go toolchain validation profile and acceptance criteria |
| NF4 | must | Test result validation | Test execution must produce exactly 4 PASS and 2 SKIP - no failures permitted | Step 5 of CRITICAL INSTRUCTIONS |

## Dependencies

- Go toolchain (1.21 or compatible) available in execution environment
- Git configured with credentials for github.com/bryanbarton525/linear-sync
- Network connectivity to GitHub for push operation
- Write access to workspace directory /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f

---

## Constitution Amendment — Cycle 1

Implementation defect: The first task (b2f997bf) created artifact metadata claiming success but never executed the write_file tool to physically write the 9 source files to workspace /var/lib/go-orca/workspaces/ac924edc-f127-4e01-88aa-cf5ae16b892f. QA validation correctly detected empty workspace (zero packages found). The Architect must create a single remediation task that executes 9 sequential write_file tool calls with the exact file content provided in the original task specification. No requirement gaps exist - the acceptance criteria and functional requirements clearly specify all 9 files must exist in the workspace before any build or test validation can proceed.
