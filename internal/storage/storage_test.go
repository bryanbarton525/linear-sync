package storage

import (
	"context"
	"os"
	"testing"

	"linear-sync/internal/linear"
)

func TestStorageUpsert(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		t.Skip("DATABASE_URL not set")
	}

	db, err := NewDB(connStr)
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS linear_issues (
		id TEXT PRIMARY KEY,
		title TEXT,
		state TEXT,
		assignee TEXT
	)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	defer db.Exec(`DELETE FROM linear_issues WHERE id IN ('test-1', 'test-2')`)

	store := New(db)
	issues := []linear.Issue{
		{ID: "test-1", Title: "First Issue", State: "Todo", Assignee: &linear.Assignee{Name: "Alice"}},
		{ID: "test-2", Title: "Second Issue", State: "In Progress", Assignee: nil},
	}

	if err := store.Upsert(context.Background(), issues); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM linear_issues WHERE id IN ('test-1', 'test-2')`).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 rows, got %d", count)
	}

	updated := []linear.Issue{
		{ID: "test-1", Title: "Updated Issue", State: "Done", Assignee: &linear.Assignee{Name: "Bob"}},
	}
	if err := store.Upsert(context.Background(), updated); err != nil {
		t.Fatalf("Upsert update: %v", err)
	}

	var title string
	if err := db.QueryRow(`SELECT title FROM linear_issues WHERE id = 'test-1'`).Scan(&title); err != nil {
		t.Fatalf("select: %v", err)
	}
	if title != "Updated Issue" {
		t.Errorf("expected 'Updated Issue', got %q", title)
	}
}

func TestStorageUpsertEmpty(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		t.Skip("DATABASE_URL not set")
	}

	db, err := NewDB(connStr)
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	defer db.Close()

	store := New(db)
	if err := store.Upsert(context.Background(), []linear.Issue{}); err != nil {
		t.Errorf("Upsert empty: %v", err)
	}
}
