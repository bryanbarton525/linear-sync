package main

import (
	"context"
	"log"
	"os"
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
		log.Fatal(err)
	}

	db, err := storage.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := storage.New(db)
	client := linear.NewClient(cfg.APIKey)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[INFO] Starting linear-sync service")

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	sync := func() {
		syncCtx, syncCancel := context.WithTimeout(ctx, 2*time.Minute)
		defer syncCancel()

		issues, err := client.FetchIssues(syncCtx, cfg.TeamID)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch issues: %v", err)
			return
		}

		if err := store.Upsert(syncCtx, issues); err != nil {
			log.Printf("[ERROR] Failed to upsert issues: %v", err)
			return
		}

		log.Printf("[INFO] Synced %d issues", len(issues))
	}

	sync()

	for {
		select {
		case <-ticker.C:
			sync()
		case <-ctx.Done():
			log.Println("[INFO] Shutdown signal received")
			goto shutdown
		}
	}

shutdown:
	log.Println("[INFO] Shutdown complete")
}
