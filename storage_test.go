package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Upsert(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/linear_sync_test?sslmode=disable"
	}

	db, err := newDB(connStr)
	if err != nil {
		t.Skip("PostgreSQL not available in test environment:", err)
	}
	defer db.Close()

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
	defer db.Exec("DROP TABLE IF EXISTS linear_issues")

	store := newStorage(db)
	ctx := context.Background()
	now := time.Now()

	issues := []Issue{
		{
			ID:          "test-issue-1",
			Title:       "Test Issue",
			Description: "Test Description",
			State:       "In Progress",
			Priority:    1,
			Assignee: &Assignee{
				ID:    "user-1",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	err = store.upsert(ctx, issues)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM linear_issues WHERE id = $1", "test-issue-1").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	issues[0].Title = "Updated Title"
	err = store.upsert(ctx, issues)
	assert.NoError(t, err)

	var title string
	err = db.QueryRow("SELECT title FROM linear_issues WHERE id = $1", "test-issue-1").Scan(&title)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", title)
}

func TestStorage_UpsertEmpty(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/linear_sync_test?sslmode=disable"
	}

	db, err := newDB(connStr)
	if err != nil {
		t.Skip("PostgreSQL not available in test environment:", err)
	}
	defer db.Close()

	store := newStorage(db)
	err = store.upsert(context.Background(), []Issue{})
	assert.NoError(t, err)
}
