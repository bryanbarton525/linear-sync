package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"linear-sync/internal/config"
	"linear-sync/internal/linear"
	"linear-sync/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client := linear.NewClient(cfg.APIKey)

	db, err := storage.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	store := storage.New(db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			issues, err := client.FetchIssues(ctx, cfg.TeamID)
			if err != nil {
				log.Printf("fetch: %v", err)
				continue
			}
			if err := store.Upsert(ctx, issues); err != nil {
				log.Printf("upsert: %v", err)
			}
		case <-ctx.Done():
			log.Println("shutting down")
			return
		}
	}
}
