// Package main implements a service that syncs Linear.app issues with PostgreSQL.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bryanbarton525/linear-sync/internal/poller"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

// IssuePoller is a struct that holds the necessary fields for polling and syncing issues.
type IssuePoller struct {
	pgPool   *pgx.Conn
	apiKey   string
	syncChan chan<- struct{}
}

// NewIssuePoller creates a new IssuePoller instance.
func NewIssuePoller(dbUrl, apiKey string, syncChan chan<- struct{}) (*IssuePoller, error) {
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	poller := &IssuePoller{
		pgPool:   conn,
		apiKey:   apiKey,
		syncChan: syncChan,
	}

	return poller, nil
}

// Start begins the polling loop.
func (p *IssuePoller) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.syncIssues(ctx)
		}
	}
}

// syncIssues polls the Linear API and upserts issues into PostgreSQL.
func (p *IssuePoller) syncIssues(ctx context.Context) {
	// Simulate a GraphQL API call to fetch issues.
	issues, err := p.fetchIssues(ctx)
	if err != nil {
		log.Printf("failed to fetch issues: %v", err)
		return
	}

	for _, issue := range issues {
		err = p.upsertIssue(ctx, issue)
		if err != nil {
			log.Printf("failed to upsert issue %s: %v", issue.ID, err)
			continue
		}
	}

	// Notify the main service that syncing is complete.
	p.syncChan <- struct{}{}
}

// fetchIssues simulates fetching issues from Linear GraphQL API.
func (p *IssuePoller) fetchIssues(ctx context.Context) ([]LinearIssue, error) {
	// This is a placeholder. Implement actual API call here.
	return []LinearIssue{}, nil
}

// upsertIssue inserts or updates an issue in PostgreSQL.
func (p *IssuePoller) upsertIssue(ctx context.Context, issue LinearIssue) error {
	_, err := p.pgPool.Exec(ctx,
		"INSERT INTO issues (id, title, status, priority, updated_at) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
		title = EXCLUDED.title,
		status = EXCLUDED.status,
		priority = EXCLUDED.priority,
		updated_at = EXCLUDED.updated_at",
		issue.ID, issue.Title, issue.Status, issue.Priority, issue.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to upsert issue: %w", err)
	}

	return nil
}

// LinearIssue represents an issue fetched from the Linear API.
type LinearIssue struct {
	ID        string `json:"id"
	Title     string `json:"title"
	Status    string `json:"status"
	Priority  int    `json:"priority"
	UpdatedAt int64  `json:"updatedAt"
}

// main initializes the service and starts polling.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	apiKey := os.Getenv("LINEAR_API_KEY")
	if dbUrl == "" || apiKey == "" {
		log.Fatal("DATABASE_URL and LINEAR_API_KEY must be set")
	}

	syncChan := make(chan struct{})
	poller, err := NewIssuePoller(dbUrl, apiKey, syncChan)
	if err != nil {
		log.Fatalf("failed to create poller: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go poller.Start(ctx)

	// Start HTTP server here.
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
