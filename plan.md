# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This workflow deploys a complete Linear.app to PostgreSQL synchronization service by writing 9 pre-specified Go source files to the workspace, building the module with the Go toolchain, running the test suite to verify correct behavior, and pushing the codebase to the target GitHub repository. No architectural decisions are required—all source code is provided verbatim and must be written exactly as specified.

## Delivery Target

Complete Go module pushed to github.com/bryanbarton525/linear-sync main branch with all 9 source files, successful build, and verified test outcomes

## Tech Stack

- Go 1.21
- github.com/lib/pq v1.10.9 (PostgreSQL driver)
- github.com/stretchr/testify v1.9.0 (testing assertions)
- Git version control
- GitHub remote repository

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Source File Deployment | Write all 9 Go source files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) to the workspace with exact character-for-character content matching the specification. No modifications are permitted. | Verbatim file content from workflow request | 9 Go source files in workspace root |
| Build Verification | Execute Go build with -buildvcs=false flag to compile the module. Verify exit code 0 and presence of compiled binary artifact. Build must succeed without errors using the provided go.mod and go.sum. | 9 source files in workspace | Compiled binary executable, Build success confirmation |
| Test Execution and Validation | Run Go test suite with -buildvcs=false flag in verbose mode. Parse test output to confirm exactly 4 PASS results (config_test.go and linear_test.go) and 2 SKIP results (storage_test.go tests that require PostgreSQL). Fail if counts don't match specification. | Compiled module, Test files in workspace | Test execution report with pass/skip counts |
| Git Repository Management | Initialize git repository if not present, configure remote URL to github.com/bryanbarton525/linear-sync, create main branch if needed, commit all changes with Co-authored-by trailer, and push to main. Handle git configuration errors without force-pushing. | All source files in workspace, Build and test verification | Git repository committed and pushed to remote main branch |

## Architectural Decisions

1. **Write all 9 files in a single parallel batch before executing any build or test commands**
   - Rationale: Matriarch specified parallel file writing to maximize efficiency and ensure the workspace is complete before build/test phases. This prevents partial workspace states and reduces total execution time.
   - Tradeoffs: Single batch write requires all file content to be available upfront, but this is satisfied since all content is provided in the workflow request specification.
2. **Parse test output to enforce exactly 4 PASS and 2 SKIP results**
   - Rationale: The specification explicitly requires this exact test outcome. Any deviation (e.g., additional tests appearing, unexpected failures, or different skip counts) indicates a problem with the workspace or test execution that must be caught and reported.
   - Tradeoffs: Strict validation prevents silent failures but requires precise output parsing. However, the test suite is fixed so output format is predictable.
