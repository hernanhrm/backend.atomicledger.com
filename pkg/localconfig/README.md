# LocalConfig Package

Flexible configuration management using godotenv for .env files, designed for monorepo environments.

## Features

- **Configurable .env file location** via environment variables
- **Support for multiple env file names** (`.env`, `.env.dev`, `.env.production`, etc.)
- **Automatic fallback to system environment variables** (perfect for production)
- **Monorepo-friendly** - each service can have its own .env file

## Configuration Structure

```go
type LocalConfig struct {
    Service  Service
    Database Database
}

type Service struct {
    Port int
    Name string
}

type Database struct {
    Host     string
    Port     int
    Username string
    Password string
    Name     string
    SSLMode  string
}
```

## Environment Variables

### Configuration File Location (meta variables)
- `CONFIG_ENV_PATH` - Relative path to directory containing the .env file (default: `"."`)
- `CONFIG_ENV_FILENAME` - Name of the env file to load (default: `".env"`)

### Service
- `SERVICE_PORT` - API server port
- `SERVICE_NAME` - Service name

### Database  
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USERNAME` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSL_MODE` - SSL mode (disable, require, verify-ca, verify-full)

## Usage

### Basic Usage (Default Behavior)

```go
// Loads ./.env by default
config := localconfig.GetConfig()

fmt.Printf("Starting %s on port %d\n", config.Service.Name, config.Service.Port)
```

### Monorepo Usage with Environment Variables

Perfect for local development in a monorepo where each service has its own `.env` file:

```bash
# Terminal 1: Run API service
CONFIG_ENV_PATH="./cmd/api" go run ./cmd/api

# Terminal 2: Run worker service  
CONFIG_ENV_PATH="./cmd/worker" go run ./cmd/worker

# Terminal 3: Run webhook service with custom env file
CONFIG_ENV_PATH="./cmd/webhook" CONFIG_ENV_FILENAME=".env.dev" go run ./cmd/webhook
```

Project structure:
```
backend.atomicledger.com/
├── cmd/
│   ├── api/
│   │   ├── .env
│   │   └── main.go
│   ├── worker/
│   │   ├── .env
│   │   └── main.go
│   └── webhook/
│       ├── .env.dev
│       ├── .env.production
│       └── main.go
└── pkg/
    └── localconfig/
```

### Programmatic Usage

For scenarios where you need to override configuration programmatically:

```go
// Load from a specific directory with custom filename
config := localconfig.GetConfigWithOptions(localconfig.ConfigOptions{
    EnvPath:     "./cmd/api",
    EnvFileName: ".env.production",
})
```

### Production Usage

In production, set environment variables directly (no .env file needed):

```bash
export SERVICE_PORT=8080
export SERVICE_NAME=api-prod
export DB_HOST=prod.database.com
export DB_PORT=5432
# ... etc

go run ./cmd/api
```

The package will automatically fall back to system environment variables if the .env file is not found.

## Local Development Setup

### Option 1: Shell Environment Variables

Add to your `.zshrc` or `.bashrc`:

```bash
export CONFIG_ENV_PATH="./cmd/api"
export CONFIG_ENV_FILENAME=".env"
```

### Option 2: Per-Command

```bash
CONFIG_ENV_PATH="./cmd/api" go run ./cmd/api
```

### Option 3: Makefile

```makefile
.PHONY: run-api
run-api:
	CONFIG_ENV_PATH="./cmd/api" go run ./cmd/api

.PHONY: run-worker
run-worker:
	CONFIG_ENV_PATH="./cmd/worker" go run ./cmd/worker
```

## Advanced: Multiple Databases

If you need multiple database configurations (e.g., primary, replica, cache), you can extend the `LocalConfig` struct:

```go
type AppConfig struct {
    Service   localconfig.Service
    PrimaryDB localconfig.Database
    ReplicaDB ReplicaDatabase
}

type ReplicaDatabase struct {
    Host     string
    Port     int
    Username string
    Password string
    Name     string
    SSLMode  string
}
```

Then in your `.env`:
```env
# Primary Database
DATABASE_HOST=primary.db.com
DATABASE_PORT=5432
DATABASE_USERNAME=user
DATABASE_PASSWORD=pass

# Replica Database
REPLICA_HOST=replica.db.com
REPLICA_PORT=5432
REPLICA_USERNAME=user
REPLICA_PASSWORD=pass
```

You'll need to extend the loading logic in a similar way to how `GetConfigWithOptions()` loads the standard config, manually mapping environment variables to your custom struct fields.
