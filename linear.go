package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Issue represents a Linear issue.
type Issue struct {
	ID          string
	Title       string
	Description string
	State       string
	Priority    int
	Assignee    *Assignee
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Assignee represents a Linear user assigned to an issue.
type Assignee struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// rawState holds the decoded state object from the Linear API.
type rawState struct {
	Name string `json:"name"`
}

// rawNode holds one decoded issue node from the Linear API.
type rawNode struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       rawState  `json:"state"`
	Priority    int       `json:"priority"`
	Assignee    *Assignee `json:"assignee"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// rawIssues holds the decoded issues list from the Linear API.
type rawIssues struct {
	Nodes []rawNode `json:"nodes"`
}

// rawTeam holds the decoded team object from the Linear API.
type rawTeam struct {
	Issues rawIssues `json:"issues"`
}

// rawData holds the decoded data field from the Linear API.
type rawData struct {
	Team rawTeam `json:"team"`
}

// rawError holds a single error from the Linear API response.
type rawError struct {
	Message string `json:"message"`
}

// rawResponse is the top-level decoded Linear API response.
type rawResponse struct {
	Data   rawData    `json:"data"`
	Errors []rawError `json:"errors"`
}

type linearClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func newClient(apiKey string) *linearClient {
	return &linearClient{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		baseURL:    "https://api.linear.app/graphql",
	}
}

func (c *linearClient) fetchIssues(ctx context.Context, teamID string) ([]Issue, error) {
	query := fmt.Sprintf(`{"query":"query { team(id: \"%s\") { issues { nodes { id title description state { name } priority assignee { id name email } createdAt updatedAt } } } }"}`, teamID)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result rawResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, errors.New(result.Errors[0].Message)
	}

	nodes := result.Data.Team.Issues.Nodes
	issues := make([]Issue, len(nodes))
	for i, raw := range nodes {
		issues[i] = Issue{
			ID:          raw.ID,
			Title:       raw.Title,
			Description: raw.Description,
			State:       raw.State.Name,
			Priority:    raw.Priority,
			Assignee:    raw.Assignee,
			CreatedAt:   raw.CreatedAt,
			UpdatedAt:   raw.UpdatedAt,
		}
	}
	return issues, nil
}
