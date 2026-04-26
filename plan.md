# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This design outlines the implementation of a Go service that syncs Linear.app issues to PostgreSQL via the Linear GraphQL API. The service will authenticate with Linear using an API key from an environment variable, fetch issues from a configurable team, upsert them into a PostgreSQL table, and expose HTTP endpoints for querying issues, triggering a full sync, and checking health.

## Delivery Target

git@github.com:bryanbarton525/linear-sync.git#workflow/711a06f2-e678-4d96-92f2-4932e43e6a30

## Tech Stack

- Go
- pgx/v5
- net/http

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| LinearClient | Handles authentication and fetching of Linear.app issues using the GraphQL API. Manages pagination to fetch all issues from the specified team. | LINEAR_API_KEY, LINEAR_TEAM_KEY | issues |
| PostgresClient | Handles connecting to and interacting with PostgreSQL database. Upserts issues into a table defined by DATABASE_URL. | DATABASE_URL | upsert_result |
| SyncService | Orchestrates the synchronization of issues between Linear.app and PostgreSQL. Manages background sync loops and HTTP endpoints. | issues, upsert_result |  |

## Architectural Decisions

1. **Use pgx/v5 for PostgreSQL interactions**
   - Rationale: pgx/v5 is a well-maintained, idiomatic Go library for PostgreSQL that supports context-aware operations.
   - Tradeoffs: None significant; pgx/v5 offers good performance and flexibility.
2. **Use standard library net/http for HTTP server**
   - Rationale: net/http is the standard Go library for building HTTP servers, providing a robust foundation for our endpoints.
   - Tradeoffs: None significant; net/http is widely used and well-documented.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 09b4a902 | backend | Define LinearClient Interface | Initiali | Produce artifact kind `code`, name `internal/linear/client.go`. Define the LinearClient interface with methods for authentication and fetching issues from Linear.app using the GraphQL API. Ensure it handles pagination. |
| 5e28a3ae | backend | Implement LinearClient | Define L | Produce artifact kind `code`, name `internal/linear/client.go`. Implement the LinearClient interface to authenticate using LINEAR_API_KEY and fetch issues from the team specified by LINEAR_TEAM_KEY. Include error handling for network failures and invalid responses. |
| e7a3e234 | backend | Define PostgresClient Interface | Implemen | Produce artifact kind `code`, name `internal/postgres/client.go`. Define the PostgresClient interface with methods for connecting to PostgreSQL and upserting issues into a specified table. |
| 411aeea6 | backend | Implement PostgresClient | Define P | Produce artifact kind `code`, name `internal/postgres/client.go`. Implement the PostgresClient interface to connect using DATABASE_URL and upsert issues into a PostgreSQL table. Ensure it handles database errors gracefully. |
| f6d6940c | backend | Define SyncService Interface | Implemen | Produce artifact kind `code`, name `internal/sync/service.go`. Define the SyncService interface with methods for orchestrating background sync loops and handling HTTP endpoints. |
| bed474a7 | backend | Implement SyncService | Define S | Produce artifact kind `code`, name `internal/sync/service.go`. Implement the SyncService interface to manage background sync loops with context-aware cancellation and expose HTTP endpoints for issues, syncing, and health checks. |
| 2736d625 | backend | Setup Main Function | Implemen | Produce artifact kind `code`, name `main.go`. Setup the main function to initialize and run the HTTP server, integrating LinearClient, PostgresClient, and SyncService. Ensure all components are initialized with environment variables. |
| 0231195e | backend | Add Tests for LinearClient | Implemen | Produce artifact kind `code`, name `_test.go`. Add table-driven tests covering authentication and issue fetching using httptest.NewServer() with defer ts.Close(). Ensure all test artifacts use separate packages. |
| e8cde1eb | backend | Add Tests for PostgresClient | Implemen | Produce artifact kind `code`, name `_test.go`. Add table-driven tests covering upsert operations using httptest.NewServer() with defer ts.Close(). Ensure all test artifacts use separate packages. |
| b3b1289f | backend | Add Tests for SyncService | Implemen | Produce artifact kind `code`, name `_test.go`. Add table-driven tests covering background sync and HTTP endpoints using httptest.NewServer() with defer ts.Close(). Ensure all test artifacts use separate packages. |
| 7d829f6b | ops | Initialize Git Repository | - | Produce artifact kind `git_commit`, name `initial_repo_setup.git`. Initialize the git repository in the workspace path `/var/lib/go-orca/workspaces/711a06f2-e678-4d96-92f2-4932e43e6a30` with a README.md file. Ensure all subsequent tasks are tracked by this git repository. |
| b266ef70 | ops | Commit Initial Codebase | Setup Ma, Add Test, Add Test, Add Test | Produce artifact kind `git_commit`, name `initial_codebase_commit.git`. Commit all implemented code, interfaces, and tests to the git repository. Ensure the commit message is descriptive. |

