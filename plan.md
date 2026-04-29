# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

linear-sync is a Go service that polls the Linear GraphQL API every 5 minutes and upserts issues into a PostgreSQL table. It exposes /issues, /sync, and /healthz endpoints.

## Delivery Target

/var/lib/go-orca/workspaces/032ebf6d-08b9-401b-b283-44eab9daf752

## Tech Stack

- Go
- GraphQL
- PostgreSQL

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| LinearSyncService | Main struct for running the linear-sync service. |  | linear-sync |
| IssuePoller | Component responsible for polling the Linear GraphQL API and upserting issues into PostgreSQL. | LINEAR_API_KEY, DATABASE_URL | postgres_table_upserted |
| HTTPHandler | Handles HTTP requests for /issues, /sync, and /healthz endpoints. |  | http_responses |

## Architectural Decisions

1. **Use Go's standard library for HTTP client and server implementation.**
   - Rationale: Standard library is well-maintained, idiomatic, and sufficient for the required functionality.
   - Tradeoffs: No third-party dependencies, increased development time.
2. **Use a PostgreSQL driver from the Go ecosystem.**
   - Rationale: Popular and battle-tested, provides robust database connection management.
   - Tradeoffs: Learning curve, potential for licensing issues with certain drivers.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 6c81e540 | backend | Design the LinearSyncService struct | - | Produce artifact kind `go`, name `linear-sync/service.go`. Define the `LinearSyncService` struct with methods for running the service. This task does not involve upserting issues or handling HTTP requests, so it is a backend task. |
| 06a30604 | backend | Design the IssuePoller component | Design t | Produce artifact kind `go`, name `linear-sync/issue-poller.go`. Implement a `IssuePoller` struct with methods for polling the Linear GraphQL API and upserting issues into PostgreSQL. This task involves database operations, so it is a backend task. |
| 7549c818 | backend | Design the HTTPHandler component | Design t | Produce artifact kind `go`, name `linear-sync/http-handler.go`. Implement a `HTTPHandler` struct with methods for handling /issues, /sync, and /healthz endpoints. This task involves HTTP routing and response generation, so it is a backend task. |
| 01d9cb46 | backend | Initialize database connection in IssuePoller | Design t | Extend the `IssuePoller` artifact produced by the Design the IssuePoller component task. Implement a method to initialize the PostgreSQL database connection using the DATABASE_URL env var. |
| 22c5d2bd | backend | Implement issue polling and upserting in IssuePoller | Initiali | Extend the `IssuePoller` artifact produced by the Initialize database connection in IssuePoller task. Implement the logic to poll the Linear GraphQL API every 5 minutes and upsert issues into the PostgreSQL table. |
| 68fee1c3 | backend | Implement /issues GET endpoint in HTTPHandler | Design t | Extend the `HTTPHandler` artifact produced by the Design the HTTPHandler component task. Implement the /issues GET endpoint to return a list of issues. |
| 7fcae529 | backend | Implement /sync POST endpoint in HTTPHandler | Implemen | Extend the `HTTPHandler` artifact produced by the Implement /issues GET endpoint in HTTPHandler task. Implement the /sync POST endpoint to trigger an immediate sync of issues. |
| db951fca | backend | Implement /healthz GET endpoint in HTTPHandler | Implemen | Extend the `HTTPHandler` artifact produced by the Implement /sync POST endpoint in HTTPHandler task. Implement the /healthz GET endpoint to return service health status. |

---

## Remediation Cycle 1 — PM Triage

QA blocking issues classified as validation/environment failures (go_mod_tidy, go_build, go_test) and requirement gaps (service.go, issue-poller.go, http-handler.go, main.go). Remediation: create go.mod file, fix formatting, set up correct directory structure, implement missing code.

**QA blocking issues being triaged:**

