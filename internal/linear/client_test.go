package linear

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchIssues(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"team": {
					"issues": {
						"nodes": [
							{"id": "ISS-1", "title": "Fix bug", "state": {"name": "In Progress"}, "assignee": {"name": "Alice"}},
							{"id": "ISS-2", "title": "Add feature", "state": {"name": "Todo"}, "assignee": null}
						]
					}
				}
			}
		}`))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := NewClient("test-key")
	client.httpClient = &http.Client{}

	// override endpoint by using a transport that rewrites the URL
	client.httpClient.Transport = &rewriteTransport{base: ts.URL}

	issues, err := client.FetchIssues(context.Background(), "TEAM-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].ID != "ISS-1" {
		t.Errorf("expected ISS-1, got %s", issues[0].ID)
	}
	if issues[0].Assignee == nil || issues[0].Assignee.Name != "Alice" {
		t.Errorf("expected assignee Alice")
	}
	if issues[1].Assignee != nil {
		t.Errorf("expected nil assignee for ISS-2")
	}
}

func TestFetchIssuesHTTPError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := NewClient("test-key")
	client.httpClient = &http.Client{
		Transport: &rewriteTransport{base: ts.URL},
	}

	_, err := client.FetchIssues(context.Background(), "TEAM-1")
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestFetchIssuesGraphQLError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"errors": [{"message": "team not found"}]}`))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := NewClient("test-key")
	client.httpClient = &http.Client{
		Transport: &rewriteTransport{base: ts.URL},
	}

	_, err := client.FetchIssues(context.Background(), "TEAM-1")
	if err == nil {
		t.Fatal("expected error for GraphQL errors")
	}
}

// rewriteTransport redirects all requests to the test server.
type rewriteTransport struct {
	base string
}

func (r *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme = "http"
	req.URL.Host = r.base[len("http://"):]
	return http.DefaultTransport.RoundTrip(req)
}
