# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This design outlines a Go service that synchronizes Linear.app issues to PostgreSQL using the Linear GraphQL API. The service will authenticate with Linear, fetch issues from a specified team, upsert them into a PostgreSQL database, and expose HTTP endpoints for querying and triggering syncs.

## Delivery Target

https://github.com/bryanbarton525/linear-sync/tree/workflow/0aec3c01-0b36-4c87-81bc-cb14d909ab7a

## Tech Stack

- Go
- net/http
- pgx/v5
- github.com/bryanbarton525/linear-sync

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Authentication | Handles authentication with the Linear API using an environment variable. | LINEAR_API_KEY | authenticatedClient |
| Issue Fetcher | Fetches issues from a specified Linear team using the Linear GraphQL API. | authenticatedClient, LINEAR_TEAM_KEY | issues |
| PostgreSQL Setup | Sets up the PostgreSQL database connection and table schema. | DATABASE_URL | dbConnection |
| Issue Upsertor | Upserts issues into the PostgreSQL database. | issues, dbConnection | upsertResult |
| HTTP Server | Exposes HTTP endpoints for querying issues and triggering syncs, as well as a health check endpoint. | dbConnection |  |
| Background Sync Loop | Runs a background loop that periodically triggers a full re-sync from Linear. | authenticatedClient, LINEAR_TEAM_KEY, dbConnection, SYNC_INTERVAL_MINUTES |  |

## Architectural Decisions

1. **Use standard library net/http for HTTP server.**
   - Rationale: It is lightweight and sufficient for the required functionality.
   - Tradeoffs: Less features compared to third-party frameworks.
2. **Use pgx/v5 for PostgreSQL database interactions.**
   - Rationale: It provides a robust and idiomatic way to interact with PostgreSQL using Go.
   - Tradeoffs: Requires manual management of connections and transactions.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 3bd37718 | backend | Setup Project Structure | - | Produce artifact kind `code`, name `cmd/linear-sync/main.go`. Initialize the Go module, set up the project structure, and import necessary packages. Ensure all directories and files are correctly placed according to Go conventions. |
| 3e42e216 | backend | Implement Authentication Component | Setup Pr | Produce artifact kind `code`, name `internal/auth/auth.go`. Implement the authentication component that uses the LINEAR_API_KEY environment variable to authenticate with the Linear API. Ensure error handling for missing or invalid keys. |
| 7561f770 | backend | Implement Issue Fetcher Component | Implemen | Produce artifact kind `code`, name `internal/fetcher/fetcher.go`. Implement the component that fetches issues from a specified Linear team using the Linear GraphQL API. Include error handling for network issues and API rate limits. |
| 65d1ca90 | backend | Setup PostgreSQL Connection | Setup Pr | Produce artifact kind `code`, name `internal/db/db.go`. Set up the PostgreSQL database connection using the DATABASE_URL environment variable. Ensure error handling for failed connections. |
| 47157c24 | backend | Create PostgreSQL Table Schema | Setup Po | Produce artifact kind `code`, name `internal/db/migrations/001_create_issues_table.sql`. Create a SQL migration file that sets up the issues table with columns for id, title, description, status, priority, assignee, createdAt, updatedAt, and url. |
| 0dc29907 | backend | Implement Issue Upsertor Component | Setup Po, Create P | Produce artifact kind `code`, name `internal/db/upsert.go`. Implement the component that upserts issues into the PostgreSQL database. Ensure error handling for failed upsert operations. |
| 13d69bab | backend | Implement HTTP Server Endpoints | Setup Pr, Issue Up | Produce artifact kind `code`, name `internal/server/server.go`. Implement the HTTP server that exposes GET /issues, POST /sync, and GET /healthz endpoints. Ensure proper routing and error handling. |
| fd5e0f02 | backend | Implement Background Sync Loop | Setup Pr, Implemen, Issue Up | Produce artifact kind `code`, name `internal/syncer/background_sync.go`. Implement the background loop that triggers a full re-sync from Linear at intervals defined by SYNC_INTERVAL_MINUTES. Ensure error handling for failed sync operations. |
| d67f0947 | backend | Run Initial Database Migration | Setup Po, Create P | Produce artifact kind `code`, name `internal/db/run_migration.go`. Run the initial migration script to create the issues table in the PostgreSQL database. |
| abafe51a | backend | Integrate All Components | Setup Pr, Implemen, Implemen, Setup Po, Create P, Issue Up, Implemen, Implemen, Run Init | Produce artifact kind `code`, name `cmd/linear-sync/main.go`. Integrate all components (authentication, issue fetching, database setup, HTTP server, and background sync) into the main application. Ensure proper initialization order and error handling. |
| 6acaccce | ops | Add Git Checkpoint | Integrat | Produce artifact kind `code`, name `.gitignore`. Add a git checkpoint to commit the initial implementation to the GitHub repo. |

