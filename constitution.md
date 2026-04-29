# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

To create a reliable, scalable Go service that synchronizes Linear.app issues to PostgreSQL, ensuring data consistency, performance, and compliance with the toolchain validation profile.

## Goals

- Implement a Go service that polls the Linear GraphQL API every 5 minutes
- Upsert issues into a PostgreSQL table with columns: id, title, status, priority, updated_at
- Expose GET /issues and POST /sync and GET /healthz endpoints
- Use DATABASE_URL env var for PostgreSQL connection

## Constraints

- Use LINEAR_API_KEY environment variable for API authentication
- Ensure service runs with proper permissions for PostgreSQL access
- Handle API rate limiting and errors gracefully
- Validate against the toolchain validation profile (e.g., tests, build, formatting)

## Audience

Developers, system administrators, and operational engineers

## Output Medium

Go service with HTTP endpoints and PostgreSQL database integration

## Acceptance Criteria

- Service polls Linear GraphQL API every 5 minutes
- Upsert issues into PostgreSQL table with correct schema
- Exposes /issues GET endpoint with correct data
- Exposes /sync POST endpoint with proper validation
- Exposes /healthz GET endpoint with health check
- Uses DATABASE_URL env var for PostgreSQL connection
- All endpoints are accessible and functional
- Service passes all toolchain validation checks

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Implement Linear GraphQL API polling | Poll Linear GraphQL API every 5 minutes using LINEAR_API_KEY env var | Workflow requirements |
| F2 | must | Upsert issues into PostgreSQL table | Upsert issues into PostgreSQL table with columns: id, title, status, priority, updated_at | Workflow requirements |
| F3 | must | Expose /issues GET endpoint | Expose GET /issues endpoint with correct data | Workflow requirements |
| F4 | must | Expose /sync POST endpoint | Expose POST /sync endpoint with proper validation | Workflow requirements |
| F5 | must | Expose /healthz GET endpoint | Expose GET /healthz endpoint with health check | Workflow requirements |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Toolchain compliance | Service passes all toolchain validation checks (e.g., tests, build, formatting) | Workflow requirements |

## Dependencies

- github.com/jackc/dbc/blob/master/dbc.go
- github.com/graphql-go/graphql

---

## Constitution Amendment — Cycle 1

QA blocking issues classified as validation/environment failures (go_mod_tidy, go_build, go_test) and requirement gaps (service.go, issue-poller.go, http-handler.go, main.go). Remediation: create go.mod file, fix formatting, set up correct directory structure, implement missing code.
