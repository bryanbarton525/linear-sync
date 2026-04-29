package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Issue struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Status    string       `json:"status"`
	Priority  int          `json:"priority"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type HTTPHandler struct {
	db *sql.DB
}

func NewHTTPHandler(db *sql.DB) *HTTPHandler {
	return &HTTPHandler{db: db}
}

func (h *HTTPHandler) GetIssues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := h.db.QueryContext(ctx, "SELECT id, title, status, priority, updated_at FROM issues")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var issues []Issue
	for rows.Next() {
		var issue Issue
		if err := rows.Scan(&issue.ID, &issue.Title, &issue.Status, &issue.Priority, &issue.UpdatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		issues = append(issues, issue)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error reading rows: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(issues); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *HTTPHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	db, err := sql.Open("postgres", "DATABASE_URL")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	h := NewHTTPHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("/issues", h.GetIssues)
	mux.HandleFunc("/healthz", h.Healthz)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
