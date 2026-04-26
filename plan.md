# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

The service will authenticate with Linear using an API key, fetch issues from a specified team using the Linear GraphQL API, and upsert them into a PostgreSQL database. It will expose HTTP endpoints for querying and syncing issues, and run a background sync loop at configurable intervals.

## Delivery Target

git@github.com:bryanbarton525/linear-sync.git

## Tech Stack

- Go
- net/http
- pgx/v5
- github.com/bryanbarton525/linear-sync

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| LinearClient | Wraps the Linear GraphQL API with context propagation. | LINEAR_API_KEY | authToken, issues |
| PostgreSQLStorage | Handles database connection and upsert operations for issues. | DATABASE_URL, issues |  |
| HTTPServer | Exposes endpoints for querying and syncing issues, and health checks. | issues |  |
| BackgroundSync | Runs a background loop to sync issues at specified intervals. | LINEAR_TEAM_KEY, SYNC_INTERVAL_MINUTES |  |

## Architectural Decisions

1. **Use pgx/v5 pool for PostgreSQL operations.**
   - Rationale: pgx/v5 is a high-performance PostgreSQL driver with connection pooling support, suitable for the task.
   - Tradeoffs: None significant.
2. **Implement background sync using errgroup.Group.**
   - Rationale: errgroup provides graceful shutdown and context propagation for background tasks.
   - Tradeoffs: Requires careful management of contexts to avoid deadlocks.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| c6744f14 | backend | Create LinearClient | - | Produce artifact kind `code`, name `linear_client.go`. Implement the LinearClient struct to authenticate with LINEAR_API_KEY and fetch issues using the Linear GraphQL API. Include context propagation in all API calls. Ensure the client is reusable for fetching multiple queries. |
| fb1b2b44 | backend | Create PostgreSQLStorage | - | Produce artifact kind `code`, name `postgres_storage.go`. Implement the PostgreSQLStorage struct to handle database connection using DATABASE_URL and upsert issues into a PostgreSQL table. Define the table schema with fields: id, title, description, status/state name, priority, assignee name, createdAt, updatedAt, url. |
| 0e4e7cad | backend | Create HTTPServer | - | Produce artifact kind `code`, name `http_server.go`. Implement the HTTPServer struct to expose the following endpoints: GET /issues (returns issues as JSON), POST /sync (triggers a full re-sync), and GET /healthz (service health check). Ensure all handlers use context for cancellation propagation. |
| 3e181781 | backend | Create BackgroundSync | - | Produce artifact kind `code`, name `background_sync.go`. Implement the BackgroundSync struct to run a ticker-based sync loop every SYNC_INTERVAL_MINUTES (default 15 minutes). Use errgroup.Group for graceful shutdown and context propagation. Ensure the background task fetches issues using LinearClient and upserts them using PostgreSQLStorage. |
| b2a250ff | backend | Create main.go | Create L, Create P, Create H, Create B | Produce artifact kind `code`, name `main.go`. Implement the entry point of the application. Initialize LinearClient, PostgreSQLStorage, HTTPServer, and BackgroundSync. Validate environment variables LINEAR_API_KEY, LINEAR_TEAM_KEY, DATABASE_URL, and SYNC_INTERVAL_MINUTES at startup. Start the HTTP server and background sync loop. |

---

## Remediation Cycle 1 — PM Triage

QA blocking issues: 1. Missing go.mod file (requirement gap), 2. Duplicate function (implementation defect), 3. Incomplete struct initialization (implementation defect), 4. Field name inconsistency (implementation defect), 5. Inconsistent package names (design gap), 6. Server handlers not properly initialized (implementation defect).

**QA blocking issues being triaged:**

