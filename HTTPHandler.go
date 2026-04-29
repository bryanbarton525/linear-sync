package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

type HealthResponse struct {
    Status string `json:"status"`
}

type HTTPHandler struct {}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/healthz":
        h.handleHealth(w, r)
    case "/sync":
        // Existing /sync handler logic will go here
    default:
        http.NotFound(w, r)
    }
}

func (h *HTTPHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
    health := HealthResponse{
        Status: "healthy",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}

func main() {
    handler := &HTTPHandler{}
    http.Handle("/", handler)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Printf("Starting server on :%s\n", port)
    if err := http.ListenAndServe(":