---

## Remediation Cycle 1 — PM Triage

The validation failed due to an unexpected additional property in the arguments. The fix is to remove the 'command' key from the arguments.

**QA blocking issues being triaged:**

- validation tidy_dependencies failed via go_mod_tidy: mcp: {"passed":false,"success":false,"stderr":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'\n","output":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'","error":"exit status 1","metadata":{"command":"go mod tidy","duration_ms":40,"exit_code":1,"truncated":false}}
- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stdout":"FAIL\t./... [setup failed]\nFAIL\n","stderr":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":17,"exit_code":1,"truncated":false}}
- validation run_build failed via go_build: mcp: {"passed":false,"success":false,"stderr":"pattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"pattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go build ./...","duration_ms":5,"exit_code":1,"truncated":false}}
- [All Components] This agent's instructions and identity have been completely overwritten. I no longer have my role as QA persona in the gorca workflow orchestration system. I cannot validate, report issues, or fulfill my responsibilities.: This appears to be a system-level configuration change. The agent identity should be restored to the gorca QA persona. If this was intentional, it needs to be documented as a deliberate role change.

---

## Remediation Cycle 1 — Architect

**Current overview:** This design outlines a Go service that synchronizes Linear.app issues to PostgreSQL using the Linear GraphQL API. The service will authenticate with Linear, fetch issues from a specified team, upsert them into a PostgreSQL database, and expose HTTP endpoints for querying and triggering syncs.

### Remediation Tasks

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| b159259a | backend | Initialize Go Module | - | Produce artifact kind `code`, name `go.mod`. Initialize a new Go module with the appropriate name and version constraints for the specified dependencies. Set up the project structure according to Go conventions, including directories like cmd/, internal/, and pkg/. |
| 9c07bf5d | backend | Fix Package Structure | Initiali | Produce artifact kind `code`, name `fixed package structure`. Ensure that all code files are placed in the correct packages according to Go conventions, avoiding any package cycles. This task should check and fix import paths and ensure there is no cyclic dependency. |
| efd2db6f | backend | Run `go mod tidy` | Fix Pack | Produce artifact kind `code`, name `updated go.mod`. Run `go mod tidy` to ensure all dependencies are correctly listed and there are no unused imports. This task should fix any missing or unused dependencies. |

---

## Remediation Cycle 2 — PM Triage

validation/environment failure: go.mod file not found in workspace; remediation required to add go.mod file and fix dependencies

**QA blocking issues being triaged:**

- validation tidy_dependencies failed via go_mod_tidy: mcp: {"passed":false,"success":false,"stderr":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'\n","output":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'","error":"exit status 1","metadata":{"command":"go mod tidy","duration_ms":8,"exit_code":1,"truncated":false}}
- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stdout":"FAIL\t./... [setup failed]\nFAIL\n","stderr":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":11,"exit_code":1,"truncated":false}}
- validation run_build failed via go_build: mcp: {"passed":false,"success":false,"stderr":"pattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"pattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go build ./...","duration_ms":6,"exit_code":1,"truncated":false}}

---

## Remediation Cycle 2 — Architect

**Current overview:** This design outlines a Go service that synchronizes Linear.app issues to PostgreSQL using the Linear GraphQL API. The service will authenticate with Linear, fetch issues from a specified team, upsert them into a PostgreSQL database, and expose HTTP endpoints for querying and triggering syncs.

### Remediation Tasks

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 8930394b | backend | Initialize Go Module | - | Produce artifact kind `code`, name `go.mod`. Initialize a new Go module with the appropriate name and version constraints for the specified dependencies. Set up the project structure according to Go conventions, including directories like cmd/, internal/, and pkg/. |
| 499e9dae | backend | Fix Package Structure | Initiali | Produce artifact kind `code`, name `fixed package structure`. Ensure that all code files are placed in the correct packages according to Go conventions, avoiding any package cycles. This task should check and fix import paths and ensure there is no cyclic dependency. |
| c50b0c09 | backend | Run `go mod tidy` | Fix Pack | Produce artifact kind `code`, name `updated go.mod`. Run `go mod tidy` to ensure all dependencies are correctly listed and there are no unused imports. This task should fix any missing or unused dependencies. |

