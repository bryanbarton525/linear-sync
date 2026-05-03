package linear

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchIssues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"team": {
					"issues": {
						"nodes": [
							{
								"id": "issue-1",
								"title": "Test Issue",
								"description": "Test Description",
								"state": {"name": "In Progress"},
								"priority": 2,
								"assignee": {
									"id": "user-1",
									"name": "John Doe",
									"email": "john@example.com"
								},
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-02T00:00:00Z"
							}
						]
					}
				}
			}
		}`))
	}))
	defer ts.Close()

	client := NewClient("test-api-key", WithBaseURL(ts.URL))
	issues, err := client.FetchIssues(context.Background(), "test-team-id")

	assert.NoError(t, err)
	assert.Len(t, issues, 1)
	assert.Equal(t, "issue-1", issues[0].ID)
	assert.Equal(t, "Test Issue", issues[0].Title)
	assert.Equal(t, "In Progress", issues[0].State)
	assert.Equal(t, 2, issues[0].Priority)
	assert.NotNil(t, issues[0].Assignee)
	assert.Equal(t, "user-1", issues[0].Assignee.ID)
	assert.Equal(t, "John Doe", issues[0].Assignee.Name)
	assert.Equal(t, "john@example.com", issues[0].Assignee.Email)
}

func TestFetchIssuesError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer ts.Close()

	client := NewClient("test-api-key", WithBaseURL(ts.URL))
	issues, err := client.FetchIssues(context.Background(), "test-team-id")

	assert.Error(t, err)
	assert.Nil(t, issues)
}
