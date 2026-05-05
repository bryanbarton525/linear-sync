package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

// Storage handles persistence of Linear issues to PostgreSQL.
type Storage struct {
	db *sql.DB
}

// newDB opens and pings a PostgreSQL database connection.
func newDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

// newStorage creates a new Storage with the given database connection.
func newStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// upsert inserts or updates the given issues in the database.
func (s *Storage) upsert(ctx context.Context, issues []Issue) error {
	if len(issues) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO linear_issues (id, title, description, state, priority, assignee, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			state = EXCLUDED.state,
			priority = EXCLUDED.priority,
			assignee = EXCLUDED.assignee,
			updated_at = EXCLUDED.updated_at
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, issue := range issues {
		assigneeJSON, err := json.Marshal(issue.Assignee)
		if err != nil {
			return fmt.Errorf("failed to marshal assignee for issue %s: %w", issue.ID, err)
		}
		if _, err = stmt.ExecContext(ctx, issue.ID, issue.Title, issue.Description,
			issue.State, issue.Priority, assigneeJSON, issue.CreatedAt, issue.UpdatedAt); err != nil {
			return fmt.Errorf("failed to upsert issue %s: %w", issue.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
