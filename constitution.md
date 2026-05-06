# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

Build a Go service that syncs Linear.app issues to PostgreSQL, ensuring robust error handling, test coverage, and production-ready functionality.

## Goals

- Implement a service to sync Linear.app issues to PostgreSQL
- Ensure seamless integration with Linear API
- Pass all automated tests and validation checks
- Maintain data consistency and integrity through upsert operations

## Constraints

- Use provided Go modules (github.com/lib/pq, github.com/stretchr/testify)
- Ensure service runs in production environment
- Handle all error conditions gracefully
- Maintain test isolation and consistent test artifacts

## Audience

Developers, testers, and operational personnel

## Output Medium

Go service with PostgreSQL database integration

## Acceptance Criteria

- Service starts and runs without errors
- Issues are synced from Linear API to PostgreSQL
- All tests pass with 4 PASS and 2 SKIP
- Service commits to main branch with proper changelog

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Implement Linear API integration | Fetch issues from Linear API using GraphQL | linear.go |
| F2 | must | Data persistence | Persist issues to PostgreSQL with upsert operations | storage.go |
| F3 | must | Error handling | Handle API errors and database errors gracefully | linear.go |
| F4 | must | Test coverage | Implement test suite covering all edge cases | linear_test.go |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Performance | Sync issues at 500 per minute with 99th percentile latency < 200ms | constitution |
| NF2 | must | Test isolation | No shared state between tests | constitution |
| NF3 | must | Validation | All test artifacts use httptest.NewServer() with defer ts.Close() | constitution |

## Dependencies

- github.com/lib/pq v1.10.9
- github.com/stretchr/testify v1.9.0
- github.com/davecgh/go-spew v1.1.1
- github.com/pmezard/go-difflib v1.0.0
- gopkg.in/yaml.v3 v3.0.1

