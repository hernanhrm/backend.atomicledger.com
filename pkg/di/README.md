# Dependency Injection Package

This package provides dependency injection for the backend.atomicledger.com application using [`github.com/samber/do/v2`](https://github.com/samber/do).

## Overview

The DI container manages the lifecycle of all application services, providing:
- Type-safe service registration and resolution
- Automatic dependency graph resolution
- Graceful shutdown in correct dependency order
- Easy testing with service mocking

## Architecture

```
Injector (Root Scope)
  ├─> Logger (Lazy Singleton)
  ├─> Config (Lazy Singleton)
  ├─> Database (Lazy Singleton)
  └─> Server (Lazy Singleton)
```

## Usage

### Basic Usage

```go
// Create DI container
injector := di.NewInjector()

// Get services
logger := di.MustInvokeLogger(injector)
server, err := do.InvokeNamed[*server.Server](injector, "server")

// Graceful shutdown
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
di.Shutdown(ctx, injector)
```

### Adding New Services

1. **Create the service** in its own package (e.g., `pkg/cache`)
2. **Add a provider function** in `providers.go`:

```go
func ProvideCache() func(do.Injector) {
    return do.LazyNamed[*cache.Cache]("cache", func(i do.Injector) (*cache.Cache, error) {
        log := do.MustInvokeNamed[logger.Logger](i, "logger")
        config := do.MustInvokeNamed[*localconfig.ConfigService](i, "config")
        
        return cache.NewCache(config, log)
    })
}
```

3. **Register in `container.go`**:

```go
func NewInjector() do.Injector {
    injector := do.New()
    
    // ... existing providers
    ProvideCache()(injector)
    
    return injector
}
```

### Named Services

For multiple instances of the same type (e.g., multiple databases):

```go
// Register
ProvideNamedDatabase("primary")(injector)
ProvideNamedDatabase("analytics")(injector)

// Use
primaryDB := do.MustInvokeNamed[*database.Database](injector, "primary")
analyticsDB := do.MustInvokeNamed[*database.Database](injector, "analytics")
```

## Service Lifecycles

| Lifecycle | Description | Use Case |
|-----------|-------------|----------|
| Lazy | Created on first use | Most services (database, server) |
| Eager | Created immediately | Critical validation (config) |
| Transient | Created every time | Request-scoped services |

### Implementing Shutdown

Services that need cleanup should implement graceful shutdown:

```go
type MyService struct {
    // fields
}

func (s *MyService) Shutdown(ctx context.Context) error {
    // cleanup logic
    return nil
}
```

The DI container automatically calls `Shutdown()` in reverse dependency order.

## Testing

### Test Injector

```go
func TestMyHandler(t *testing.T) {
    injector := do.New()
    
    // Provide test services
    do.ProvideNamed(injector, "logger", func(i do.Injector) (logger.Logger, error) {
        return logger.NewNoop(), nil
    })
    
    // Test code
}
```

### Mocking Services

```go
func TestWithMockDB(t *testing.T) {
    injector := do.New()
    
    // Override with mock
    mockDB := &MockDatabase{}
    do.OverrideNamed(injector, "database", func(i do.Injector) (*database.Database, error) {
        return mockDB, nil
    })
    
    // Your test
}
```

## Environment Configuration

The DI container respects the `APP_ENV` environment variable:

- `production` → JSON logging, optimized settings
- `development` (default) → Text logging, debug mode

Access helpers:
```go
di.GetEnvironment()  // Returns current env
di.IsProduction()    // Returns true if prod
di.IsDevelopment()   // Returns true if dev
```

## Error Handling

All providers use `oops` for rich error context:

```go
func ProvideService() func(do.Injector) {
    return do.LazyNamed[*MyService]("service", func(i do.Injector) (*MyService, error) {
        dep, err := do.InvokeNamed[*Dependency](i, "dependency")
        if err != nil {
            return nil, oops.
                Code("service_dependency_failed").
                With("dependency", "dependency").
                Wrapf(err, "failed to resolve dependency")
        }
        
        return NewMyService(dep)
    })
}
```

## Best Practices

1. **Use named services** for clarity and multiple instances
2. **Lazy load by default** unless eager loading is required
3. **Implement Shutdown** for resources needing cleanup
4. **Keep providers simple** - complex logic belongs in service constructors
5. **Use MustInvoke** only in main(), use Invoke elsewhere
6. **Test with real injector** when possible, mock sparingly

## Future Enhancements

- [ ] Health check aggregation endpoint
- [ ] Service dependency visualization
- [ ] Hot-reload for configuration
- [ ] Request-scoped injectors for per-request dependencies
- [ ] Metrics integration (service invocation counts, etc.)
