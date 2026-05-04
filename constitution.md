# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Transcribe 9 provided Go source files for a Linear.app sync service exactly as specified, validate through the Go toolchain, and push to the GitHub repository.

## Goals

- Write all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) character-for-character as provided
- Validate the implementation builds successfully with go build -buildvcs=false
- Validate tests execute with expected results: 4 PASS, 2 SKIP using go test -buildvcs=false . -v
- Commit all files and push to github.com/bryanbarton525/linear-sync on the main branch

## Constraints

- Do NOT modify, add, or remove any characters from the provided file contents
- Do NOT run go mod init, go mod tidy, or go work init — dependencies are already specified
- Use write_file tool for all 9 files BEFORE running any build commands
- Storage tests are expected to skip without PostgreSQL — 2 SKIP results are correct

## Audience

Workflow automation system executing file transcription and validation

## Output Medium

Git repository at github.com/bryanbarton525/linear-sync, branch main

## Acceptance Criteria

- All 9 source files are written to the workspace using write_file tool with exact character-for-character match to provided content
- go build -buildvcs=false . completes successfully with exit code 0
- go test -buildvcs=false . -v produces exactly 4 PASS and 2 SKIP test results
- All files are committed to git and pushed to github.com/bryanbarton525/linear-sync main branch
- The final result passes the configured toolchain validation profile (build and test)
- The final codebase is self-contained — no placeholder text, no missing imports, no compilation errors

## Out of Scope

- Modifying the provided source code in any way
- Adding additional files beyond the 9 specified
- Configuring PostgreSQL for storage tests (tests will skip, which is expected)
- Optimizing or refactoring the provided implementation

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write all 9 source files verbatim | Use the write_file tool to create go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, and storage_test.go with exact content as provided. No modifications permitted. | User-provided file specifications |
| F2 | must | Execute go build validation | After all files are written, run go build -buildvcs=false . and verify it completes successfully with exit code 0. | Toolchain validation profile |
| F3 | must | Execute go test validation | Run go test -buildvcs=false . -v and verify output shows exactly 4 PASS results and 2 SKIP results (storage tests skip without PostgreSQL). | Toolchain validation profile |
| F4 | must | Commit and push to GitHub | Commit all files with appropriate message and push to github.com/bryanbarton525/linear-sync on the main branch. | Workflow destination requirement |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Character-perfect file transcription | All written files must be byte-for-byte identical to the provided content. No formatting changes, no line ending modifications, no whitespace adjustments. | Critical instruction: write files VERBATIM |
| NF2 | must | No Go module initialization | Do not run go mod init, go mod tidy, or go work init as dependencies are already specified in provided go.mod and go.sum. | Critical instruction constraint |
| NF3 | must | Sequential execution order | Write all files first, then run build, then run tests, then commit/push. Build and test commands must use -buildvcs=false flag. | Execution sequence requirement |

## Dependencies

- write_file tool for creating source files
- Go toolchain (go build, go test) in workspace environment
- Git repository access to github.com/bryanbarton525/linear-sync
- GitHub authentication for push access to main branch

