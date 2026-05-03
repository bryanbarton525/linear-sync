# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

The design involves cloning an existing Go service that syncs Linear.app issues to PostgreSQL via the Linear GraphQL API. The service will poll issues every 5 minutes and handle errors and rate limits gracefully. The workflow ensures that the code builds, runs tests (with expected skips), and commits to GitHub without modifying existing file structures.

## Delivery Target

GitHub repository

## Tech Stack

- Go
- PostgreSQL
- Linear GraphQL API

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| linear-sync | A Go service that syncs Linear.app issues to PostgreSQL using the Linear GraphQL API. | Linear.app credentials, PostgreSQL connection string | Synced issues in PostgreSQL |
| config | Configuration management for linear-sync, loading configuration from environment variables and a config file. | Environment variables, Config file | *Config object |
| linear-client | Integration with the Linear GraphQL API to fetch issues. | API key, Team ID | List of issues |
| storage | Storage management for syncing issues to PostgreSQL. | List of issues | Upserted issues in PostgreSQL |

## Architectural Decisions

1. **Clone the existing linear-sync repository without changes.**
   - Rationale: The code is already complete and correct as per PM validation.
   - Tradeoffs: No changes means no risk of introducing new bugs.
2. **Run tests expecting skips due to missing PostgreSQL.**
   - Rationale: Storage tests will skip without a running PostgreSQL instance, which is acceptable per PM constitution.
   - Tradeoffs: Skipped tests may not cover storage functionality, but they are expected and acceptable for this workflow.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 83413eb6 | backend | Clone linear-sync repository | - | Produce artifact kind `code`, name `linear-sync`. Clone the existing linear-sync repository from https://github.com/bryanbarton525/linear-sync into the workspace. Do not modify any file structure or add unused imports. |
| 6c458c6c | backend | Run tests | 83413eb6 | Produce artifact kind `test_results`, name `run_tests`. Run all Go tests in the workspace using `go test ./...`. Expect storage tests to skip due to missing PostgreSQL. Ensure all other tests pass. |
| 477c41e1 | backend | Build service | 6c458c6c | Produce artifact kind `build_results`, name `run_build`. Build the linear-sync service using `go build ./...` to ensure it compiles successfully. |
| de15a9c8 | backend | Stage and commit changes | 477c41e1 | Produce artifact kind `commit`, name `stage_and_commit`. Stage all files in the workspace and commit them to the branch `workflow/e297036a-25fb-4408-9394-57ad85bccc51` with a message indicating completion of the workflow. |
| 8f0f61c1 | backend | Push to GitHub | de15a9c8 | Produce artifact kind `push`, name `push_to_github`. Push the committed changes from the workflow branch to the remote repository at https://github.com/bryanbarton525/linear-sync. |

