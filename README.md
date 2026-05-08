# linear-sync

`linear-sync` is a small Go service that polls Linear issues for a team and upserts them into PostgreSQL.

## Requirements

- Go 1.21+
- PostgreSQL
- A Linear API key and team ID

## Configuration

Set the following environment variables before running:

- `LINEAR_API_KEY`: Linear API key used for GraphQL requests
- `LINEAR_TEAM_ID`: Linear team identifier to sync
- `DATABASE_URL`: PostgreSQL connection string

Example:

```bash
export LINEAR_API_KEY="lin_api_xxx"
export LINEAR_TEAM_ID="team-id"
export DATABASE_URL="postgres://user:pass@localhost:5432/linear?sslmode=disable"
```

## Database schema

The service expects this table:

```sql
CREATE TABLE IF NOT EXISTS linear_issues (
  id TEXT PRIMARY KEY,
  title TEXT,
  state TEXT,
  assignee TEXT
);
```

## Run

```bash
go run ./cmd/linear-sync
```

The service polls Linear every 5 minutes and upserts issues into `linear_issues` until it receives `SIGINT` or `SIGTERM`.
