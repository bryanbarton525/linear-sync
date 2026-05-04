package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := newDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := newStorage(db)
	client := newClient(cfg.APIKey)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[INFO] Starting linear-sync service")

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	doSync := func() {
		syncCtx, syncCancel := context.WithTimeout(ctx, 2*time.Minute)
		defer syncCancel()

		issues, err := client.fetchIssues(syncCtx, cfg.TeamID)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch issues: %v", err)
			return
		}

		if err := store.upsert(syncCtx, issues); err != nil {
			log.Printf("[ERROR] Failed to upsert issues: %v", err)
			return
		}

		log.Printf("[INFO] Synced %d issues", len(issues))
	}

	doSync()

	for {
		select {
		case <-ticker.C:
			doSync()
		case <-ctx.Done():
			log.Println("[INFO] Shutdown signal received")
			return
		}
	}
}
