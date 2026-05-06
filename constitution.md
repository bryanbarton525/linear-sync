# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Create a production-ready Go service that synchronizes Linear.app issues to PostgreSQL, providing persistent storage and queryability for Linear issue data outside the Linear platform.

## Goals

- Write exactly 9 Go source files verbatim to the workspace without modification
- Build the service successfully using the provided go.mod and go.sum
- Execute tests and verify expected results (4 PASS, 2 SKIP for storage tests)
- Commit all source files and push to github.com/bryanbarton525/linear-sync on the main branch

## Constraints

- Files must be written verbatim - no additions, removals, or character changes
- Do NOT run go mod init, go mod tidy, or go work init - go.mod and go.sum are already provided
- Build command must use -buildvcs=false flag
- Test command must use -buildvcs=false flag
- Storage tests are expected to skip without PostgreSQL (2 SKIP is correct behavior)
- All files must be written BEFORE running any build or test commands

## Audience

Development team maintaining the linear-sync service and repository maintainers at github.com/bryanbarton525/linear-sync

## Output Medium

Git repository at github.com/bryanbarton525/linear-sync, main branch

## Acceptance Criteria

- All 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) are written to the workspace verbatim
- Command 'go build -buildvcs=false .' completes successfully with exit code 0
- Command 'go test -buildvcs=false . -v' shows exactly 4 PASS and 2 SKIP
- All source files are committed to the workflow branch
- Code is pushed to github.com/bryanbarton525/linear-sync on the main branch
- Repository contains all 9 files at the specified paths in the root directory

## Out of Scope

- Modifying or improving the provided source code
- Adding additional files beyond the specified 9
- Running PostgreSQL for storage tests (SKIP is expected and correct)
- Running go mod tidy or similar dependency management commands
- Setting up CI/CD pipelines or deployment configurations

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write 9 source files verbatim | Use write_file tool to write all 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to the workspace exactly as specified, preserving every character, without any modifications | User request explicit instruction #1 |
| F2 | must | Build the Go service | After all files are written, execute 'go build -buildvcs=false .' command to build the linear-sync service binary | User request explicit instruction #3 |
| F3 | must | Run test suite | Execute 'go test -buildvcs=false . -v' command to run all tests and verify 4 PASS and 2 SKIP results | User request explicit instruction #4 |
| F4 | must | Commit and push to GitHub | Commit all 9 source files to the workflow branch and push to github.com/bryanbarton525/linear-sync on the main branch | User request explicit instruction #6 |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Byte-for-byte file accuracy | Each of the 9 files must match the provided content exactly - no trailing whitespace, encoding changes, or character substitutions | User request critical instruction: 'Write each file VERBATIM' |
| NF2 | must | Build success without errors | The go build command must complete with exit code 0 and produce a valid binary | Acceptance criteria for software workflows |
| NF3 | must | Test results validation | Test output must show exactly 4 tests passing and 2 tests skipping (storage tests skip when PostgreSQL is unavailable, which is expected and correct) | User request explicit instruction #5 |
| NF4 | must | Execution order constraint | All 9 files must be written to disk before executing any go build or go test commands | User request critical instruction #1 |
| NF5 | must | No dependency management | Must NOT run go mod init, go mod tidy, or go work init as dependencies are pre-configured in go.mod and go.sum | User request explicit instruction #2 |

## Dependencies

- Go 1.21 toolchain installed and available in PATH
- Git configured with credentials for github.com/bryanbarton525/linear-sync
- Network connectivity to github.com for git push
- Workspace at /var/lib/go-orca/workspaces/cb5e88ce-5a8f-4951-9893-9e07c607d87d with write permissions