- validation tidy_dependencies failed via go_mod_tidy: mcp: {"passed":false,"success":false,"stderr":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'\n","output":"go: go.mod file not found in current directory or any parent directory; see 'go help modules'","error":"exit status 1","metadata":{"command":"go mod tidy","duration_ms":31,"exit_code":1,"truncated":false}}
- validation format_code failed via go_fmt: mcp: {"passed":false,"success":false,"stderr":"HTTPHandler.go:44:35: string literal not terminated\nissue_poller.go:86:3: string literal not terminated\nissue_poller.go:87:6: missing ',' in argument list\nissue_poller.go:91:35: string literal not terminated\nissue_poller.go:93:2: missing ',' in argument list\nissue_poller.go:94:3: expected operand, found 'return'\nissue_poller.go:95:3: missing ',' before newline in argument list\nissue_poller.go:97:2: expected operand, found 'return'\nissue_poller.go:98:1: expected operand, found '}'\nissue_poller.go:101:1: missing ',' in argument list\nissue_poller.go:106:19: raw string literal not terminated\nissue_poller.go:139:3: missing ',' in argument list\n","output":"issue_poller.go:95:3: missing ',' before newline in argument list\nissue_poller.go:97:2: expected operand, found 'return'\nissue_poller.go:98:1: expected operand, found '}'\nissue_poller.go:101:1: missing ',' in argument list\nissue_poller.go:106:19: raw string literal not terminated\nissue_poller.go:139:3: missing ',' in argument list","error":"exit status 2","metadata":{"command":"gofmt -w .","duration_ms":163,"exit_code":2,"truncated":false}}
- validation run_tests failed via go_test: mcp: {"passed":false,"success":false,"stdout":"FAIL\t./... [setup failed]\nFAIL\n","stderr":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go test ./...","duration_ms":13,"exit_code":1,"truncated":false}}
- validation run_build failed via go_build: mcp: {"passed":false,"success":false,"stderr":"pattern ./...: directory prefix . does not contain main module or its selected dependencies\n","output":"pattern ./...: directory prefix . does not contain main module or its selected dependencies","error":"exit status 1","metadata":{"command":"go build ./...","duration_ms":8,"exit_code":1,"truncated":false}}
- [service.go] The provided service.go file does not contain any implementation details and appears to be a stub.: Implement the LinearSyncService struct with methods for running the service, including polling the Linear GraphQL API every 5 minutes and upserting issues into PostgreSQL.
- [linear-sync/issue-poller.go] The provided issue-poller.go file does not contain any implementation details and appears to be a stub.: Implement the IssuePoller component with methods for polling the Linear GraphQL API and upserting issues into PostgreSQL.
- [linear-sync/http-handler.go] The provided http-handler.go file does not contain any implementation details and appears to be a stub.: Implement the HTTPHandler component with methods for handling /issues, /sync, and /healthz endpoints.
- [main.go] The provided main.go file does not contain any implementation details and appears to be a stub.: Implement the /issues GET endpoint in HTTPHandler to return a list of issues from PostgreSQL, and ensure that all endpoints are accessible and functional.
- [HTTPHandler.go] The provided HTTPHandler.go file does not contain any implementation details and appears to be a stub.: Implement the /healthz GET endpoint in HTTPHandler to return service health status, and ensure that all endpoints are accessible and functional.

---

## Remediation Cycle 1 — Architect

**Current overview:** linear-sync is a Go service that polls the Linear GraphQL API every 5 minutes and upserts issues into a PostgreSQL table. It exposes /issues, /sync, and /healthz endpoints.

### Remediation Tasks

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 8e8dec14 | backend | Create go.mod file | - | Produce artifact kind `go`, name `go.mod`. Add necessary dependencies for Go modules, HTTP client, server, and PostgreSQL driver. |
| fa13c027 | backend | Fix formatting in code files | Create g | Produce artifact kind `go`, name `*.go`. Fix all formatting issues identified by go_fmt, ensuring all string literals are properly terminated and missing commas are added. |
| a8acc634 | backend | Initialize main.go with service initialization | Fix form | Produce artifact kind `go`, name `main.go`. Initialize LinearSyncService, IssuePoller, and HTTPHandler in the main function. |
| cbdca833 | backend | Implement LinearSyncService struct | Initiali | Produce artifact kind `go`, name `linear-sync/service.go`. Define the LinearSyncService struct with methods for running the service, including initialization of IssuePoller and HTTPHandler. |
| 59dc2ad8 | backend | Implement IssuePoller component | Initiali | Produce artifact kind `go`, name `linear-sync/issue-poller.go`. Implement the IssuePoller struct with methods for polling the Linear GraphQL API and upserting issues into PostgreSQL. |
| ccc791f8 | backend | Implement HTTPHandler component | Initiali | Produce artifact kind `go`, name `linear-sync/http-handler.go`. Implement the HTTPHandler struct with methods for handling /issues, /sync, and /healthz endpoints. |
| 5b63be3a | backend | Implement /issues GET endpoint in HTTPHandler | Initiali | Produce artifact kind `go`, name `linear-sync/http-handler.go`. Implement the /issues GET endpoint to return a list of issues from PostgreSQL. |
| 52589443 | backend | Implement /sync POST endpoint in HTTPHandler | Initiali | Produce artifact kind `go`, name `linear-sync/http-handler.go`. Implement the /sync POST endpoint to trigger an immediate sync of issues. |
| 4a3cd13c | backend | Implement /healthz GET endpoint in HTTPHandler | Initiali | Produce artifact kind `go`, name `linear-sync/http-handler.go`. Implement the /healthz GET endpoint to return service health status. |

