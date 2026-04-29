# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

To create a robust, scalable Go service that synchronizes Linear.app issues to PostgreSQL using the Linear GraphQL API, ensuring secure authentication, efficient data retrieval, and reliable endpoint exposure.

## Goals

- Implement authentication with Linear API key
- Fetch all issues from a configurable team
- Upsert issues into PostgreSQL with matching fields
- Expose GET /issues HTTP endpoint for querying
- Expose POST /sync HTTP endpoint for re-sync
- Run background sync loop every SYNC_INTERVAL_MINUTES
- Implement GET /healthz endpoint

## Constraints

- Use standard library net/http, pgx/v5 for Postgres, and github.com/bryanbarton525/linear-sync module
- Environment variables LINEAR_API_KEY, LINEAR_TEAM_KEY, DATABASE_URL, SYNC_INTERVAL_MINUTES must be set
- No external dependencies beyond specified modules

## Audience

Developers, system administrators, and operational engineers

## Output Medium

HTTP server (net/http), PostgreSQL (pgx/v5)

## Acceptance Criteria

- Service starts with LINEAR_API_KEY and LINEAR_TEAM_KEY set
- POST /sync triggers full sync with Linear API
- GET /issues returns JSON with all specified fields
- Background sync loop runs every SYNC_INTERVAL_MINUTES
- GET /healthz returns 200 status code
- Error handling for missing environment variables

## Out of Scope

- Third-party dependencies beyond specified modules
- Advanced features not requested

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Authentication with Linear API key | Service must authenticate with Linear using the environment variable LINEAR_API_KEY | user request |
| F2 | must | Fetch issues from Linear team | Service must fetch all issues from the environment variable LINEAR_TEAM_KEY using the Linear GraphQL API | user request |
| F3 | must | Upsert issues into PostgreSQL | Service must upsert issues into the PostgreSQL table defined by DATABASE_URL | user request |
| F4 | must | Expose GET /issues endpoint | Service must expose a GET /issues endpoint that queries the local PostgreSQL table and returns JSON | user request |
| F5 | must | Expose POST /sync endpoint | Service must expose a POST /sync endpoint that triggers a full re-sync from Linear | user request |
| F6 | must | Background sync loop | Service must run a background sync loop every SYNC_INTERVAL_MINUTES | user request |
| F7 | must | Expose GET /healthz endpoint | Service must expose a GET /healthz endpoint | user request |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Error handling | Service must handle errors gracefully, including missing environment variables and API rate limits | user request |
| NF2 | should | Performance | Service must efficiently handle upserts and query operations with low latency | user request |
| NF3 | should | Consistency | Service must ensure data consistency between Linear and PostgreSQL tables | user request |

## Dependencies

- github.com/bryanbarton525/linear-sync
- github.com/jackc/pgx/v4
- net/http

