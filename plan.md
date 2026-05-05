# Plan

> Authored by the Architect. The initial section below is the primary plan; remediation cycles append `## Remediation Cycle N` sections and never rewrite this header.

## Overview

This workflow implements a Go-based tool orchestration system with MCP integration, rate-limiting via token bucket algorithm, and a RESTful API for managing tools and workflows. The system uses idiomatic Go patterns including early returns, error wrapping, and channel-based concurrency. Key components include the rate limiter, tool registry, MCP server configuration, and workflow engine.

## Delivery Target

go-orca binary with embedded MCP tool discovery and rate-limited execution

## Tech Stack

- Go 1.22+
- context package
- sync/atomic
- net/http
- context.Context
- errgroup

## Components

| Name | Description | Inputs | Outputs |
|---|---|---|---|
| Rate Limiter | Token bucket rate limiter with leaky bucket replenishment logic using mutex-protected state. Implements context-aware allow() with periodic token replenishment based on elapsed time since last replenishment. | tokenBucketConfig | bool, time.Time |
| Tool Registry | In-memory registry for MCP and local tools with discovery, validation, and lifecycle management. Supports both streamable and SSE transport modes for MCP servers. | MCPConfig, localToolsConfig | []*ToolInfo, error |
| Workflow Engine | Orchestrates tool execution with rate limiting, context propagation, and structured logging. Supports fan-out patterns with errgroup for parallel tool execution. | []*ToolInfo, context.Context | *WorkflowResult, error |
| Configuration Manager | Manages go-orca.yaml configuration parsing, MCP server endpoints, and rate limit settings with validation. | string | *Config, error |

## Architectural Decisions

1. **Use token bucket with mutex protection rather than atomic operations**
   - Rationale: Token counts represent compound state requiring mutex. Atomic operations only work for single-value operations without side effects.
   - Tradeoffs: Slight performance overhead for mutex locking vs data race hazards with atomic+mutex mixing
2. **Use channels over shared memory for goroutine coordination**
   - Tradeoffs: Slightly higher latency vs shared memory but guarantees thread safety
3. **Implement error wrapping with %w rather than string concatenation**
   - Rationale: Standard library expects wrapped errors for diagnostics and errors.Is() patterns
   - Tradeoffs: Slightly more verbose but follows Go best practices

## Task Graph

| ID | Specialty | Title | Depends On | Description |
|---|---|---|---|---|
| 4c64351f | backend | Create rate limiter package | - | Produce artifact kind `code`, name `pkg/limiter/token_bucket.go`. Implement token bucket rate limiter with `Bucket` struct containing `tokens int64`, `burst int64`, `rate int64`, `lastReplenish int64`, `mutex sync.Mutex`. Include methods `Tokens(ctx context.Context) (int64, error)` and `Allow(ctx context.Context) bool`. Use `time.Now().UnixMilli()` for timestamps. Handle context cancellation with select on ctx.Done(). Include sentinel error `ErrBucketFull` and `ErrBucketEmpty` at package level. Add table-driven tests in `token_bucket_test.go` covering happy path, context cancellation, and rate limit exhaustion scenarios. |
| c8aa12e0 | backend | Create tool registry package | pkg/limi | Produce artifact kind `code`, name `pkg/registry/tools.go`. Implement `ToolRegistry` struct with `mcpServers []*MCPConfig`, `localTools []LocalToolConfig`, `tools map[string]*ToolInfo`. Include methods `Load(ctx context.Context) error` for MCP discovery and `RegisterLocalTool(name string, tool LocalTool) error`. Use `sync.RWMutex` for read-heavy access patterns. Implement `Get(name string) (*ToolInfo, error)` and `Execute(ctx context.Context, name string, args map[string]any) (*ToolResult, error)`. Add tests in `tools_test.go` covering MCP server registration and local tool execution. |
| 21dc9fd7 | backend | Create workflow engine package | pkg/limi, pkg/regi | Produce artifact kind `code`, name `pkg/engine/workflow.go`. Implement `WorkflowEngine` struct with `limiter *RateLimiter`, `registry *ToolRegistry`, `log *log.Logger`. Include methods `Run(ctx context.Context, tools []string) (*WorkflowResult, error)` and `RunParallel(ctx context.Context, tools []string) (*WorkflowResult, error)`. Use `errgroup.Group` for parallel execution. Emit structured logs with `zap` or standard `log`. Return proper HTTP status codes for API endpoints. Add tests in `workflow_test.go` using `httptest.NewServer()` with deferred cleanup. |
| 33b0a08e | backend | Create REST API handlers | pkg/engi, pkg/regi | Produce artifact kind `code`, name `handlers/tools.go`. Implement HTTP handlers for `/tools/register`, `/tools/{name}`, `/workflow/run` endpoints. Use `httptest.NewServer()` for testing with `defer ts.Close()`. Return proper HTTP status codes (201 for create, 200 for success, 500 for errors). Use `json.NewEncoder(w)` for response serialization. Add tests in `tools_handler_test.go` using `testing.T` with table-driven approach. |
| b19a9e92 | ops | Create configuration management package | pkg/limi, pkg/regi | Produce artifact kind `code`, name `pkg/config/config.go`. Implement `Config` struct with `Tools MCPConfig`, `RateLimit RateLimitConfig`. Add YAML parsing using `yaml.Unmarshal`. Implement `Validate()` method for schema constraints. Include `cmd/init` main package that parses `go-orca.yaml` and registers MCP servers and local tools. |
| 2468d012 | ops | Create main application entry point | pkg/conf, handlers, pkg/engi | Produce artifact kind `code`, name `cmd/orca/main.go`. Implement `main()` function that loads configuration, initializes registry and engine, starts HTTP server, and handles shutdown gracefully. Use `sync.WaitGroup` for signal handling. Emit startup logs with workflow summary. Include graceful shutdown with context cancellation on OS signals. |
| 1412e762 | writer | Create documentation for MCP configuration | - | Produce artifact kind `markdown`, name `docs/mcp-configuration.md`. Write 300–500 word guide covering MCP server configuration in go-orca.yaml, available transport modes (streamable, sse, command), and tool discovery process. Include code example showing `mcp-config.yaml` with `streamable` and `command` transport examples. Reference [mcp-server-catalog.md](references/mcp-server-catalog.md) for server list. Use clear headings and code blocks with language specification. |
| 56cea65c | writer | Create README with usage examples | docs/mcp | Produce artifact kind `markdown`, name `README.md`. Write 800–1200 word README covering project overview, installation, configuration examples, rate limiting behavior, and tool execution workflows. Include `go-orca.yaml` example configuration and CLI usage examples. Reference `docs/mcp-configuration.md` for MCP details. Use active voice and clear call-to-action sections. |

