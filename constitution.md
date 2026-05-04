# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Deliver a production-ready Go service that syncs Linear.app issues to PostgreSQL, with full test coverage and validated build toolchain.

## Goals

- Write all 9 source files (main.go, config.go, linear.go, storage.go, and corresponding tests) to the workspace
- Validate that the code compiles successfully without errors
- Confirm that the test suite executes with 4 PASS and 2 SKIP results (storage tests skip without PostgreSQL)
- Push the working code to github.com/bryanbarton525/linear-sync on the main branch

## Constraints

- Write each file VERBATIM — no additions, removals, or character modifications
- Do not run go mod init, go mod tidy, or go work init — dependencies are pre-configured
- Use -buildvcs=false flag for build and test commands
- Expected test outcome: 4 PASS (config_test, linear_test happy paths) and 2 SKIP (storage_test without PostgreSQL environment)
- Push to main branch of github.com/bryanbarton525/linear-sync, not to the workflow branch

## Audience

Go developers and DevOps engineers deploying the Linear.app sync service

## Output Medium

Go source code repository (GitHub)

## Acceptance Criteria

- All 9 source files exist in the workspace with exact content matching provided specifications
- Command 'go build -buildvcs=false .' completes without errors
- Command 'go test -buildvcs=false . -v' shows exactly 4 PASS and 2 SKIP (no failures)
- Code is committed and pushed to github.com/bryanbarton525/linear-sync on the main branch
- Git history reflects a single commit with the service implementation

## Out of Scope

- Modifying go.mod or go.sum
- Adding new dependencies or dependency management steps
- Creating database schemas or migration scripts
- Deploying the service to any environment
- Modifying test expectations or test files

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Write source files | Write all 9 source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to the workspace directory using the exact content provided in the specification. | Task specification |
| F2 | must | Build the service | Run 'go build -buildvcs=false .' to compile the service. Build must complete without errors or warnings. | Task specification |
| F3 | must | Run tests | Run 'go test -buildvcs=false . -v' and verify output shows 4 PASS and 2 SKIP. No test failures are acceptable. | Task specification |
| F4 | must | Push to GitHub | Commit all source files with appropriate message and push to github.com/bryanbarton525/linear-sync on the main branch. | Task specification |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Test validation profile compliance | All test output must conform to the 'go' toolchain default validation profile: 4 PASS from unit tests (config_test, linear_test) and 2 SKIP from integration tests (storage_test) when PostgreSQL is unavailable. | Toolchain configuration |
| NF2 | must | Code fidelity | Every source file must match the provided specification character-for-character. No automatic formatting, linting corrections, or content modifications are permitted. | Task specification |
| NF3 | must | Dependency immutability | The go.mod and go.sum must remain unchanged. No 'go mod tidy', 'go mod init', or 'go work init' commands are allowed. | Task specification |

## Dependencies

- github.com/lib/pq (PostgreSQL driver)
- github.com/stretchr/testify (assertion library)