- validation tidy_dependencies failed via go_mod_tidy: mcp: {"passed":false,"success":false,"stderr":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'\n","output":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'","error":"exit status 1","metadata":{"command":"go mod tidy","duration_ms":49,"exit_code":1,"truncated":false}}
- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stdout":"FAIL\t./... [setup failed]\nFAIL\n","stderr":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":47,"exit_code":1,"truncated":false}}
- validation run_build failed via go_build: mcp: {"passed":false,"success":false,"stderr":"pattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"pattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go build ./...","duration_ms":7,"exit_code":1,"truncated":false}}
- [Module Initialization] Missing go.mod file in workspace. All validation steps (go mod tidy, go test, go build) fail because the workspace is not a valid Go module. This is a fundamental setup issue.: Create go.mod file with module path github.com/bryanbarton525/linear-sync and run go mod tidy to add dependencies
- [linear_client.go] Duplicate NewLinearClient function defined twice in the same file (lines ~62 and ~104). This violates Go conventions.: Remove duplicate function definition, keep only one
- [postgres_storage.go] Struct initialization incomplete: 'pool:' is missing closing brace. Also missing 'sync' import despite using sync.Mutex in struct.: Complete struct initialization with proper closing brace and add missing import
- [PostgreSQL Storage] Field name inconsistency: Issue struct uses StatusName while linear_client.go uses StatusID. This will cause mismatch between client and storage.: Align field names between linear_client.go and postgres_storage.go Issue structs
- [Package Organization] Inconsistent package names across artifacts (linear, storage, orca, main). Orca package does not match other package conventions.: Use consistent package names that match module structure: 'internal/linear', 'internal/storage', 'internal/server'
- [HTTP Server] Server handlers never properly initialized. HandleIssues, handleSync methods receive nil handlers that are never set. Healthz handler also has method routing issues.: Properly wire handlers or use ServeMux pattern correctly. Fix handleIssues, handleSync methods to accept proper handlers

---

## Remediation Cycle 1 — Architect

**Current overview:** The service will authenticate with Linear using an API key, fetch issues from a specified team using the Linear GraphQL API, and upsert them into a PostgreSQL database. It will expose HTTP endpoints for querying and syncing issues, and run a background sync loop at configurable intervals.

### Remediation Tasks

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 9ce013dd | backend | Create go.mod file | - | Produce artifact kind `code`, name `go.mod`. Create a go.mod file with module path github.com/bryanbarton525/linear-sync. Run go mod tidy to add dependencies. |
| 9d6db9c9 | backend | Remove duplicate NewLinearClient function | Create g | Produce artifact kind `code`, name `linear_client.go`. Remove the duplicate NewLinearClient function definition, keeping only one. |
| 3c7f315f | backend | Complete PostgreSQLStorage struct initialization | Create g | Produce artifact kind `code`, name `postgres_storage.go`. Complete the struct initialization with proper closing brace and add missing import for sync.Mutex. |
| 4ec3c35f | backend | Align field names between linear_client.go and postgres_storage.go Issue structs | Create g | Produce artifact kind `code`, name `linear_client.go` and `postgres_storage.go`. Align the field names between the two Issue structs (e.g., use StatusName in both). |
| 5890f394 | backend | Rename packages to match module structure | Create g | Produce artifact kind `code`, name `linear_client.go` and other related files. Rename the package names to 'internal/linear', 'internal/storage', 'internal/server' as per Go convention. |
| adb8b933 | backend | Properly wire HTTP server handlers | Create g | Produce artifact kind `code`, name `http_server.go`. Properly wire the HTTP server handlers to accept proper handlers. Ensure HandleIssues, handleSync methods are correctly set and wired. |

---

## Remediation Cycle 2 — PM Triage

QA blocking issues: 1. Missing go.mod file (requirement gap), 2. Duplicate function (implementation defect), 3. Incomplete struct initialization (implementation defect), 4. Field name inconsistency (implementation defect), 5. Inconsistent package names (design gap), 6. Server handlers not properly initialized (implementation defect).

**QA blocking issues being triaged:**

- validation tidy_dependencies failed via go_mod_tidy: mcp: {"passed":false,"success":false,"stderr":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'\n","output":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'","error":"exit status 1","metadata":{"command":"go mod tidy","duration_ms":5,"exit_code":1,"truncated":false}}
- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stdout":"FAIL\t./... [setup failed]\nFAIL\n","stderr":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":10,"exit_code":1,"truncated":false}}
- validation run_build failed via go_build: mcp: {"passed":false,"success":false,"stderr":"pattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"pattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go build ./...","duration_ms":6,"exit_code":1,"truncated":false}}

---

## Remediation Cycle 2 — Architect

**Current overview:** The service will authenticate with Linear using an API key, fetch issues from a specified team using the Linear GraphQL API, and upsert them into a PostgreSQL database. It will expose HTTP endpoints for querying and syncing issues, and run a background sync loop at configurable intervals.

### Remediation Tasks

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| c1b27235 | backend | Create go.mod file | - | Produce artifact kind `code`, name `go.mod`. Create a go.mod file with module path github.com/bryanbarton525/linear-sync. Run go mod tidy to add dependencies. |
| ba9e99fd | backend | Remove duplicate NewLinearClient function | - | Produce artifact kind `code`, name `linear_client.go`. Remove the duplicate NewLinearClient function definition, keeping only one. |
| 48177a22 | backend | Complete PostgreSQLStorage struct initialization | - | Produce artifact kind `code`, name `postgres_storage.go`. Complete the struct initialization with proper closing brace and add missing import for sync.Mutex. |
| 775809ff | backend | Align field names between linear_client.go and postgres_storage.go Issue structs | - | Produce artifact kind `code`, name `linear_client.go` and `postgres_storage.go`. Align the field names between the two Issue structs (e.g., use StatusName in both). |
| ac55f4d1 | backend | Rename packages to match module structure | - | Produce artifact kind `code`, name `linear_client.go` and other related files. Rename the package names to 'internal/linear', 'internal/storage', 'internal/server' as per Go convention. |
| 01a1aa95 | backend | Properly wire HTTP server handlers | - | Produce artifact kind `code`, name `http_server.go`. Properly wire the HTTP server handlers to accept proper handlers. Ensure HandleIssues, handleSync methods are correctly set and wired. |

