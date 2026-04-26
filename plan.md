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

