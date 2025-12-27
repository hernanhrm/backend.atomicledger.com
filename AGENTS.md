# Agent Guidelines for backend.atomicledger.com

This document provides essential information for AI coding agents working on this Go backend project.

## Project Overview

**Tech Stack:**
- **Language:** Go 1.25.1
- **Module:** `backend.atomicledger.com`
- **Web Framework:** Echo v4
- **Database:** PostgreSQL with pgx/v5
- **DI Container:** samber/do/v2
- **Error Handling:** samber/oops
- **Logging:** log/slog via custom Logger interface
- **Testing:** testify/assert

**Architecture:**
- Dependency injection pattern for service management
- Interface-based abstractions (Database, Logger)
- Package-based organization (pkg/ for libraries, internal/ for app code, cmd/ for executables)
- Structured logging with context propagation

## Build, Test, and Lint Commands

### Build
```bash
make build              # Build API binary to bin/api
go build -o bin/api ./cmd/api
```

### Testing
```bash
make test               # Run all tests with verbose output
make test-cover         # Run tests with coverage

# Run tests in specific package
go test -v ./pkg/database

# Run a single test
go test -v ./pkg/database -run TestDatabase_NewConnection

# Run a single subtest
go test -v ./pkg/database -run TestDatabase_NewConnection/invalid_connection_string

# Run tests with race detector
go test -race ./...
```

### Linting and Formatting
```bash
make lint               # Run golangci-lint
make fmt                # Format all Go code
make vet                # Run go vet
make check              # Run fmt, vet, lint, and test

# Auto-fix linting issues where possible
golangci-lint run --fix
```

### Running Services
```bash
make run service=api    # Run API service
go run ./cmd/api        # Alternative
```

### Dependency Management
```bash
make tidy               # Clean and organize go.mod
go get <package>        # Add new dependency
```

## Code Style Guidelines

### Import Organization
Imports are organized in three sections using `gci` (enforced by golangci-lint):
1. Standard library
2. External packages
3. Internal packages (prefix: `backend.atomicledger.com`)

**Example:**
```go
import (
    "context"
    "fmt"
    
    "github.com/labstack/echo/v4"
    "github.com/samber/oops"
    
    "backend.atomicledger.com/pkg/database"
    "backend.atomicledger.com/pkg/logger"
)
```

### Package Documentation
Every package must have a doc comment in its main file or a dedicated `doc.go`:
```go
// Package database provides PostgreSQL database connection and query functionality.
package database
```

### Formatting
- Use `gofmt` standard formatting (tabs for indentation)
- Maximum line length: No hard limit, but keep it reasonable (~120 chars)
- Use `goimports` for automatic import management

### Types and Interfaces

**Interfaces:**
- Name interfaces with "Interface" suffix (e.g., `DatabaseInterface`, `PoolInterface`)
- Document all interface methods
- Use interface compile-time assertions:
```go
var _ DatabaseInterface = (*Database)(nil)
```

**Structs:**
- Unexported fields by default
- Document exported fields
- Use struct tags for JSON/DB mapping

**Constants:**
- Use typed constants with descriptive names
- Group related constants in const blocks
- Use camelCase with semantic prefixes:
```go
const (
    healthStatusHealthy = "healthy"
    sqlPreviewMaxLen    = 100
)
```

### Naming Conventions

- **Packages:** Short, lowercase, single word (e.g., `logger`, `database`, `sqlcraft`)
- **Files:** Lowercase with underscores for test files (e.g., `core.go`, `core_test.go`)
- **Types:** PascalCase (e.g., `SelectQuery`, `Database`)
- **Functions/Methods:** PascalCase for exported, camelCase for unexported
- **Variables:** camelCase for locals, PascalCase for exported
- **Interfaces:** PascalCase with "Interface" suffix
- **Constants:** camelCase with semantic prefixes

### Error Handling

Use `github.com/samber/oops` for rich error context:

```go
if err != nil {
    return oops.
        Code("db_query_failed").
        With("operation", "query").
        With("sql_preview", truncateSQL(sql)).
        Wrapf(err, "database query failed")
}
```

**Guidelines:**
- Always wrap errors with context using `oops.Wrapf()`
- Use descriptive error codes in snake_case
- Add relevant metadata with `With(key, value)`
- Never ignore errors without explicit reason (golangci-lint enforces this)
- Use `errors.New()` for creating new sentinel errors

### Logging

Use the custom `Logger` interface (backed by slog):

```go
// Structured logging with key-value pairs
log.Info("database connection established",
    "max_connections", poolConfig.MaxConns,
    "min_connections", poolConfig.MinConns,
)

log.Error("query execution failed",
    "error", err,
    "sql_preview", truncateSQL(sql),
)

// Create child loggers with context
requestLogger := log.With("request_id", requestID, "user_id", userID)
```

**Guidelines:**
- Always inject loggers via DI (never use global loggers)
- Use structured key-value pairs, not formatted strings
- Create child loggers with `With()` for context
- Choose appropriate levels: Debug, Info, Warn, Error
- Use `logger.NewNoop()` in tests

## Testing Guidelines

Use table-driven tests with subtests:

```go
func TestDatabase_NewConnection(t *testing.T) {
    tests := []struct {
        name        string
        connString  string
        expectError bool
    }{
        {
            name:        "invalid connection string",
            connString:  "invalid://connection",
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            log := logger.NewNoop()
            _, err := NewConnection(context.Background(), tt.connString, log)

            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**Guidelines:**
- Use `github.com/stretchr/testify/assert` for assertions
- Test files named `*_test.go` in the same package
- Use `logger.NewNoop()` to avoid test output noise
- Test edge cases (nil values, empty strings, etc.)
- Use `context.Background()` or `context.TODO()` in tests

## Common Patterns

### Dependency Injection
All services are registered in `pkg/di/providers.go` and resolved via the injector:
```go
injector := di.NewInjector()
log := di.MustInvokeLogger(injector)
db, err := do.InvokeNamed[*database.Database](injector, "database")
```

### Graceful Shutdown
Services implement shutdown for cleanup:
```go
func (db *Database) Shutdown(_ context.Context) error {
    db.logger.Info("shutting down database connection")
    if db.Pool != nil {
        db.Pool.Close()
    }
    return nil
}
```

### SQL Query Building
Use `pkg/sqlcraft` for type-safe SQL query construction:
```go
result, err := sqlcraft.
    Select("id", "name", "email").
    From("users").
    Where(filters...).
    OrderBy(sorts...).
    Limit(10).
    ToSQL()
```

## File Organization

```
backend.atomicledger.com/
├── cmd/                    # Executables (main packages)
│   └── api/               # API server entry point
├── internal/              # Private application code
│   ├── core/              # Core business logic
│   └── shared/            # Shared internal utilities
├── pkg/                   # Public libraries (reusable)
│   ├── database/          # Database connectivity
│   ├── di/                # Dependency injection
│   ├── logger/            # Logging abstraction
│   ├── server/            # HTTP server
│   └── sqlcraft/          # SQL query builder
├── .golangci.yml          # Linter configuration
├── Makefile               # Build commands
└── go.mod                 # Go module definition
```

When adding new packages, include a README.md explaining usage and architecture.