3. **Do not force-push if remote diverges from local history**
   - Rationale: Force-pushing can destroy remote work and should require explicit user approval. If the remote has commits not in local history, treat as a blocking issue and fail safely.
   - Tradeoffs: This may block deployment if the remote repository already has content, but it prevents data loss and ensures safe git operations.

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 1390de8d | backend | Write All 9 Source Files to Workspace | - | Produce artifact kind `code` for 9 separate files: go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go. Each file must be written to the workspace root directory (/var/lib/go-orca/workspaces/11760721-25da-47cd-86ed-547a652c5dc4) with content matching the workflow specification exactly—byte-for-byte identical with no modifications, additions, deletions, or whitespace normalization. Use the write_file tool for each file. Acceptance criteria: All 9 files exist in workspace root with exact content from specification. Quality standard: Character-for-character match (no encoding changes, no line ending conversion). Workspace materialization: Files must be written to workspace root, not a subdirectory. This task must write all files before any build or test commands are executed. |
| 8d08d197 | backend | Verify All Source Files Present | 1390de8d | Produce artifact kind `validation_report` named `file_presence_check.txt`. List the workspace directory and verify that all 9 files (go.mod, go.sum, main.go, config.go, linear.go, storage.go, config_test.go, linear_test.go, storage_test.go) exist in the workspace root. Acceptance criteria: All 9 files must be present; if any file is missing, fail with a blocking error listing which files are absent. Quality standard: Use 'ls -la' or equivalent to list files and confirm presence. This verification step ensures the workspace is complete before attempting build. |
| e5480b17 | backend | Build Go Module | 8d08d197 | Produce artifact kind `binary` named `linear-sync` (the compiled executable). Execute 'go build -buildvcs=false .' in the workspace directory (/var/lib/go-orca/workspaces/11760721-25da-47cd-86ed-547a652c5dc4). Acceptance criteria: Build must complete with exit code 0, produce no compilation errors, and generate a binary executable file in the workspace root. Quality standard: All imports must resolve using the provided go.mod and go.sum without running go mod tidy or go mod init. Verify the binary artifact exists after build completes. If build fails, capture the full error output and fail the task with a blocking error. |
| cabbef48 | backend | Run Tests and Verify Outcomes | e5480b17 | Produce artifact kind `test_report` named `test_results.txt`. Execute 'go test -buildvcs=false . -v' in the workspace directory. Parse the test output and count PASS and SKIP results. Acceptance criteria: Exactly 4 PASS results (from config_test.go and linear_test.go) and exactly 2 SKIP results (from storage_test.go tests that skip without PostgreSQL). If counts do not match, fail with a blocking error showing actual vs expected counts. Quality standard: Capture full verbose test output. Record which tests passed and which were skipped. Tests must execute with exit code 0 or 1 (exit 1 is acceptable if skips are present). Verify no test failures (FAIL results). |
| f67454ac | ops | Initialize Git Repository | cabbef48 | Produce artifact kind `git_repository` in workspace root. Check if .git directory exists. If not, run 'git init' to initialize a new repository. If git is already initialized, verify it is a valid repository by running 'git status'. Acceptance criteria: A valid .git directory must exist in the workspace root. Quality standard: After initialization or verification, 'git status' must execute without errors. If git init fails, capture the error and fail the task with a blocking error. |
| ad04baf9 | ops | Configure Git Remote and Branch | f67454ac | Produce artifact kind `git_config` for remote and branch setup. Configure the git remote named 'origin' to point to https://github.com/bryanbarton525/linear-sync. If remote exists with a different URL, update it using 'git remote set-url origin'. Verify the main branch exists locally; if not, create it from the current HEAD using 'git checkout -b main' or 'git branch -M main'. Acceptance criteria: 'git remote -v' must show origin pointing to https://github.com/bryanbarton525/linear-sync for both fetch and push. 'git branch' must show main branch exists. Quality standard: Use 'git remote get-url origin' to verify URL matches specification exactly. |
| 9273952c | ops | Commit All Changes | ad04baf9 | Produce artifact kind `git_commit` with a descriptive commit message. Stage all 9 source files using 'git add .'. Create a commit with message 'Initial commit: Linear.app to PostgreSQL sync service' and include the Co-authored-by trailer 'Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>' at the end of the commit message. Acceptance criteria: 'git log -1' must show a commit with the specified message and trailer. 'git status' must show a clean working tree with no uncommitted changes. Quality standard: All 9 files must be included in the commit. Use 'git diff --cached' before committing to verify all expected files are staged. |
| d825cb2d | ops | Push to GitHub Main Branch | 9273952c | Produce artifact kind `deployment` by pushing the committed code to the remote repository. Execute 'git push -u origin main' to push the main branch to github.com/bryanbarton525/linear-sync. Acceptance criteria: Push must succeed with exit code 0. The remote main branch must reflect the committed changes. Quality standard: If push is rejected due to divergent remote history (e.g., remote has commits not in local), fail with a blocking error—do not force-push without explicit user approval. If push fails due to authentication errors, fail with a clear error message indicating credentials must be configured. After successful push, verify using 'git ls-remote origin main' that the remote ref matches the local commit. |

