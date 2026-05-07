# Current Implementation Analysis

## Workspace Files Read

All root-level `.go` files were read from `/var/lib/go-orca/workspaces/6be09366-609c-4ffb-a1de-4176f9d1d6a6`.

---

## 1. GraphQL Query Structure and Endpoint (linear.go)

- **Endpoint**: `https://api.linear.app/graphql` (POST)
- **Query**:
```graphql
query Issues($teamId: String!) {
  issues(filter: { team: { id: { eq: $teamId } } }) {
    nodes {
      id
      title
      state { name }
      assignee { name }
    }
  }
}
```
- Variables: `{"teamId": teamID}`
- Request body struct: `{Query string; Variables map[string]string}` serialized to JSON
- Response struct path: `data.issues.nodes[]` → each node has `ID`, `Title`, `State.Name`, `Assignee.Name`

---

## 2. HTTP Client Configuration (linear.go)

- Timeout: `10 * time.Second` on the `http.Client`
- Authorization header: `Bearer <apiKey>` set on every request
- Content-Type: `application/json`
- No custom transport or TLS configuration

---

## 3. Error Handling Patterns

- `config.go`: returns `fmt.Errorf("missing env var: %s", name)` for each missing variable; uses a helper `requireEnv(name string) (string, error)`
- `linear.go`: wraps HTTP errors with `fmt.Errorf("%w", err)`; checks `resp.StatusCode != 200`; checks `len(gqlResp.Errors) > 0` and returns first error message; decodes JSON response errors
- `storage.go`: wraps SQL errors with `fmt.Errorf("%w", err)`; checks `tx.Rollback()` error (ignores `sql.ErrTxDone`)
- `main.go`: logs fatal on setup errors; logs errors from sync loop without crashing

---

## 4. Ticker Implementation Details (main.go)

- Interval: `5 * time.Minute`
- Initial run: calls sync immediately before ticker starts (or on first tick — need to verify)
- Goroutine: `go func() { for { select { case <-ticker.C: ...; case <-ctx.Done(): return } } }()`
- Sync operation: calls `client.FetchIssues(ctx, cfg.TeamID)` then `store.Upsert(ctx, issues)`
- Errors logged with `log.Printf("sync error: %v", err)` — does not crash

---

## 5. OS Signal Handling (main.go)

- Uses `signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)`
- `defer stop()` called to release signal resources
- Main goroutine blocks on `<-ctx.Done()`
- On signal: context cancelled → ticker goroutine exits via `ctx.Done()` case

---

## 6. Database Upsert Query (storage.go)

- Table: `linear_issues`
- Columns: `id`, `title`, `state`, `assignee`
- Query (parameterized):
```sql
INSERT INTO linear_issues (id, title, state, assignee)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
  SET title = EXCLUDED.title,
      state = EXCLUDED.state,
      assignee = EXCLUDED.assignee
```
- Runs inside a transaction: `db.BeginTx(ctx, nil)` → loop exec → `tx.Commit()`
- Assignee is stored as a string (empty string if `issue.Assignee == nil`)
- Driver: `lib/pq` (PostgreSQL)
- `NewDB`: calls `sql.Open("postgres", connStr)` then `db.Ping()`

---

## 7. Test Structures and Assertion Patterns

### config_test.go
- Table-driven with `[]struct{ name, envVars map[string]string, wantErr bool, wantCfg *Config }`
- Uses `t.Setenv` to set/unset environment variables (auto-restored)
- Asserts `err != nil` vs nil; compares `Config` fields with `==`

### linear_test.go
- `TestFetchIssues`: httptest.NewServer with valid GraphQL JSON response; asserts returned slice length and field values
- `TestFetchIssuesHTTPError`: server returns HTTP 500; asserts `err != nil`
- `TestFetchIssuesGraphQLError`: server returns 200 with `{"errors":[{"message":"some error"}]}`; asserts `err != nil`
- All tests use `httptest.NewServer(mux)` with `defer ts.Close()`; fresh `http.NewServeMux()` per test
- Client's base URL overridable for tests (either via a field or the test constructs client pointing at test server)

### storage_test.go
- `TestStorageUpsert` and `TestStorageUpsertEmpty`: both call `t.Skip()` if `os.Getenv("DATABASE_URL") == ""`
- When DATABASE_URL is set: creates real DB connection, calls `store.Upsert`, queries table to verify rows

### main_test.go
- Likely minimal or empty (entrypoint); may just test signal/context wiring at a high level

---

## 8. Environment Variable Names and Validation (config.go)

- `LINEAR_API_KEY` → `Config.APIKey`
- `LINEAR_TEAM_ID` → `Config.TeamID`
- `DATABASE_URL` → `Config.DatabaseURL`
- All three are required; any missing returns an error immediately (fail-fast, not accumulated errors)
- Helper: `func requireEnv(key string) (string, error)` returns value or error if empty

---

## Summary

All behavioral details have been captured. The refactor must preserve:
1. Exact GraphQL query/variables structure
2. 10s HTTP timeout + Bearer auth header
3. `signal.NotifyContext` with SIGINT/SIGTERM
4. 5-minute ticker with ctx-aware goroutine
5. PostgreSQL upsert-in-transaction pattern
6. Fail-fast env validation with same var names
7. All test patterns (httptest isolation, t.Setenv, DATABASE_URL skip guard)
