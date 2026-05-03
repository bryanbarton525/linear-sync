# Constitution

> Authored by the Project Manager. Immutable for the rest of the workflow; remediation triage may append a `Constitution Amendment` section but cannot rewrite this file.

## Vision

To provide a robust, scalable, and maintainable Go service that synchronizes Linear.app issues to PostgreSQL using the Linear GraphQL API, ensuring reliable and efficient data synchronization.

## Goals

- Implement a service that polls Linear.app issues at 5-minute intervals
- Integrate with the Linear GraphQL API to fetch and store issues
- Ensure data consistency and correctness through proper error handling and retries
- Support configuration via environment variables and a config file

## Constraints

- Use the Linear GraphQL API for data retrieval
- Poll Linear.app issues every 5 minutes
- Handle potential API rate limits and errors gracefully
- Ensure data integrity through proper transaction management
- Do not modify existing file structures or add unused imports

## Audience

Developers, sysadmins, and CI/CD pipeline operators

## Output Medium

A Go service with built-in functionality and test coverage

## Acceptance Criteria

- The service builds and runs successfully with all tests passing
- The service polls Linear.app issues every 5 minutes
- The service handles API errors and rate limits appropriately
- The service logs critical errors and warnings
- The service supports configuration via environment variables and a config file

## Out of Scope

- Validation and logging enhancements
- Additional error handling scenarios
- Performance optimizations

## Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| F1 | must | Implement Linear GraphQL API integration | Integrate with the Linear GraphQL API to fetch and store issues | code |
| F2 | must | Poll Linear.app issues every 5 minutes | Poll Linear.app issues at a fixed interval of 5 minutes | code |
| F3 | must | Handle API errors and rate limits | Implement error handling for API rate limits and potential failures | code |
| F4 | must | Ensure data integrity | Maintain data consistency and correctness through proper transaction management | code |
| F5 | must | Support configuration via environment variables | Allow configuration via environment variables for flexibility | code |

## Non-Functional Requirements

| ID | Priority | Title | Description | Source |
|---|---|---|---|---|
| NF1 | must | Error handling and logging | Implement robust error handling and logging for critical issues | code |
| NF2 | must | Test coverage | Ensure comprehensive test coverage for all functionality | code |

## Dependencies

- internal/config/config.go
- internal/linear/client.go
- internal/storage/storage.go

