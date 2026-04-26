# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

To create a reliable, scalable Go service that synchronizes Linear.app issues to PostgreSQL using the Linear GraphQL API, ensuring data consistency, efficient operations, and robust error handling.

## Goals

- Implement authentication with Linear using an API key from an environment variable.
- Fetch all issues from a configurable team using the Linear GraphQL API.
- Upsert issues into a PostgreSQL table with matching columns.
- Expose a GET /issues HTTP endpoint to query the local PostgreSQL table.
- Expose a POST /sync HTTP endpoint to trigger a full re-sync from Linear.
- Run a background sync loop every SYNC_INTERVAL_MINUTES.
- Include a GET /healthz endpoint for system health checks.

## Constraints

- Use standard library net/http, pgx/v5 for Postgres, and the github.com/bryanbarton525/linear-sync module.
- Environment variables LINEAR_API_KEY, LINEAR_TEAM_KEY, DATABASE_URL, and SYNC_INTERVAL_MINUTES must be set.
- The service must handle errors gracefully and return appropriate HTTP status codes.
- The background sync loop must be non-blocking and use context.Context for cancellation.

## Audience

Developers, system administrators, and operational engineers.

## Output Medium

HTTP and PostgreSQL database.

## Acceptance Criteria

- The service authenticates with Linear using the provided API key.
- The service fetches all issues from the specified team and includes the required fields.
- The service upserts issues into the PostgreSQL table correctly.
- The GET /issues endpoint returns the issues as JSON.
- The POST /sync endpoint triggers a full re-sync from Linear.
- The background sync loop runs periodically with the specified interval.
- The GET /healthz endpoint returns a 200 status code.

## Out of Scope

- Third-party dependencies beyond the specified module.
- Custom database schemas or advanced query capabilities.

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Authentication with Linear API Key | The service must read LINEAR_API_KEY from the environment and authenticate with Linear using the GraphQL API. | environment variable |
| F2 | must | Fetch Issues from Linear Team | The service must fetch all issues from the specified LINEAR_TEAM_KEY team using the Linear GraphQL API, including the required fields. | environment variable |
| F3 | must | Upsert Issues into PostgreSQL | The service must upsert issues into the PostgreSQL table defined by DATABASE_URL, with columns matching the fetched fields. | environment variable |
| F4 | must | Expose GET /issues Endpoint | The service must expose a GET /issues endpoint that queries the local PostgreSQL table and returns the issues as JSON. | service |
| F5 | must | Expose POST /sync Endpoint | The service must expose a POST /sync endpoint that triggers a full re-sync from Linear on demand. | service |
| F6 | must | Background Sync Loop | The service must run a background sync loop every SYNC_INTERVAL_MINUTES. | service |
| F7 | must | Expose GET /healthz Endpoint | The service must expose a GET /healthz endpoint for system health checks. | service |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Error Handling | The service must handle errors gracefully, including network failures, invalid API responses, and database errors. | service |
| NF2 | must | Performance | The service must efficiently handle the background sync loop and ensure minimal latency for the GET /issues endpoint. | service |
| NF3 | must | Security | The service must securely store and use API keys, and avoid exposing sensitive data in logs or error messages. | service |

## Dependencies

- github.com/bryanbarton525/linear-sync
- github.com/jackc/pgx/v5
- net/http

