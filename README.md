# migrate

A fork of [golang-migrate/migrate](https://github.com/golang-migrate/migrate) — Database migrations written in Go. Use as CLI or import as library.

[![Go Reference](https://pkg.go.dev/badge/github.com/your-org/migrate.svg)](https://pkg.go.dev/github.com/your-org/migrate)
[![CI](https://github.com/your-org/migrate/actions/workflows/ci.yaml/badge.svg)](https://github.com/your-org/migrate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/migrate)](https://goreportcard.com/report/github.com/your-org/migrate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Stateless** — no need for a separate migration tracking table (uses a single version table)
- **Multiple database drivers** — PostgreSQL, MySQL, SQLite, MongoDB, and more
- **Multiple source drivers** — filesystem, Go embed, S3, GitHub, and more
- **CLI & library** — use from the command line or import as a Go package
- **Up & Down migrations** — apply and rollback migrations
- **Dirty state detection** — detects and reports failed migrations

## Supported Databases

| Database   | Driver Import Path                          |
|------------|---------------------------------------------|
| PostgreSQL | `github.com/your-org/migrate/database/postgres` |
| MySQL      | `github.com/your-org/migrate/database/mysql`    |
| SQLite3    | `github.com/your-org/migrate/database/sqlite3`  |
| MongoDB    | `github.com/your-org/migrate/database/mongodb`  |

## Supported Sources

| Source     | Driver Import Path                         |
|------------|--------------------------------------------|
| File       | `github.com/your-org/migrate/source/file`  |
| Go Embed   | `github.com/your-org/migrate/source/iofs`  |
| GitHub     | `github.com/your-org/migrate/source/github`|

## Installation

### CLI

```bash
go install github.com/your-org/migrate/cmd/migrate@latest
```

Or download a pre-built binary from the [releases page](https://github.com/your-org/migrate/releases).

### Library

```bash
go get github.com/your-org/migrate/v4
```

## Quick Start

### CLI Usage

```bash
# Apply all up migrations
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" up

# Rollback the last migration
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" down 1

# Check current migration version
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" version

# Force set version (use with caution)
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" force 1
```

### Library Usage

```go
import (
    "github.com/your-org/migrate/v4"
    _ "github.com/your-org/migrate/v4/database/postgres"
    _ "github.com/your-org/migrate/v4/source/file"
)

func main() {
    m, err := migrate.New(
        "file://migrations",
        "postgres://localhost:5432/mydb?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

## Migration Files

Migration files follow the naming convention:

```
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

Example:

```
migrations/
  000001_create_users_table.up.sql
  000001_create_users_table.down.sql
  000002_add_email_index.up.sql
  000002_add_email_index.down.sql
```

## Development

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run

# Build CLI
go build -o migrate ./cmd/migrate
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](LICENSE)
