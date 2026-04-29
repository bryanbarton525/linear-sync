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

