# Logger Package

A structured logging abstraction for backend.atomicledger.com that provides a clean interface without direct dependencies on specific logging implementations.

## Design Philosophy

- **Interface-based**: Application code depends on `Logger` interface, not concrete implementations
- **Dependency Injection**: Loggers are injected into components for testability and flexibility
- **Structured Logging**: All logging uses key-value pairs for better observability
- **Swappable Implementations**: Can switch from slog to other loggers without changing application code

## Architecture

```
Application Code
       ↓
Logger Interface (abstraction)
       ↓
SlogAdapter (implementation using log/slog)
```

## Usage

### Basic Usage

```go
import "backend.atomicledger.com/pkg/logger"

// Create a logger
log := logger.NewDevelopment()

// Log messages with structured data
log.Info("user logged in", "user_id", 123, "username", "john")
log.Error("failed to process request", "error", err, "request_id", reqID)
```

### Dependency Injection

```go
// In your main.go
func main() {
    log := logger.NewProduction()
    
    // Inject into components
    db, err := database.NewConnection(ctx, connString, log)
    config, err := localconfig.GetConfig(log)
}
```

### Creating Child Loggers

```go
// Create a base logger with common fields
requestLogger := log.With("request_id", requestID, "user_id", userID)

// All logs from this logger will include those fields
requestLogger.Info("processing request")
requestLogger.Error("request failed", "error", err)
```

### Environment-based Configuration

```go
func initLogger() logger.Logger {
    env := os.Getenv("APP_ENV")
    
    if env == "production" {
        return logger.NewProduction() // JSON format, Info level
    }
    return logger.NewDevelopment() // Text format, Debug level
}
```

## Available Implementations

### SlogAdapter (Default)

Uses Go's standard `log/slog` package. Two convenience constructors:

- `NewDevelopment()`: Text format, Debug level - great for local development
- `NewProduction()`: JSON format, Info level - optimized for production logging systems

### NoopLogger

A logger that does nothing. Useful for testing:

```go
func TestMyFunction(t *testing.T) {
    log := logger.NewNoop()
    result := MyFunction(log)
    // ... assertions
}
```

## Testing

When testing components that require a logger:

```go
func TestDatabaseConnection(t *testing.T) {
    // Use NoopLogger to avoid log output during tests
    log := logger.NewNoop()
    
    db, err := database.NewConnection(ctx, connString, log)
    // ... test assertions
}
```

## Custom Implementations

To create a custom logger implementation:

1. Implement the `Logger` interface
2. Provide all required methods: `Debug`, `Info`, `Warn`, `Error`, `With`, `WithContext`
3. Ensure thread-safety for concurrent use

Example:

```go
type MyCustomLogger struct {
    // your fields
}

func (l *MyCustomLogger) Info(msg string, keysAndValues ...interface{}) {
    // your implementation
}

// ... implement other methods
```

## Best Practices

1. **Always inject loggers**: Don't create loggers inside packages, inject them
2. **Use structured fields**: Always log with key-value pairs, not formatted strings
3. **Create child loggers**: Use `With()` to add context rather than repeating fields
4. **Choose appropriate levels**:
   - `Debug`: Detailed information for debugging
   - `Info`: General informational messages
   - `Warn`: Warning messages that don't stop execution
   - `Error`: Error conditions that need attention

## Example Integration

```go
package main

import (
    "context"
    "os"
    
    "backend.atomicledger.com/pkg/database"
    "backend.atomicledger.com/pkg/localconfig"
    "backend.atomicledger.com/pkg/logger"
)

func main() {
    // Initialize logger based on environment
    var log logger.Logger
    if os.Getenv("APP_ENV") == "production" {
        log = logger.NewProduction()
    } else {
        log = logger.NewDevelopment()
    }
    
    log.Info("application starting")
    
    // Load config with logger
    config, err := localconfig.GetConfig(log)
    if err != nil {
        log.Error("config load failed", "error", err)
        os.Exit(1)
    }
    
    // Connect to database with logger
    db, err := database.NewConnection(context.Background(), connString, log)
    if err != nil {
        log.Error("database connection failed", "error", err)
        os.Exit(1)
    }
    defer db.Close()
    
    log.Info("application ready")
}
```
