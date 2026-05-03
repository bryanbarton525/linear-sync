package linear

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

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

type Assignee struct {
	ID    string
	Name  string
	Email string
}

type rawIssue struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       struct {
		Name string `json:"name"`
	} `json:"state"`
	Priority  int       `json:"priority"`
	Assignee  *Assignee `json:"assignee"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

type ClientOption func(*Client)

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		baseURL:    "https://api.linear.app/graphql",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

func (c *Client) FetchIssues(ctx context.Context, teamID string) ([]Issue, error) {
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

	var result struct {
		Data struct {
			Team struct {
				Issues struct {
					Nodes []rawIssue `json:"nodes"`
				} `json:"issues"`
			} `json:"team"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, errors.New(result.Errors[0].Message)
	}

	issues := make([]Issue, len(result.Data.Team.Issues.Nodes))
	for i, raw := range result.Data.Team.Issues.Nodes {
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
