package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchIssues(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	resp := map[string]interface{}{
		"data": map[string]interface{}{
			"team": map[string]interface{}{
				"issues": map[string]interface{}{
					"nodes": []interface{}{
						map[string]interface{}{
							"id":          "LIN-1",
							"title":       "Test issue",
							"description": "desc",
							"state":       map[string]interface{}{"name": "Todo"},
							"priority":    2,
							"assignee":    nil,
							"createdAt":   now.Format(time.RFC3339),
							"updatedAt":   now.Format(time.RFC3339),
						},
					},
				},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := &linearClient{
		apiKey:     "test-key",
		httpClient: ts.Client(),
		baseURL:    ts.URL,
	}

	issues, err := c.fetchIssues(context.Background(), "team-1")
	require.NoError(t, err)
	require.Len(t, issues, 1)
	assert.Equal(t, "LIN-1", issues[0].ID)
	assert.Equal(t, "Test issue", issues[0].Title)
	assert.Equal(t, "Todo", issues[0].State)
	assert.Equal(t, 2, issues[0].Priority)
}

func TestFetchIssues_APIError(t *testing.T) {
	resp := map[string]interface{}{
		"errors": []interface{}{
			map[string]interface{}{"message": "not found"},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := &linearClient{
		apiKey:     "test-key",
		httpClient: ts.Client(),
		baseURL:    ts.URL,
	}

	_, err := c.fetchIssues(context.Background(), "team-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
