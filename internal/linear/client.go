package linear

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Assignee struct {
	Name string
}

type Issue struct {
	ID       string
	Title    string
	State    string
	Assignee *Assignee
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) FetchIssues(ctx context.Context, teamID string) ([]Issue, error) {
	query := `{
		"query": "query($teamId: String!) { team(id: $teamId) { issues { nodes { id title state { name } assignee { name } } } } }",
		"variables": {"teamId": "` + teamID + `"}
	}`

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.linear.app/graphql", bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result struct {
		Data struct {
			Team struct {
				Issues struct {
					Nodes []struct {
						ID    string `json:"id"`
						Title string `json:"title"`
						State struct {
							Name string `json:"name"`
						} `json:"state"`
						Assignee *struct {
							Name string `json:"name"`
						} `json:"assignee"`
					} `json:"nodes"`
				} `json:"issues"`
			} `json:"team"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", result.Errors[0].Message)
	}

	nodes := result.Data.Team.Issues.Nodes
	issues := make([]Issue, 0, len(nodes))
	for _, n := range nodes {
		issue := Issue{
			ID:    n.ID,
			Title: n.Title,
			State: n.State.Name,
		}
		if n.Assignee != nil {
			issue.Assignee = &Assignee{Name: n.Assignee.Name}
		}
		issues = append(issues, issue)
	}
	return issues, nil
}
