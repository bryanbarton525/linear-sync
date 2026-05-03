package storage

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"linear-sync/internal/linear"
)

func TestStorage_Upsert(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/linear_sync_test?sslmode=disable"
	}

	db, err := NewDB(connStr)
	if err != nil {
		t.Skip("PostgreSQL not available in test environment:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("PostgreSQL not available in test environment:", err)
	}

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS linear_issues (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			state TEXT NOT NULL,
			priority INTEGER,
			assignee JSONB NOT NULL DEFAULT '{}',
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	// Cleanup
	defer func() {
		_, _ = db.Exec("DROP TABLE IF EXISTS linear_issues")
	}()

	store := New(db)
	ctx := context.Background()

	now := time.Now()
	issues := []linear.Issue{
		{
			ID:          "test-issue-1",
			Title:       "Test Issue",
			Description: "Test Description",
			State:       "In Progress",
			Priority:    1,
			Assignee: &linear.Assignee{
				ID:    "user-1",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	err = store.Upsert(ctx, issues)
	assert.NoError(t, err)

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM linear_issues WHERE id = $1", "test-issue-1").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Test upsert (update)
	issues[0].Title = "Updated Title"
	err = store.Upsert(ctx, issues)
	assert.NoError(t, err)

	var title string
	err = db.QueryRow("SELECT title FROM linear_issues WHERE id = $1", "test-issue-1").Scan(&title)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", title)
}
