# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

A Go service to synchronize Linear.app issues with PostgreSQL, providing HTTP endpoints for querying and syncing, and a background sync loop for real-time updates.

## Goals

- Implement a secure authentication mechanism using Linear API key
- Fetch and map Linear GraphQL issues to PostgreSQL columns
- Expose HTTP endpoints for querying and syncing issues
- Run a background sync loop every SYNC_INTERVAL_MINUTES
- Ensure compliance with the specified Go toolchain and libraries

## Constraints

- Use environment variables LINEAR_API_KEY and LINEAR_TEAM_KEY for authentication and team selection
- Use DATABASE_URL for PostgreSQL connection string
- Limit sync intervals to SYNC_INTERVAL_MINUTES (default 15 minutes)
- Comply with Go standard library and pgx/v5 for PostgreSQL operations
- Include a GET /healthz endpoint for service health checks

## Audience

Developers, system administrators, and workflow engineers

## Output Medium

HTTP

## Acceptance Criteria

- Service authenticates using LINEAR_API_KEY and LINEAR_TEAM_KEY
- Fetches all issues from LINEAR_TEAM_KEY using Linear GraphQL API
- Upserts issues into PostgreSQL table with matching columns
- Exposes GET /issues HTTP endpoint returning JSON data
- Exposes POST /sync endpoint to trigger full re-sync
- Background sync loop runs every SYNC_INTERVAL_MINUTES
- Includes GET /healthz endpoint

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Authentication with Linear API key | Service must authenticate using LINEAR_API_KEY environment variable | user request |
| F2 | must | Fetch Linear GraphQL issues | Service must fetch all issues from LINEAR_TEAM_KEY using Linear GraphQL API | user request |
| F3 | must | Upsert issues into PostgreSQL | Service must upsert issues into PostgreSQL table with matching columns | user request |
| F4 | must | HTTP endpoints for querying and syncing | Service must expose GET /issues and POST /sync endpoints | user request |
| F5 | must | Background sync loop | Service must run a background sync loop every SYNC_INTERVAL_MINUTES | user request |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Compliance with Go toolchain | Service must use standard library net/http, pgx/v5 for PostgreSQL, and github.com/bryanbarton525/linear-sync module | user request |

## Dependencies

- github.com/bryanbarton525/linear-sync
- github.com/jackc/pgx/v5
- net/http

---

## Constitution Amendment — Cycle 1

QA blocking issues: 1. Missing go.mod file (requirement gap), 2. Duplicate function (implementation defect), 3. Incomplete struct initialization (implementation defect), 4. Field name inconsistency (implementation defect), 5. Inconsistent package names (design gap), 6. Server handlers not properly initialized (implementation defect).
