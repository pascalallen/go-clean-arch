# go-clean-arch

A reference implementation of Clean Architecture in Go. This repository is the companion project for the Medium article [*Structuring Go Projects With Clean Architecture*](https://medium.com/@pascalallen/structuring-go-projects-with-clean-architecture-2c46d7e58ac3).

## Project Structure

```
├── bin/       # Helper scripts for common Docker commands
├── cmd/       # Go main packages (entry points)
├── docs/      # Documentation
├── internal/  # All application code
└── web/       # Frontend assets
```

## Architecture

The application code in `internal/` is organized around three layers:

### Domain Layer

Contains the core business logic — entities, repository interfaces, and service interfaces. No framework dependencies.

```
internal/app/domain/
├── logger/
├── pagination/
├── password/
├── permission/
├── role/
└── user/
    ├── user.go        # User entity with behavior methods
    └── repository.go  # Repository interface
```

### Application Layer

Orchestrates the domain using CQRS and event-driven patterns.

```
internal/app/application/
├── command/           # Intent structs (RegisterUser, DeleteUser)
├── command_handler/   # Command handlers
├── event/             # Domain events (UserRegistered)
├── listener/          # Event listeners
├── query/             # Query structs (GetUserById, ListUsers)
└── query_handler/     # Query handlers
```

### Infrastructure Layer

Concrete implementations of domain interfaces.

```
internal/app/infrastructure/
├── container/    # Google Wire DI container
├── database/     # Postgres session, migrations, seeders
├── messaging/    # Command bus, event dispatcher, query bus
├── repository/   # Postgres implementations
├── routes/       # Gin HTTP router
├── service/      # JWT token service
└── websocket/    # WebSocket hub for real-time push updates
```

## Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Development Setup

### Clone the repository

```bash
git clone https://github.com/pascalallen/go-clean-arch.git && cd go-clean-arch
```

### Copy the environment file

```bash
cp .env.example .env
```

### Bring up the environment

```bash
bin/up
```

The app will be running at [http://localhost:8080](http://localhost:8080).

### Seed the database

```bash
bin/exec go run cmd/seed/seed.go
```

### Take down the environment

```bash
bin/down
```

## Testing

```bash
docker compose exec go go test ./... -covermode=count -coverprofile=coverage.out
```

## API

| Method | Endpoint                    | Description        |
|--------|-----------------------------|--------------------|
| POST   | `/api/v1/auth/register`    | Register a user    |
| POST   | `/api/v1/auth/login`       | Log in             |

## License

[MIT](LICENSE)
