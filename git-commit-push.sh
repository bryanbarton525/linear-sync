#!/bin/bash
set -e

cd /var/lib/go-orca/workspaces/13ebcc13-3f88-4186-bc3c-2ecb581e4ceb

# Add all source files
git add go.mod go.sum main.go config.go linear.go storage.go config_test.go linear_test.go storage_test.go

# Commit with the specified message and co-author trailer
git commit -m "Add Linear.app to PostgreSQL sync service (workflow 13ebcc13-3f88-4186-bc3c-2ecb581e4ceb)

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"

# Push to main branch
git push origin main

echo "Successfully pushed 9 files to github.com/bryanbarton525/linear-sync"
