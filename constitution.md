# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Deploy a production-ready Go service that syncs Linear.app issues to PostgreSQL with complete source code, dependency management, and automated testing.

## Goals

- Write all 9 source files exactly as specified without any modifications
- Verify the project builds successfully using the Go toolchain
- Confirm test suite executes with expected outcomes (4 PASS, 2 SKIP for database tests)
- Commit and push the complete codebase to the specified GitHub repository on the main branch

## Constraints

- All files must be written verbatim—no additions, deletions, or modifications to the provided content
- Do not run go mod init, go mod tidy, or go work init—use the provided go.mod and go.sum exactly as given
- Build command must use -buildvcs=false flag as specified
- Test command must use -buildvcs=false flag as specified
- Target repository: github.com/bryanbarton525/linear-sync
- Target branch: main

## Audience

Development team responsible for deploying and maintaining the Linear.app synchronization service

## Output Medium

Git repository containing a complete Go module with source code, tests, and dependency declarations

## Acceptance Criteria

- All 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are written to the workspace exactly as provided
- go build -buildvcs=false . executes successfully with exit code 0
- go test -buildvcs=false . -v executes successfully and produces exactly 4 PASS results and 2 SKIP results (storage tests skip when PostgreSQL is unavailable)
- All changes are committed to Git with a descriptive commit message
- Code is pushed to github.com/bryanbarton525/linear-sync on the main branch
- The final repository state contains all 9 files with identical content to the specification

## Out of Scope

- Modifying or improving the provided source code
- Adding additional dependencies or tooling
- Creating database schemas or PostgreSQL setup instructions
- Implementing additional features beyond the provided specification
- Refactoring or optimizing the code structure

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write all source files exactly as specified | Use the write_file tool to create all 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) in the workspace with content matching the specification character-for-character. No modifications, additions, or deletions are permitted. | Workflow request specification |
| F2 | must | Build the Go module | Execute 'go build -buildvcs=false .' in the workspace after all files are written. The build must complete successfully with exit code 0 and produce a valid executable binary. | Workflow request specification |
| F3 | must | Execute test suite | Run 'go test -buildvcs=false . -v' after the build succeeds. The test output must show exactly 4 PASS results (config_test.go tests and linear_test.go tests) and 2 SKIP results (storage_test.go tests that require PostgreSQL). | Workflow request specification |
| F4 | must | Commit and push to GitHub | Commit all 9 files to Git with a descriptive commit message, then push to github.com/bryanbarton525/linear-sync on the main branch. The push must succeed and the remote repository must reflect all changes. | Workflow request specification |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Build must succeed without errors | The Go build process must complete with exit code 0 and produce no compilation errors. All imports must resolve correctly using the provided go.mod and go.sum. | Go toolchain validation requirements |
| NF2 | must | Test outcomes must match specification | The test execution must produce exactly 4 PASS and 2 SKIP results. This validates that config and linear components are tested successfully, while storage tests correctly skip when PostgreSQL is unavailable in the test environment. | Workflow request specification |
| NF3 | must | File content must be exact | Each of the 9 files must be written with byte-for-byte identical content to the specification. No whitespace normalization, line ending conversion, or character encoding changes are permitted. | Workflow request specification—verbatim requirement |
| NF4 | must | Git operations must target correct repository | All Git commit and push operations must target github.com/bryanbarton525/linear-sync repository on the main branch. No other branches or repositories are acceptable. | Workflow request specification |

## Dependencies

- GitHub repository github.com/bryanbarton525/linear-sync must exist and be accessible
- Git credentials must be configured for pushing to the target repository
- Go toolchain version 1.21 or later must be available in the workspace environment

