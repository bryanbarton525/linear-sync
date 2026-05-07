package storage

import (
	"context"
	"database/sql"
	"fmt"

	"linear-sync/internal/linear"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Upsert(ctx context.Context, issues []linear.Issue) error {
	for _, issue := range issues {
		assigneeName := ""
		if issue.Assignee != nil {
			assigneeName = issue.Assignee.Name
		}
		_, err := s.db.ExecContext(ctx,
			`INSERT INTO linear_issues (id, title, state, assignee)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (id) DO UPDATE
			 SET title = EXCLUDED.title,
			     state = EXCLUDED.state,
			     assignee = EXCLUDED.assignee`,
			issue.ID, issue.Title, issue.State, assigneeName,
		)
		if err != nil {
			return fmt.Errorf("upsert issue %s: %w", issue.ID, err)
		}
	}
	return nil
}
